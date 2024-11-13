package ver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVer(t *testing.T) {
	assert.Equal(t, "dev", BuildVersion())
	assert.WithinDuration(t, time.Now().UTC(), BuildTime(), 100*time.Microsecond)
}
