"""
Unit tests for ticker parsing
"""

import pytest

from . import parse_tickers

test_data = [
    ("Buy lots and lots of GME! Stonks only go up!", ["GME"]),
    ("How I made money off of $BABA, a success store that I DID MYSELF", ["BABA"]),
    ("THIS POST IS IN ALL CAPS WE DID IT BOYS GO DO AAPL", []),
    ("Apple is going to get bigger $AAPL $GOOGL $BABA $MSFT $TSLA $BB", []),
    ("I think BB is going to be bigger than GME", ["BB", "GME"]),
]


@pytest.mark.parametrize("text,true_tickers", test_data)
def test_ticker_parsing(text, true_tickers):
    """Unit test parse_tickers method"""
    tickers = parse_tickers(text)

    assert set(tickers) == set(true_tickers)
