"""
Entrypoint for running a reddit submissions apache beam pipeline
"""

from beam.comments import (
    create_test_comments_pipeline,
    create_dataflow_comments_pipeline,
)


def main():
    """Entrypoint method"""

    create_test_comments_pipeline()
    # create_dataflow_comments_pipeline()


if __name__ == "__main__":
    main()
