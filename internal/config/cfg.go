package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/playwright-community/playwright-go"
)

// Структура для хранения настроек из окружения или установленных по умолчанию
type Config struct {
	Port int `env:"SERVER_PORT" envDefault:"13005"`

	PgPort   string `env:"PG_PORT" envDefault:"5432"`
	PgHost   string `env:"PG_HOST" envDefault:"127.0.0.1"`
	PgDBName string `env:"PG_DB_NAME" envDefault:"db"`
	PgUser   string `env:"PG_USER" envDefault:"konstantin"`
	PgPwd    string `env:"PG_PWD" envDefault:"konstantin"`
}

// Структура для хранения объекта Config и обеспечения доступа к его значениям
type ConfigProvider struct {
	config *Config
}

// Функция для создания объекта ConfigProvider и инициализации Config из окружения или установки значений по умолчанию
func NewConfigProvider() (*ConfigProvider, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return &ConfigProvider{
		config: cfg,
	}, nil
}

// Функция для получения объекта Config
func (cp *ConfigProvider) GetConfig() *Config {
	return cp.config
}

var PW *playwright.Playwright

// Функция для инициализации Playwright
func Init() {
	pw, err := playwright.Run()
	if err != nil {
		panic(err) // Можно обработать ошибку по-другому, но пока просто завершим программу
	}
	PW = pw
}
