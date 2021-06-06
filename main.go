package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("step failed: %v", err)
	}

	os.Exit(0)
}

func run() error {
	actionID := os.Getenv("uses")
	if actionID == "" {
		return fmt.Errorf("uses must be set")
	}

	path, cleanupFunc, err := cloneAction(actionID)
	defer cleanupFunc()
	if err != nil {
		return err
	}

	config, err := readActionConfig(path)
	if err != nil {
		return err
	}

	if err := loadEnvs(config); err != nil {
		return err
	}

	if err := runAction(path, config); err != nil {
		return err
	}

	return nil
}
