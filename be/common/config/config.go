package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/cmd"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/postgres"
)

type (
	Otel struct {
		Enabled bool
	}
	GrpcConfig struct {
		Address string
	}

	RestConfig struct {
		Address string
	}

	JwtConfig struct {
		CrtFile string
		KeyFile string
		UseMock string
		Ttl     struct {
			AccessToken  time.Duration
			RefreshToken time.Duration
		}
	}

	AppConfig struct {
		Otel            Otel
		Postgres        postgres.Config
		Aws             awsclient.AwsConfig
		Dynamodb        awsclient.DynamodbConfig
		Grpc            GrpcConfig
		Jwt             JwtConfig
		Rest            RestConfig
		Mgo             mgo.Config
		ShutdownTimeout time.Duration `default:"20s"`
	}
)

func InitConfig(cfg any) (err error) {
	for _, p := range cmd.GlobalFlags.ConfigPaths {
		viper.AddConfigPath(p)
	}
	viper.SetConfigName(cmd.GlobalFlags.ConfigName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	err = viper.Unmarshal(&cfg)
	return err
}
