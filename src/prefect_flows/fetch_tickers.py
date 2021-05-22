"""
Prefect flow for fetching the list of tickers from IEXCloud
"""

from typing import Dict, List, Tuple
from os import environ as env
import re
import requests

import psycopg2
from psycopg2.extras import execute_values

from prefect import task, Flow
from prefect.storage import GCS
from prefect.tasks.secrets import PrefectSecret
from prefect.run_configs import KubernetesRun


IEX_BASE = "https://cloud.iexapis.com/stable"
TICKER_COLUMNS = [
    "symbol",
    "exchange",
    "exchange_name",
    "name",
    "scraped_on",
    "is_enabled",
    "type",
    "region",
    "currency",
    "cik",
    "short_name",
]

USERNAME = env.get("POSTGRES_USERNAME")
HOST = env.get("POSTGRES_HOST")
DB = env.get("POSTGRES_DB")
DEPLOYMENT = env.get("DEPLOYMENT", "dev")


assert USERNAME is not None
assert HOST is not None
assert DB is not None
assert DEPLOYMENT in ["dev", "prod"]


@task
def extract_tickers(iex_token: str) -> List[Dict[str, str]]:
    """Fetch the list of active tickers from IEXCloud"""

    response = requests.get(
        f"{IEX_BASE}/ref-data/symbols", params=dict(token=iex_token)
    )
    data = response.json()

    return data


@task
def transform_tickers(tickers: List[Dict[str, str]]) -> List[Dict[str, str]]:
    """Transform the tickers from the format provided by IEXCloud to the format stored in PostgreSQL"""

    new_tickers = []

    for ticker in tickers:
        new_ticker = dict(
            symbol=ticker["symbol"],
            exchange=ticker["exchange"],
            exchange_name=ticker["exchangeName"],
            name=ticker["name"],
            scraped_on=ticker["date"],
            is_enabled=ticker["isEnabled"],
            type=ticker["type"],
            region=ticker["region"],
            currency=ticker["currency"],
            cik=ticker["cik"],
        )

        new_tickers.append(new_ticker)

    return new_tickers


@task
def parse_postgres_tuples(
    columns: List[str], tickers: List[Dict[str, str]]
) -> List[Tuple]:
    """
    Parse the list of tuples in dictionary form into a
    list of tuples with values in the order required by PostgreSQL
    """

    new_tickers = [tuple(ticker[column] for column in columns) for ticker in tickers]

    return new_tickers


@task
def add_short_name(tickers: List[Dict[str, str]]) -> List[Dict[str, str]]:
    """Add the shortname column to the list of tickers"""

    # Get list of all words in the tickers' names
    names = [ticker["name"] for ticker in tickers]
    names = " ".join(names)
    words = names.split(" ")

    # Regex expressions
    name_rex = re.compile("")

    words = [word for word in words if word.isalpha()]
    words = [word for word in words if bool(name_rex.match(word))]

    # Find unique words
    word_counts = dict()
    for word in words:
        count = word_counts.get(word, 0) + 1
        word_counts[word] = count

    # Only select names that come up very infrequently
    unique_words = [word for word in words if word_counts[word] <= 10]

    # Set short_name field
    for ticker in tickers:
        name_words = ticker["name"].split(" ")
        short_name = None

        for word in name_words:
            if word in unique_words:
                short_name = word
                break

        ticker["short_name"] = short_name

    return tickers


def generate_query(columns: List[str]) -> str:
    """Generate a PostgreSQL query for inserting new rows into the Tickers column"""

    sql_query = f"""
        INSERT INTO tickers ({", ".join(columns)})
        VALUES %s
        ON CONFLICT (symbol) DO
            UPDATE SET name = EXCLUDED.name, is_enabled = EXCLUDED.is_enabled, scraped_on = EXCLUDED.scraped_on,
                short_name = EXCLUDED.short_name;
    """

    return sql_query


@task
def insert_query(
    sql_query: str,
    data: List[Tuple],
    username: str,
    database: str,
    host: str,
    password: str,
):
    """Run an insert query with psycopg2.  The provided Postgres insert methods are simply too slow."""

    # Create connection
    conn = psycopg2.connect(
        dbname=database, user=username, password=password, host=host
    )
    cur = conn.cursor()

    # Run query
    execute_values(cur, sql_query, data)
    conn.commit()

    # Close connection
    cur.close()
    conn.close()


bucket = GCS(bucket=f"varity-prefect-{DEPLOYMENT}")
run_config = KubernetesRun(
    env={"EXTRA_PIP_PACKAGES": "psycopg2-binary"},
    image="prefecthq/prefect:latest-python3.8",
)

with Flow("Get Tickers", run_config=run_config, storage=bucket) as flow:
    # Secrets
    token = PrefectSecret("IEX_TOKEN")
    postgres_password = PrefectSecret("POSTGRES_PASSWORD")

    # Fetch and transform tickers
    ticker_list = extract_tickers(iex_token=token)
    ticker_list = transform_tickers(ticker_list)
    ticker_list = add_short_name(ticker_list)
    ticker_list = parse_postgres_tuples(TICKER_COLUMNS, ticker_list)

    # Insert into postgres
    query = generate_query(TICKER_COLUMNS)
    insert_query(
        sql_query=query,
        data=ticker_list,
        username=USERNAME,
        password=postgres_password,
        database=DB,
        host=HOST,
    )

flow.run()
# flow.register(project_name=f"varity-{DEPLOYMENT}", labels=[])
