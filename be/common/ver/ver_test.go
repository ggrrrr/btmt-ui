package ver

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

	assert.Equal(t, "dev_build", BuildVersion(fmt.Sprintf("%s/be/common/ver/build_ver.txt", pwd)))

	assert.WithinDuration(t, time.Now().UTC(), BuildTime(""), 500*time.Microsecond)
	assert.WithinDuration(t, time.Now().UTC(), BuildTime("asd"), 500*time.Microsecond)

	assert.WithinDuration(t, time.Now().UTC(), BuildTime(fmt.Sprintf("%s/be/common/ver/build_ts_nok.txt", pwd)), 500*time.Microsecond)

	ts1, err := time.Parse(time.RFC3339Nano, "2024-11-14T08:29:13Z")
	require.NoError(t, err)

	assert.WithinDuration(t, ts1, BuildTime(fmt.Sprintf("%s/be/common/ver/build_ts.txt", pwd)), 500*time.Microsecond)
}
