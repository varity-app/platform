"""
SQLAlchemy model for Tickers and their metadata, scraped from IEXCloud.
More documentation found at: https://iexcloud.io/docs/api/#symbols
"""

from sqlalchemy import Column, String, Date, Boolean

from .base import Base


class Ticker(Base):
    """Ticker ORM model"""

    __tablename__ = "tickers"

    symbol = Column(String, primary_key=True)
    exchange = Column(String, primary_key=True)
    exchange_name = Column(String)
    name = Column(String)
    scraped_on = Column(Date)
    type = Column(String)
    region = Column(String)
    currency = Column(String)
    is_enabled = Column(Boolean)
    cik = Column(String)




