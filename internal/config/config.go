package config

type Config struct {
	Bind      string
	Mode      string
	Pg        string
	Nats      string
	Messaging string
	MonBind   string
	Log       struct {
		Level string
	}
}

var C Config
