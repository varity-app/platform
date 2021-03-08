import logging
import csv
import boto3
from datetime import datetime

from messages.sns import SNSMessage
from messages.reddit.submission import SubmissionMessage
from consumers.reddit.submission import SubmissionConsumer
from util.constants.aws import S3

logger = logging.getLogger()
logger.setLevel(logging.INFO)

consumer = SubmissionConsumer()

file_name = f'{datetime.now().strftime("%Y-%m-%dT%H-%M-%S")}-submissions.csv'
file_path = f'/tmp/{file_name}'
s3_client = boto3.client('s3')


def lambda_handler(event, context):
    """
    :param event: The event dict that contains the parameters sent when the function
                  is invoked.
    :param context: The context in which the function is called.
    :return: The result of the specified action.
    """
    logger.info('Event: %s', event)

    with open(file_path, 'w', encoding='UTF8') as f:
        writer = csv.DictWriter(f, fieldnames=SubmissionMessage.fields)
        writer.writeheader()

        for record in event['Records']:
            # Parse submission from message
            sub_obj = SNSMessage.parse_message(record)
            submission = SubmissionMessage.from_obj(sub_obj)

            # Parse tickers and post text
            has_tickers = consumer.run(submission)

            # Write to CSV
            if has_tickers:
                writer.writerow(submission.serialize(as_str=False))

    # Publish to S3
    s3_client.upload_file(file_path, S3.REDDIT_SUBMISSIONS, file_name)

    response = 1
    return response
