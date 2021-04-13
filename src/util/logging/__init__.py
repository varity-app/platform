"""
Helper module for a custom logger
"""

import logging


def set_log_config() -> None:
    """Sets the basic logging config"""
    logging.basicConfig(
        format="%(name)s:%(levelname)s:%(message)s",
        level=logging.INFO,
    )
