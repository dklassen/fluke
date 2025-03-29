package main

import (
	"io/ioutil"
	"os"
)

func ReadFile(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SaveMigration(filepath string, data string) error {
	return os.WriteFile(filepath, []byte(data), 0644)
}
