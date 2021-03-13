package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// SettingsReader retrieves settings from the DB
type SettingsReader interface {
	Get(ctx context.Context, settingName string) (string, error)
}

// SettingsWriter updates a setting and returns its value
type SettingsWriter interface {
	Update(ctx context.Context, settingName string, value string) (string, error)
}

// SettingsReaderWriter can both read and write to settings
type SettingsReaderWriter interface {
	SettingsReader
	SettingsWriter
}

// Get will retrieve the setting from the DB
func (*DB) Get(ctx context.Context, settingName string) (string, error) {
	dbpool, err := connect(ctx)
	if err != nil {
		return "", err
	}
	defer dbpool.Close()

	var settingValue string
	err = dbpool.QueryRow(ctx, "SELECT value FROM settings WHERE name=$1", settingName).Scan(&settingValue)
	if err != nil {
		return "", err
	}
	return settingValue, nil
}

// Update will update a setting in the DB
func (*DB) Update(ctx context.Context, settingName string, value string) (string, error) {
	dbpool, err := connect(ctx)
	if err != nil {
		return "", err
	}
	defer dbpool.Close()

	var settingValue string
	err = dbpool.QueryRow(ctx, "UPDATE settings SET value=$1 WHERE name=$2 RETURNING value", value, settingName).Scan(&settingValue)
	if err != nil {
		return "", err
	}
	return settingValue, nil
}

func connect(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, os.Getenv("POSTGRES_URL"))
}
