"""
Apache beam pipeline definitions for processing reddit submisisons
"""

import os
import json
from typing import List, Dict

import apache_beam as beam
from apache_beam.pvalue import PCollection
from apache_beam.options.pipeline_options import PipelineOptions

from util.tickers import parse_tickers
from util.constants.scraping import (
    ScrapedPostConstants as SPC,
    TickerMentionsConstants as TMC,
    MentionTypes,
)
from util.constants.pubsub import Topics, BeamSubscriptions as BS
from util.constants.bigquery import Tables

from . import print_collection, publish_to_pubsub, standard_options

LOCATION = os.path.realpath(os.path.join(os.getcwd(), os.path.dirname(__file__)))


def parse_ticker_mentions(scraped_post: Dict) -> List[Dict]:
    """Parse ticker mention objects from a scraped post object"""

    ticker_mentions = []
    for ticker in parse_tickers(scraped_post[SPC.TEXT]):
        ticker_mention = {
            TMC.TICKER: ticker,
            TMC.DATA_SOURCE: scraped_post[SPC.DATA_SOURCE],
            TMC.PARENT_SOURCE: scraped_post[SPC.PARENT_SOURCE],
            TMC.PARENT_ID: scraped_post[SPC.PARENT_ID],
            TMC.TIMESTAMP: scraped_post[SPC.TIMESTAMP],
            TMC.MENTION_TYPE: MentionTypes.TICKER,
        }

        ticker_mentions.append(ticker_mention)

    return ticker_mentions


def extract_ticker_mentions(scraped_posts: PCollection) -> PCollection:
    """Apply the parse_ticker_mentions method to a scraped_posts collection"""
    return scraped_posts | beam.FlatMap(parse_ticker_mentions)


def create_test_scraped_posts_pipeline():
    """
    Create the local test pipeline that uses data
    saved in a json file inside the 'test' directory
    """

    # Load test data from json
    with open(f"{LOCATION}/test/scraped_posts.json", "r") as test_file:
        test_data = json.load(test_file)

    # Create pipeline
    pipeline_options = PipelineOptions(standard_options)

    with beam.Pipeline(options=pipeline_options) as pipeline:
        scraped_posts = pipeline | beam.Create(test_data)
        ticker_mentions = extract_ticker_mentions(scraped_posts)
        print_collection(scraped_posts)
        print_collection(ticker_mentions)


def create_scraped_posts_pipeline():
    """
    Create the local test pipeline that uses data
    saved in a json file inside the 'test' directory
    """

    # Create pipeline
    pipeline_options = PipelineOptions(
        [*standard_options, "--direct_num_workers=0", "--streaming"]
    )

    with beam.Pipeline(options=pipeline_options) as pipeline:
        # Load from Pub/Sub
        scraped_posts = (
            pipeline
            | beam.io.ReadFromPubSub(subscription=f"{BS.PREFIX}/{BS.SCRAPED_POSTS}")
            | beam.Map(lambda d: d.decode("utf-8"))
            | beam.Map(json.loads)
        )

        # Write posts to bigquery
        _ = scraped_posts | "Save Posts to BQ" >> beam.io.WriteToBigQuery(
            Tables.SCRAPED_POSTS,
            write_disposition=beam.io.BigQueryDisposition.WRITE_APPEND,
        )

        # Extract ticker_mentions
        ticker_mentions = extract_ticker_mentions(scraped_posts)

        # Write mentions to bigquery
        _ = ticker_mentions | "Save Tickers to BQ" >> beam.io.WriteToBigQuery(
            Tables.TICKER_MENTIONS,
            write_disposition=beam.io.BigQueryDisposition.WRITE_APPEND,
        )

        # Write to Pub/Sub
        publish_to_pubsub(ticker_mentions, f"{Topics.PREFIX}/{Topics.TICKER_MENTIONS}")
