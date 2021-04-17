"""
Entrypoint script for picking and running a beam pipeline
"""

import argparse

from beam.submissions import (
    create_test_submissions_pipeline,
    create_submissions_pipeline,
)

from beam.comments import (
    create_test_comments_pipeline,
    create_comments_pipeline,
)

from beam.scraped_posts import (
    create_test_scraped_posts_pipeline,
    create_scraped_posts_pipeline,
)

SUBMISSIONS = "submissions"
COMMENTS = "comments"
SCRAPED_POSTS = "scraped_posts"
PIPELINES = [SUBMISSIONS, COMMENTS, SCRAPED_POSTS]


def parse_args():
    """Parse arguments from command line"""
    parser = argparse.ArgumentParser()

    parser.add_argument("pipeline", choices=PIPELINES)
    parser.add_argument("--debug", action="store_true")

    args = parser.parse_args()

    return args


def main():
    """Entrypoint method"""
    args = parse_args()

    if args.pipeline == SUBMISSIONS:
        if args.debug:
            create_test_submissions_pipeline()
        else:
            create_submissions_pipeline()

    elif args.pipeline == COMMENTS:
        if args.debug:
            create_test_comments_pipeline()
        else:
            create_comments_pipeline()

    elif args.pipeline == SCRAPED_POSTS:
        if args.debug:
            create_test_scraped_posts_pipeline()
        else:
            create_scraped_posts_pipeline()


if __name__ == "__main__":
    main()
