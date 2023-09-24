package db

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx"
)

// DSNToConnConfig converts a DSN string to pgx.ConnConfig.
func DSNToConnConfig(dsn string) (pgx.ConnConfig, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return pgx.ConnConfig{}, err
	}

	port, err := strconv.Atoi(parsedURL.Port())
	if err != nil {
		return pgx.ConnConfig{}, err
	}

	password, _ := parsedURL.User.Password()

	return pgx.ConnConfig{
		Host:                 parsedURL.Hostname(),
		Port:                 uint16(port),
		User:                 parsedURL.User.Username(),
		Password:             password,
		Database:             strings.TrimLeft(parsedURL.Path, "/"),
		PreferSimpleProtocol: true,
	}, nil
}

// ConnConfigToDSN converts a pgx.ConnConfig to a DSN string.
func ConnConfigToDSN(config pgx.ConnConfig) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}
