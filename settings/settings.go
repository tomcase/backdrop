package settings

import (
	"backdropGo/db"
	"context"

	"github.com/google/uuid"
)

// GetClientID returns the Client ID used to identify the client
func GetClientID(ctx context.Context, sg db.SettingsReader) (string, error) {
	return getSetting(ctx, sg, "CLIENT_ID")
}

// GetDeviceID returns the Device ID use to identify the device
func GetDeviceID(ctx context.Context, sg db.SettingsReader) (string, error) {
	return getSetting(ctx, sg, "DEVICE_ID")
}

// CreateDeviceID creates and inserts a new DeviceID
func CreateDeviceID(ctx context.Context, sw db.SettingsWriter) (string, error) {
	newID := uuid.New()
	return updateSetting(ctx, sw, "DEVICE_ID", newID.String())
}

// GetOutputDirectory returns the Device ID use to identify the device
func GetOutputDirectory(ctx context.Context, sg db.SettingsReader) (string, error) {
	return getSetting(ctx, sg, "OUTPUT_DIR")
}

func getSetting(ctx context.Context, sg db.SettingsReader, name string) (string, error) {
	value, err := sg.Get(ctx, name)
	if err != nil {
		return "", err
	}
	return value, nil
}

func updateSetting(ctx context.Context, sw db.SettingsWriter, name string, value string) (string, error) {
	value, err := sw.Update(ctx, name, value)
	if err != nil {
		return "", err
	}
	return value, nil
}
