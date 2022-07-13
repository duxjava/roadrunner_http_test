package roadrunner_http_test

import "time"

type Config struct {
	MysqlConnection string        `mysql:"connection"`
	MysqlLifetime   time.Duration `mysql:"lifetime"`
	MysqlMaxidle    int           `mysql:"maxidle"`
	MysqlMaxopen    int           `mysql:"maxopen"`
}

func (c *Config) InitDefaults() {
}
