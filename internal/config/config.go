package config

type Config struct {
	Bind string
	Mode string
	Pg   string
	Log  struct {
		Level string
	}
}

var C Config