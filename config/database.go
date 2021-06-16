package config

type DatabaseConfig struct {
	User     string
	Password string
	Port     int
	DBName   string
	Host     string
	DBLog    int
}

func getDBConfig() *DatabaseConfig {
	return &DatabaseConfig{
		User:     getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", "secret"),
		Port:     getEnvAsInt("DB_PORT", 3307),
		Host:     getEnv("DB_HOST", "localhost"),
		DBName:   getEnv("DB_NAME", "release_tracker"),
		DBLog:    getEnvAsInt("ENABLE_DB_LOG", 0),
	}
}
