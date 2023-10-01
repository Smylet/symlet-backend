package db

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
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

	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		parsedURL.Hostname(),
		port,
		parsedURL.User.Username(),
		password,
		strings.TrimLeft(parsedURL.Path, "/"),
	)

	config, err := pgx.ParseConfig(connString)
	if err != nil {
		return pgx.ConnConfig{}, err
	}
	return *config, nil
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
