package config

type MailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

func GetMailConfig() *MailConfig {
	return &MailConfig{
		Host:     getEnv("MAIL_HOST", "smtp.gmail.com"),
		Port:     getEnv("MAIL_PORT", "587"),
		User:     getEnv("MAIL_USER", ""),
		Password: getEnv("MAIL_PASSWORD", ""),
		From:     getEnv("MAIL_FROM", ""),
	}
}
