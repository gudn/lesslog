package config

type Config struct {
	Bind      string
	Mode      string
	Pg        string
	Messaging string
	Log       struct {
		Level string
	}
}

var C Config
