import faust
from aiokafka.helpers import create_ssl_context

from util.constants.kafka import Config as KafkaConfig
from util.constants.faust import Constants

broker_credentials = faust.SASLCredentials(
    username=KafkaConfig.SASL_USERNAME,
    password=KafkaConfig.SASL_PASSWORD,
    mechanism=faust.types.auth.SASLMechanism.PLAIN,
    ssl_context=create_ssl_context(),
)

broker_credentials.protocol = faust.types.auth.AuthProtocol.SASL_SSL

app = faust.App(
    Constants.APP_NAME,
    broker=f"kafka://{KafkaConfig.BOOTSTRAP_SERVERS}",
    broker_credentials=broker_credentials,
    autodiscover=True,
    origin="process"
)
