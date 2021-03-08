import os

class Config:
    BOOTSTRAP_SERVERS = os.environ.get("BOOTSTRAP_SERVERS")
    SASL_MECHANISM = 'PLAIN'
    SECURITY_PROTOCOL = 'SASL_SSL'
    SASL_USERNAME = os.environ.get("SASL_USERNAME")
    SASL_PASSWORD = os.environ.get("SASL_PASSWORD")

    OBJ = {
        'boostrap.servers': BOOTSTRAP_SERVERS,
        'sasl.mechanism': SASL_MECHANISM,
        'security.protocol': SECURITY_PROTOCOL,
        'sasl.username': SASL_USERNAME,
        'sasl.password': SASL_PASSWORD,
    }

class Topics:
    REDDIT_SUBMISSIONS = 'reddit-submissions'