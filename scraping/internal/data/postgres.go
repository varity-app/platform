package data

// PostgresOpts stores configuration parameters for connecting to a Postgres database.
type PostgresOpts struct {
	Username string
	Password string
	Database string
	Address  string
}
