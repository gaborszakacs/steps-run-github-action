package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func loadEnvs(config actionConfig) error {
	envsMap, err := readFromStepInput()
	if err != nil {
		return err
	}

	for key, config := range config.Inputs {
		value, ok := lookupEnv(key, envsMap)
		if !ok {
			switch {
			case config.Default != "":
				value = config.Default
			case config.Required:
				return fmt.Errorf("input config '%s' is required by the action but is not set", key)
			default:
				continue
			}
		}

		if err := os.Setenv(convertEnvKeyToExpectedByAction(key), value); err != nil {
			return err
		}
	}
	return nil
}

func readFromStepInput() (map[string]string, error) {
	envsFromStepInput := os.Getenv("with")
	envsMap := make(map[string]string)
	if envsFromStepInput != "" {
		if err := yaml.Unmarshal([]byte(envsFromStepInput), &envsMap); err != nil {
			return nil, err
		}
	}
	return envsMap, nil
}

func convertEnvKeyToExpectedByAction(key string) string {
	return fmt.Sprintf("INPUT_%s", strings.ToUpper(key))
}

func lookupEnv(key string, envsMap map[string]string) (string, bool) {
	if value, ok := envsMap[key]; ok {
		return value, true
	}
	return os.LookupEnv(key)
}
