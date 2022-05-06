package config

type Config struct {
	Bind string
	Log  struct {
		Level string
	}
}

var C Config
