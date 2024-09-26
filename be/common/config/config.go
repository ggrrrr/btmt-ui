package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/ggrrrr/btmt-ui/be/common/awsdb"
	"github.com/ggrrrr/btmt-ui/be/common/cmd"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/postgres"
)

type (
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
		TTL     time.Duration
	}

	AppConfig struct {
		Postgres        postgres.Config
		Aws             awsdb.AwsConfig
		Grpc            GrpcConfig
		Jwt             JwtConfig
		Rest            RestConfig
		Mgo             mgo.Config
		ShutdownTimeout time.Duration `default:"10s"`
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
