package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config - структура для хранения конфигурации
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig - настройки сервера
type ServerConfig struct {
	Port string
}

// DatabaseConfig - настройки базы данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // имя файла (без расширения)
	viper.SetConfigType("yaml")   // тип файла
	viper.AddConfigPath("internal/config") // путь к файлу

	// Читаем конфиг
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
		return nil, err
	}

	// Создаем объект конфигурации
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Ошибка парсинга конфига: %v", err)
		return nil, err
	}

	return &cfg, nil
}
