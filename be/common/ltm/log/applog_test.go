package log

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

type someData1 struct{}

var _ (td.TraceDataExtractor) = (*someData1)(nil)

func (someData1) Extract() *td.TraceData {
	return &td.TraceData{
		KV: map[string]slog.Value{
			"bol": slog.BoolValue(true),
			"group1": slog.GroupValue(
				slog.String("g.str1", "val1"),
				slog.String("g.str2", "val2"),
			),
		},
	}
}

type someDataNil struct{}

var _ (td.TraceDataExtractor) = someDataNil{}

func (someDataNil) Extract() *td.TraceData {
	return nil
}

func TestLog(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	cfg := Config{
		Level:  "info",
		Format: "json",
	}
	sd1 := someData1{}
	ctx = td.Inject(ctx, sd1)
	ctx = td.InjectGroup(ctx, "sd2", slog.Bool("b2", true))

	testLogger := configureWithWriter(cfg, &buf)

	testLogger.logCtx(ctx, slog.LevelInfo, fmt.Errorf("err1"), "asd", slog.String("add", "val1"))

	fmt.Printf("%v \n", buf.String())
}
