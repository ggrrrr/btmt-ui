package people

import (
	"fmt"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/config"
)

func Test(t *testing.T) {

	config.DumpEnv()

	cfg := moduleCfg{}
	config.MustParse(&cfg)
	fmt.Printf("%+v \n", cfg)

}
