package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Config содержит настройки приложения.
type Config struct {
	Port           string   `yaml:"port"`
	DBUser         string   `yaml:"db_user"`
	DBPassword     string   `yaml:"db_password"`
	DBHost         string   `yaml:"db_host"`
	DBPort         string   `yaml:"db_port"`
	DBName         string   `yaml:"db_name"`
	FishingMethods []string `yaml:"fishing_methods"`
	FishSpecies    []string `yaml:"fish_species"`
	Baits          []string `yaml:"baits"`
}

// LoadConfig загружает настройки из файла config.yaml.
func LoadConfig() (*Config, error) {
	// Чтение содержимого файла config.yaml
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл config.yaml: %v", err)
	}

	// Декодирование YAML в структуру Config
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("не удалось декодировать config.yaml: %v", err)
	}

	// Проверка обязательных параметров
	if cfg.Port == "" {
		cfg.Port = "8080" // Значение по умолчанию
	}
	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("не указаны все необходимые параметры для подключения к базе данных")
	}

	return &cfg, nil
}

// GetDatabaseDSN формирует строку подключения к базе данных.
func (cfg *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
}
