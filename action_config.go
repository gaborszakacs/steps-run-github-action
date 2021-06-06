package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type actionConfig struct {
	Name        string                  `yaml:"name"`
	Description string                  `yaml:"description"`
	Inputs      map[string]inputConfig  `yaml:"inputs"`
	Outputs     map[string]outputConfig `yaml:"outputs"`
	Runs        runConfig               `yaml:"runs"`
}

type inputConfig struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
}

type outputConfig struct {
	Description string `yaml:"description"`
}

type runConfig struct {
	Using string `yaml:"using"`
	Main  string `yaml:"main"`
}

func readActionConfig(repoPath string) (actionConfig, error) {
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s/action.yml", repoPath))
	if err != nil {
		return actionConfig{}, err
	}

	var config actionConfig
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return actionConfig{}, err
	}
	return config, nil
}
