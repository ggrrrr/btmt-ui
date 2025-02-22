package config

import (
	"fmt"
	"os"
	"sort"

	"github.com/caarlos0/env/v11"
)

func Parse(cfg any) error {
	opts := env.Options{
		// DefaultValueTagName: "default",
		// PrefixTagName: "envPrefix",
	}
	// DefaultValueTagName
	// if prefix != "" {
	// 	opts.Prefix = fmt.Sprintf("%s_", prefix)
	// }

	return env.ParseWithOptions(cfg, opts)
}

func MustParse(cfg any) {
	if err := Parse(cfg); err != nil {
		panic(err)
	}
}

func DumpEnv() {
	l := os.Environ()

	sort.Slice(l, func(i, j int) bool {
		return l[i] < l[j]
	})

	for v := range l {
		fmt.Printf("\t %v \n", l[v])
	}
}
