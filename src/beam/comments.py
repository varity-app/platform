"""
Apache beam pipeline definitions for processing reddit submisisons
"""

import os
import json
from typing import List, Dict

import apache_beam as beam
from apache_beam.pvalue import PCollection
from apache_beam.options.pipeline_options import PipelineOptions

from util.constants.reddit import CommentConstants as CC
from util.constants.scraping import (
    DataSources as DS,
    ParentSources as PS,
    ScrapedPostConstants as SPC,
)
from util.constants.pubsub import Topics, BeamSubscriptions as BS
from util.constants.bigquery import Tables

from . import print_collection, publish_to_pubsub, standard_options

LOCATION = os.path.realpath(os.path.join(os.getcwd(), os.path.dirname(__file__)))


def parse_scraped_post(comment: Dict) -> Dict:
    """Parse scraped post objects from a comment object"""

    post = {
        SPC.PARENT_ID: comment[CC.ID],
        SPC.TIMESTAMP: comment[CC.CREATED_UTC],
        SPC.TEXT: comment[CC.BODY],
        SPC.DATA_SOURCE: DS.REDDIT,
        SPC.PARENT_SOURCE: PS.COMMENT_BODY,
    }

    return post


def extract_scraped_posts(comments: PCollection) -> PCollection:
    """Apply the parse_scraped_post method to a comments collection"""
    return comments | beam.Map(parse_scraped_post)


def create_test_comments_pipeline():
    """
    Create the local test pipeline that uses data
    saved in a json file inside the 'test' directory
    """

    # Load test data from json
    with open(f"{LOCATION}/test/comments.json", "r") as test_file:
        test_data = json.load(test_file)

    # Create pipeline
    pipeline_options = PipelineOptions(standard_options)

    with beam.Pipeline(options=pipeline_options) as pipeline:
        comments = pipeline | beam.Create(test_data)
        scraped_posts = extract_scraped_posts(comments)
        print_collection(comments)
        print_collection(scraped_posts)


def create_comments_pipeline():
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
        comments = (
            pipeline
            | beam.io.ReadFromPubSub(subscription=f"{BS.PREFIX}/{BS.REDDIT_COMMENTS}")
            | beam.Map(lambda d: d.decode("utf-8"))
            | beam.Map(json.loads)
        )

        # Write to bigquery
        _ = comments | beam.io.WriteToBigQuery(
            Tables.REDDIT_COMMENTS,
            write_disposition=beam.io.BigQueryDisposition.WRITE_APPEND,
        )

        # Extract scraped posts
        scraped_posts = extract_scraped_posts(comments)

        # Write to Pub/Sub
        publish_to_pubsub(scraped_posts, f"{Topics.PREFIX}/{Topics.SCRAPED_POSTS}")
