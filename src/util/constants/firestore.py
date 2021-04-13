"""
Constants related to Firestore
"""

import os

from . import DEPLOYMENT


DISABLE_FIRESTORE = os.environ.get("DISABLE_FIRESTORE", False)
PROJECT = "varity"


class Collections:
    """
    Names of collections in Firestore
    """

    REDDIT_SUBMISSIONS = f"reddit-submissions-{DEPLOYMENT}"
    REDDIT_COMMENTS = f"reddit-comments-{DEPLOYMENT}"
