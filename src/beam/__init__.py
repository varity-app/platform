"""
Module for declaring apache beam pipelines
"""

import json

import apache_beam as beam
from apache_beam.pvalue import PCollection

standard_options = [
    "--runner=DirectRunner",
    "--direct_running_mode=multi_processing",
]


def print_collection(collection: PCollection) -> None:
    """Nicely print out a PCollection of dict objects"""
    collection = collection | beam.Map(json.dumps, indent=4) | beam.Map(print)


def publish_to_pubsub(collection: PCollection, topic: str) -> None:
    """Helper method for publishing to a Pub/Sub topic"""
    collection = (
            collection
            | beam.Map(lambda x: json.dumps(x).encode("utf-8"))
            | beam.io.WriteToPubSub(topic=topic)
        )
