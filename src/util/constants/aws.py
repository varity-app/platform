class Misc:
    REGION = 'us-east-2'


class SNS:
    ARN_PREFIX = f'arn:aws:sns:{Misc.REGION}:178852309825'

    REDDIT_COMMENTS = f'{ARN_PREFIX}:reddit-comments'
    REDDIT_SUBMISSIONS = f'{ARN_PREFIX}:reddit-submissions'
    SCRAPED_POSTS = f'{ARN_PREFIX}:scraped-posts'
    TICKER_MENTIONS = f'{ARN_PREFIX}:ticker-mentions'

    JSON = 'json'
    DEFAULT = 'default'


class S3:
    PREFIX = "s3://"

    REDDIT_SUBMISSIONS = "varity-reddit-submissions"
    REDDIT_COMMENTS = "varity-reddit-comments"
    POST_SENTIMENT = "varity-post-sentiment"
    TICKER_MENTIONS = "varity-ticker-mentions"