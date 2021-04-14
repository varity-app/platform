"""
Entrypoint for running a reddit submissions apache beam pipeline
"""

from beam.submissions import (
    create_test_submissions_pipeline,
    create_dataflow_submissions_pipeline,
)


def main():
    """Entrypoint method"""
    # create_test_submissions_pipeline()
    create_dataflow_submissions_pipeline()


if __name__ == "__main__":
    main()
