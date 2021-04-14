"""
Apache beam pipeline definitions for processing reddit submisisons
"""

import os
import json
from typing import List, Dict

import apache_beam as beam
from apache_beam.pvalue import PCollection
from apache_beam.options.pipeline_options import PipelineOptions

from util.constants.reddit import SubmissionConstants as SC
from util.constants.scraping import (
    DataSources as DS,
    ParentSources as PS,
    MentionTypes,
    TickerMentionsConstants as TMC,
    ScrapedPostConstants as SPC,
)
from util.constants.pubsub import Topics

from . import print_collection

LOCATION = os.path.realpath(os.path.join(
    os.getcwd(), os.path.dirname(__file__)))


def parse_scraped_posts(submission: Dict) -> List[Dict]:
    """Parse scraped post objects from a submission object"""

    scraped_posts = []

    for field, source in [
        (SC.TITLE, PS.SUBMISSION_TITLE),
        (SC.SELFTEXT, PS.SUBMISSION_SELFTEXT),
    ]:
        if not submission[field]:  # If post has no selftext (image post)
            continue

        post = {
            SPC.PARENT_ID: submission[SC.ID],
            SPC.TIMESTAMP: submission[SC.CREATED_UTC],
            SPC.TEXT: submission[field],
            SPC.DATA_SOURCE: DS.REDDIT,
            SPC.PARENT_SOURCE: source,
        }

        scraped_posts.append(post)

    return scraped_posts


def extract_scraped_posts(submissions: PCollection) -> PCollection:
    """Apply the parse_scraped_posts method to a submissions collection"""
    return (
        submissions
        | beam.FlatMap(parse_scraped_posts)
    )


def create_test_submissions_pipeline():
    """
    Create the local test pipeline that uses data
    saved in a json file inside the 'test' directory
    """

    # Load test data from json
    with open(f"{LOCATION}/test/submissions.json", "r") as f:
        test_data = json.load(f)

    # Create pipeline
    pipeline_options = PipelineOptions()

    with beam.Pipeline(options=pipeline_options) as pipeline:
        submissions = (
            pipeline
            | beam.Create(test_data)
        )
        submissions = extract_scraped_posts(submissions)
        print_collection(submissions)


def create_dataflow_submissions_pipeline():
    """
    Create the local test pipeline that uses data
    saved in a json file inside the 'test' directory
    """

    # Load test data from json
    with open(f"{LOCATION}/test/submissions.json", "r") as f:
        test_data = json.load(f)

    # Create pipeline
    pipeline_options = PipelineOptions()

    with beam.Pipeline(options=pipeline_options) as pipeline:
        submissions = (
            pipeline
            | beam.io.ReadFromPubSub(topic=f"{Topics.PREFIX}/{Topics.REDDIT_SUBMISSIONS}")
            | beam.Map(lambda d: d.decode('utf-8'))
            | beam.Map(json.loads)
        )
        submissions = extract_scraped_posts(submissions)
        print_collection(submissions)