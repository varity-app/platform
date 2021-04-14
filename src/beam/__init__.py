"""
Module for declaring apache beam pipelines
"""

import json

import apache_beam as beam
from apache_beam.pvalue import PCollection


def print_collection(collection: PCollection) -> None:
    """Nicely print out a PCollection of dict objects"""
    collection = (
        collection
        | beam.Map(json.dumps, indent=4)
        | beam.Map(print)
    )
