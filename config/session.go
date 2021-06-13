package config

type Session struct {
	Secret string
}

func SessionDetails() *Session {
	return &Session{
		Secret: getEnv("SESSION_SECRET", ""),
	}
}
