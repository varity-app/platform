"""
Prefect flow for fetching the list of tickers from IEXCloud
"""

from typing import Dict, List
import requests

from prefect import task, Flow
from prefect.client.secrets import Secret


IEX_BASE = "https://cloud.iexapis.com/stable"


@task
def print_value(x: Secret) -> None:
    print(x)


@task
def fetch_tickers(iex_token: Secret) -> List[Dict[str, str]]:
    """Fetch the list of active tickers from IEXCloud"""
    response = requests.get(
        f"{IEX_BASE}/ref-data/symbols", params=dict(token=iex_token.get())
    )
    data = response.json()

    return data


with Flow("Test Secret") as flow:
    token = Secret("IEX_TOKEN")
    tickers = fetch_tickers(iex_token=token)
    print_value(tickers)

flow.run()
