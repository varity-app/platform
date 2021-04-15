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

    # Remove urls
    text = re.sub(r'https?:\/\/.*[\r\n]*', '', text)

    # Remove anything not alphanumeric
    text = re.sub(r"[^a-zA-Z0-9 \n\.]", " ", text)

    # Find tickers
    tickers = re.findall(r"[A-Z][A-Z0-9.]{0,3}[A-Z]", text)

    # Cross reference with local ticker list
    if ticker_list:
        tickers = [ticker for ticker in tickers if ticker in ticker_list]

    # Cross reference with ticker black list
    tickers = [ticker for ticker in tickers if ticker not in black_list]

    # Remove duplicates
    tickers = list(set(tickers))

    # If there are more than 5 tickers, count as spam
    if len(tickers) > 5:
        tickers = []

    return tickers


with open(os.path.join(LOCATION, "blacklist.txt"), "r") as f:
    black_list = [x.strip() for x in f.readlines()]

with open(os.path.join(LOCATION, "tickers.txt"), "r") as f:
    all_tickers = [x.strip() for x in f.readlines()]
