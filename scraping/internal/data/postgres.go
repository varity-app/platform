package data

import (
	"fmt"
	"net/url"
)

// PostgresOpts stores configuration parameters for connecting to a Postgres database.
type PostgresOpts struct {
	Username string
	Password string
	Database string
	Address  string
	Network  string
}

func ParsePostgresDSN(opts PostgresOpts) string {
	// Escape fields which may have unfriendly characters
	escapedUsername := url.QueryEscape(opts.Username)
	escapedPassword := url.QueryEscape(opts.Password)

	// Create DSN
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		escapedUsername, escapedPassword,
		opts.Address,
		opts.Database,
	)

	if opts.Network == "unix" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s",
			opts.Address,
			opts.Username,
			opts.Password,
			opts.Database,
		)
	}

	return dsn
}
