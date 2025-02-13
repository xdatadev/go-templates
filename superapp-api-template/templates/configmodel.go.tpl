package config

type Config struct {
	DB           DBConfig           `param:"database"`
	Service      any                `param:"service"`
}

type DBConfig struct {
	Host              string `param:"host"`
	User              string `param:"user"`
	DbName            string `param:"dbname"`
	Password          string `param:"password"`
	Region            string `param:"region"`
	Port              int32  `param:"port"`
	MaxPoolSize       int32  `param:"maxpoolsize"`
	MinPoolSize       int32  `param:"minpoolsize"`
	MaxIdleTime       time.Duration
	MaxLifeTime       time.Duration
	HealthCheckPeriod time.Duration
}