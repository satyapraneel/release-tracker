package config

type DatabaseConfig struct {
	User     string
	Password string
	Port     int
	DBName   string
	Host     string
}

func getDBConfig() *DatabaseConfig {
	return &DatabaseConfig{
		User:     getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Port:     getEnvAsInt("DB_PORT", 3306),
		Host:     getEnv("DB_HOST", "localhost"),
		DBName:   getEnv("DB_NAME", ""),
	}
}
