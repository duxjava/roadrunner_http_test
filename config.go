package roadrunner_http_test

import "time"

type Config struct {
	Mysql Mysql `mapstructure:"mysql"`
}

type Mysql struct {
	Connection string        `mapstructure:"connection"`
	Lifetime   time.Duration `mapstructure:"lifetime"`
	Maxidle    int           `mapstructure:"maxidle"`
	Maxopen    int           `mapstructure:"maxopen"`
}
