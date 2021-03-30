"""
Declare the helper class for querying DynamoDB
"""

from datetime import datetime, timedelta

import boto3
from boto3.dynamodb.conditions import Key

from util.constants.dynamo import Constants


class DynamoTable:
    """
    Helper class for querying a DynamoDB table
    """

    def __init__(self, table_name: str, primary_key: str, disable_ttl=False) -> None:
        self.dynamo = boto3.resource("dynamodb")
        self.table = self.dynamo.Table(table_name)

        self.primary_key = primary_key
        self.disable_ttl = disable_ttl

    def exists(self, key: str) -> bool:
        """
        Check if a key exists in the table
        """
        response = self.table.query(
            KeyConditionExpression=Key(self.primary_key).eq(key)
        )
        key_exists = response["Count"] > 0

        return key_exists

    def put(self, key: str, data=None, expiration_date=None) -> None:
        """
        Save an item to DynamoDB

        Arguments:
        key -- value for the primary key
        data -- a dict of other fields and values to store
        expiration_date -- a datetime object of when the item should expire
        """
        if expiration_date is None:
            expiration_date = datetime.now() + timedelta(
                days=Constants.DEFAULT_EXPIRATION_DAYS
            )
            expiration_date = int(expiration_date.timestamp())

        if data is None:
            data = dict()

        obj = dict(**data)
        obj[self.primary_key] = key

        if not self.disable_ttl:
            obj[Constants.EXPIRATION_DATE] = expiration_date

        self.table.put_item(Item=obj)
