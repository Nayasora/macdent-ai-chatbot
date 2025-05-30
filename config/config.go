package config

// Config - конфигурация приложения из обязательных переменных окружения
type Config struct {
	Qdrant *QdrantConfig
	App    *AppConfig
}

// QdrantConfig - настройки подключения к Qdrant
type QdrantConfig struct {
	Host string
	Port int
}

// AppConfig - общие настройки приложения
type AppConfig struct {
	Environment string
	LogLevel    string
	Debug       bool
}

// NewConfig создает конфигурацию из обязательных переменных .env файла
func NewConfig() *Config {
	env := NewEnv(".env")

	return &Config{
		Qdrant: &QdrantConfig{
			Host: env.MustString("QDRANT_HOST"),
			Port: env.MustInt("QDRANT_PORT"),
		},
		App: &AppConfig{
			Environment: env.MustString("APP_ENV"),
			LogLevel:    env.MustString("LOG_LEVEL"),
			Debug:       env.MustBool("DEBUG"),
		},
	}
}
