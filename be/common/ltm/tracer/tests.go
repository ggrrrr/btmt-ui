package tracer

import "context"

func ConfigureForTest() error {

	return Configure(context.Background(), "test-app", Config{
		Client: ClientCfg{
			Target: TargetCfg{
				Addr: "localhost:4317",
			},
		},
	})

}
