from confluent_kafka import Producer

from util.contants.kafka import Config

def main():
    producer = Producer(Config.OBJ)