package pin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
)

func TestEgn(t *testing.T) {
	egn2 := os.Getenv("PIN1")
	if egn2 == "" {
		t.Skip()
	}

	egn1 := "o645196854"
	_, err := Parse(egn1)
	assert.Error(t, err)

	res, err := Parse(egn2)
	require.NoError(t, err)
	assert.Equal(t, res, ddd.PinValidation{
		DOB:    ddd.DOB{Year: 1942, Month: 2, Day: 13},
		Gender: "male",
	})
	t.Logf("%+v", res)

	res, err = Parse(os.Getenv("PIN2"))
	require.NoError(t, err)
	assert.Equal(t, res, ddd.PinValidation{
		DOB:    ddd.DOB{Year: 1978, Month: 2, Day: 13},
		Gender: "male",
	})

}
