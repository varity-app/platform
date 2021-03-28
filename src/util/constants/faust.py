"""
Constants related to Faust
"""

import os


class Constants:
    """General Faust App configuration constants"""

    APP_NAME = "varity-faust-app"

    DEBUG = os.environ.get("FAUST_DEBUG", default=False)
    STORE_PATH = "/tmp/faust"
