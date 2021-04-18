"""
Module for declaring apache beam pipelines
"""

from typing import Dict
from datetime import datetime, timedelta
import json
import uuid

import apache_beam as beam
from apache_beam.pvalue import PCollection

standard_options = [
    "--runner=DirectRunner",
    "--direct_running_mode=multi_processing",
]


def print_collection(collection: PCollection) -> None:
    """Nicely print out a PCollection of dict objects"""
    collection = collection | str(uuid.uuid1()) >> beam.Map(
        lambda x: print(json.dumps(x, indent=4))
    )


def publish_to_pubsub(collection: PCollection, topic: str) -> None:
    """Helper method for publishing to a Pub/Sub topic"""
    collection = (
        collection
        | beam.Map(lambda x: json.dumps(x).encode("utf-8"))
        | beam.io.WriteToPubSub(topic=topic)
    )


def check_age(obj: Dict, key: str, years=5) -> bool:
    """Check if the date field of a dictionary is too old to be inserted into BigQuery"""

    value = obj.get(key, datetime.now())

    if isinstance(value, str):
        value = datetime.fromisoformat(value)

    days_delta = 365 * years - 1  # Add a buffer of one day

    valid_age = value > datetime.now() - timedelta(days=days_delta)

    return valid_age
