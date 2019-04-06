package config

/**
 * This package implements a simple configuration system which reads from a JSON
 * configuration file stored in the root directory of the project
 */

import (
    "log"
    "os"
    "encoding/json"
)

type Config struct {
    DBHost string
    DBPort string
    DBUser string
    DBPass string
    DBName string
}

func New(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        log.Printf("failed to open config file: <%s>", path)
        return nil, err
    }
    defer file.Close()

    var config Config

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        log.Println("failed to parse JSON config file")
        return nil, err
    }

    return &config, nil
}
