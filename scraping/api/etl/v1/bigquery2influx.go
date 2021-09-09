package etlv1

// BigqueryToInfluxSpec defines a request to pull data from Bigquery, aggregate,
// and store in InfluxDB.  This should be delivered via the body of a PubSub message.
type BigqueryToInfluxSpec struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
	Hour  int `json:"hour"`
}

// BigqueryToInfluxResponse is the response payload for a BigquerytoInflux service.
type BigqueryToInfluxResponse struct {
	PointsCount int `json:"points_count"`
}

// PubSubRequest is the payload of a Pub/Sub event triggered in Cloud Run.
type PubSubRequest struct {
	Message struct {
		Data string `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}
