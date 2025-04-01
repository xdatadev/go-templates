package config

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/xdatadev/superapp-packages/superapp-common/configps"
)

// Singleton instance
var (
	config Config
	once   sync.Once
)

func LoadParameters() *Config {
	once.Do(func() {
		loadParameters()
	})

	return &config
}

func loadParameters() {

	accessKey := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
	secretKey := strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	region := strings.TrimSpace(os.Getenv("AWS_REGION"))

	if region == "" {
		region = "us-east-1"
	}

	var opts []configps.Option
	opts = append(opts, configps.WithRegion(region))
	opts = append(opts, configps.WithMaxRetries(3))

	if accessKey != "" && secretKey != "" {
		opts = append(opts, configps.WithCredentials(accessKey, secretKey))
	}

	ps, err := configps.New(opts...)

	if err != nil {
		panic(err)
	}

	err = ps.LoadSettings(context.Background(), "{{.Scaffold.ParametersRoot}}", &config)
	if err != nil {
		panic(err)
	}
}
