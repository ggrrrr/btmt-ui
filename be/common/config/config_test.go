package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/config"
)

type (
	SubVar struct {
		Sub string `env:"SUB"`
	}

	TestVar struct {
		Val1   string        `env:"VAL"`
		Val2   string        `env:"VAL2" envDefault:"val2"`
		Since  time.Duration `env:"SINCE"`
		SubVar SubVar        `envPrefix:"SUB_"`
	}
)

func TestP(t *testing.T) {
	os.Setenv("VAL", "1")
	os.Setenv("SINCE", "2s")

	v1 := TestVar{}
	v2 := TestVar{}
	config.MustParse(&v1)
	config.MustParse(&v2)

	fmt.Printf("%+v \n", v1)
	fmt.Printf("%+v \n", v2)
}

func Test_Parse(t *testing.T) {
	tests := []struct {
		name      string
		prepFunc  func()
		val       TestVar
		willPanic bool
	}{
		{
			name:      "ok",
			willPanic: false,
			prepFunc: func() {
				os.Setenv("VAL", "1")
				os.Setenv("SINCE", "2s")
				os.Setenv("SUB_SUB", "sub data")
			},
			val: TestVar{Val1: "1", Val2: "val2", Since: time.Second * 2, SubVar: SubVar{Sub: "sub data"}},
		},
		{
			name:      "panic",
			willPanic: true,
			prepFunc: func() {
				os.Setenv("VAL", "1")
				os.Setenv("SINCE", "asdasd")
				os.Setenv("SUB_SUB", "sub data")
			},
			val: TestVar{Val1: "1", Val2: "val2", Since: time.Second * 2, SubVar: SubVar{Sub: "sub data"}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFunc()
			v := TestVar{}
			func() {
				defer func() {
					rvr := recover()
					if tc.willPanic {
						fmt.Printf("recover %+v \n ", rvr)
						require.NotNil(t, rvr)
					} else {
						require.Equal(t, tc.val, v)
						require.Nil(t, rvr)
						config.DumpEnv()
					}
				}()
				config.MustParse(&v)
			}()
		})
	}
}
