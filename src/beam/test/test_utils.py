"""
Unit tests for the general utilities defined in beam.__init__.py
"""

from datetime import datetime
import pytest

from .. import check_age

test_data = [
    (datetime(2021, 1, 1).isoformat(), True),
    (datetime(2000, 1, 1).isoformat(), False),
]


@pytest.mark.parametrize("date,valid", test_data)
def test_check_age(date: str, valid: bool) -> None:
    """Test the check_age utility method"""

    obj = dict(timestamp=date)
    test_valid = check_age(obj, "timestamp")

    assert test_valid == valid
