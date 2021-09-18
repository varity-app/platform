"""
SQLAlchemy model for storing Cohort Memberships, originating from a DBT model in BigQuery.
"""

from sqlalchemy import Column, String, Date, Boolean

from .base import Base


class CohortMembership(Base):
    """Cohort Membership ORM model"""

    __tablename__ = "cohort_memberships"

    author_id = Column(String, primary_key=True)
    subreddit = Column(String, primary_key=True)
    month = Column(Date, primary_key=True)
    cohort = Column(String, primary_key=True)