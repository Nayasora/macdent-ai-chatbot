package configs

type ApiServerConfig struct {
	Port     int
	Postgres *PostgresConfig
	Qdrant   *QdrantConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type QdrantConfig struct {
	Host   string
	Port   int
	ApiKey string
}

func NewConfig(env *Env) *ApiServerConfig {
	return &ApiServerConfig{
		Port: env.MustInt("APP_INTERNAL_PORT"),
		Postgres: &PostgresConfig{
			Host:     env.MustString("POSTGRES_HOST"),
			Database: env.MustString("POSTGRES_DATABASE"),
			Port:     env.MustInt("POSTGRES_PORT"),
			User:     env.MustString("POSTGRES_USER"),
			Password: env.MustString("POSTGRES_PASSWORD"),
		},
		Qdrant: &QdrantConfig{
			Host:   env.MustString("QDRANT_HOST"),
			Port:   env.MustInt("QDRANT_PORT"),
			ApiKey: env.MustString("QDRANT_API_KEY"),
		},
	}
}
