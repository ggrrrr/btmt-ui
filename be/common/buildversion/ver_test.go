package buildversion

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/help"
)

func TestVer(t *testing.T) {
	pwd := help.RepoDir()

	assert.Equal(t, "dev", BuildVersion(""))
	assert.Equal(t, "dev", BuildVersion("asdasd"))
	assert.Equal(t, "dev_build", BuildVersion(fmt.Sprintf("%s/be/common/buildversion/build_ver.txt", pwd)))

	assert.WithinDuration(t, time.Now().UTC(), BuildTime(""), 500*time.Microsecond)
	testBuildTs(t, time.Now().UTC(), BuildTime(""))

	testBuildTs(t, time.Now().UTC(), BuildTime(fmt.Sprintf("%s/be/common/buildversion/notfound", pwd)))

	testBuildTs(t, time.Now().UTC(), BuildTime(fmt.Sprintf("%s/be/common/buildversion/build_ts_nok.txt", pwd)))

	ts1, err := time.Parse(time.RFC3339Nano, "2024-11-14T08:29:13Z")
	require.NoError(t, err)
	testBuildTs(t, ts1, BuildTime(fmt.Sprintf("%s/be/common/buildversion/build_ts.txt", pwd)))
}

func testBuildTs(t *testing.T, expected time.Time, actual time.Time) {
	assert.WithinDuration(t, expected, actual, 100*time.Millisecond)
}
