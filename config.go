package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "stanwielga",
		Password: "",
		Name:     "galleria_dev",
	}
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name,
		)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACkey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "IamAsuperSecretString",
		HMACkey:  "secret-hmac-key",
		Database: DefaultPostgresConfig(),
	}
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func LoadConfig(configRequired bool) Config {
	f, err := os.Open(".config")
	if err != nil {
		if configRequired {
			panic(err)
		}

		fmt.Println("Using default config...")

		return DefaultConfig()
	}

	var c Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully loaded config")

	return c
}
