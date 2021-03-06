class Misc:
    REGION = 'us-east-2'


class SNS:
    ARN_PREFIX = f'arn:aws:sns:{Misc.REGION}:178852309825'

    REDDIT_COMMENTS = f'{ARN_PREFIX}:reddit-comments.fifo'
    REDDIT_SUBMISSIONS = f'{ARN_PREFIX}:reddit-submissions.fifo'
    SCRAPED_POSTS = f'{ARN_PREFIX}:scraped-posts.fifo'
    TICKER_MENTIONS = f'{ARN_PREFIX}:ticker-mentions.fifo'

    JSON = 'json'
    DEFAULT = 'default'
