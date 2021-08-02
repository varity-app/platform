package main

import "time"

const WAIT_INTERVAL time.Duration = 2 * time.Second

// Firestore collection for checkpointing kafka offsets
const FIRESTORE_KAFKA_OFFSETS = "kafka-offsets"

// Kafka timeouts
const CONSUMER_TIMEOUT time.Duration = 2 * time.Second
const PRODUCER_TIMEOUT_MS int = 1000 // 1 second

var TICKER_BLACKLIST []string = []string{
	"WSB",
	"USA",
	"TENDIES",
	"FUCK",
	"YOLO",
	"UK",
	"DD",
	"WTF",
	"IPO",
	"FV",
	"TD",
	"CEO",
	"IM",
	"IMO",
	"AF",
	"OK",
	"ATM",
	"BS",
	"EU",
	"US",
	"NY",
	"NYSE",
	"NYC",
	"GOP",
	"GG",
	"HK",
	"AWS",
	"ETF",
	"GTA",
	"EK",
	"RH",
	"TDA",
	"ITM",
	"OTM",
	"IG",
	"TV",
	"USD",
	"EUR",
	"LY",
	"NG",
	"EV",
	"QA",
	"OG",
	"AMA",
	"LMAO",
	"FD",
	"PC",
	"FCF",
	"DCF",
	"PE",
	"PEG",
	"IN",
	"OR",
	"TO",
	"AND",
	"UP",
	"SO",
	"WE",
	"FOR",
	"HODL",
	"YOU",
	"IQ",
	"GOOD",
	"FR",
	"XP",
	"SIX",
	"NOW",
	"TM",
	"TA",
	"HA",
}
