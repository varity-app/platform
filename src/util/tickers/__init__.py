"""
Helper methods for extracting stock tickers from text
"""

import re
import os
import string

LOCATION = os.path.realpath(os.path.join(os.getcwd(), os.path.dirname(__file__)))


def parse_tickers(text, ticker_list=None):
    """Parse stock tickers from a string"""

    if not ticker_list:
        ticker_list = all_tickers

    text = re.sub(r"[^a-zA-Z0-9 \n\.]", " ", text)

    tickers = re.findall(r"[A-Z][A-Z0-9.]{0,3}[A-Z]", text)

    if ticker_list:
        tickers = [ticker for ticker in tickers if ticker in ticker_list]

    tickers = [ticker for ticker in tickers if ticker not in black_list]

    # Remove duplicates
    tickers = list(set(tickers))

    return tickers


with open(os.path.join(LOCATION, "blacklist.txt"), "r") as f:
    black_list = f.read()

with open(os.path.join(LOCATION, "tickers.txt"), "r") as f:
    all_tickers = f.read()
