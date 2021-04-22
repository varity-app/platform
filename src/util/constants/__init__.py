"""
Module for storing constants and preventing magic variables
"""

import os


class DeploymentTypes:
    """Allowed deployment types"""

    DEV = "dev"
    TEST = "test"
    PROD = "prod"


DEPLOYMENT = os.environ.get("DEPLOYMENT", DeploymentTypes.DEV)
assert DEPLOYMENT in [DeploymentTypes.DEV, DeploymentTypes.TEST, DeploymentTypes.PROD]

RELEASE = os.environ.get("RELEASE")
assert RELEASE is not None
