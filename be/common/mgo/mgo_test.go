package mgo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ggrrrr/btmt-ui/be/common/config"
)

type (
	appCfg1 struct {
		MGO Config `envPrefix:"APP1_"`
	}
	appCfg2 struct {
		MGO Config `envPrefix:"APP2_"`
	}
)

func TestCfg(t *testing.T) {
	os.Setenv("APP1_MGO_USER", "APP1_USER")
	os.Setenv("APP2_MGO_USER", "APP2_USER")

	cfg1 := appCfg1{}
	cfg2 := appCfg2{}
	cfg22 := appCfg2{}

	config.MustParse(&cfg1)
	config.MustParse(&cfg2)
	config.MustParse(&cfg22)

	assert.Equal(t, cfg1.MGO.User, "APP1_USER")
	assert.Equal(t, cfg2.MGO.User, "APP2_USER")
	assert.Equal(t, cfg22.MGO.User, "APP2_USER")

}
