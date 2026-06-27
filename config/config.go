package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const settingsFileName = "code_settings.json"

type Config struct {
	JiraURL   string `json:"jiraURL"`
	GitlabURL string `json:"gitlabURL"`
}

func NewConfig() (*Config, error) {
	config := &Config{}
	if err := config.loadSettingsFile(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) loadSettingsFile() error {
	settingsPath, err := settingsFilePath()
	if err != nil {
		return err
	}

	file, err := os.Open(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open settings file %s: %w", settingsPath, err)
	}
	defer file.Close()

	settings := make(map[string]string)
	if err := json.NewDecoder(file).Decode(&settings); err != nil {
		return fmt.Errorf("failed to decode settings file %s: %w", settingsPath, err)
	}

	return c.applySettings(settings)
}

func settingsFilePath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Join(filepath.Dir(executablePath), settingsFileName), nil
}

func (c *Config) applySettings(settings map[string]string) error {
	value := reflect.ValueOf(c).Elem()
	valueType := value.Type()
	fields := make(map[string]reflect.Value, value.NumField()*3)

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := valueType.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		fields[fieldType.Name] = fieldValue
		fields[strings.ToLower(fieldType.Name)] = fieldValue

		jsonName := strings.Split(fieldType.Tag.Get("json"), ",")[0]
		if jsonName != "" && jsonName != "-" {
			fields[jsonName] = fieldValue
		}
	}

	for key, settingValue := range settings {
		fieldValue, ok := fields[key]
		if !ok {
			return fmt.Errorf("unknown config key %q in %s", key, settingsFileName)
		}
		if fieldValue.Kind() != reflect.String {
			return fmt.Errorf("config key %q must map to a string field", key)
		}
		fieldValue.SetString(settingValue)
	}

	return nil
}
