package pin

import (
	"os"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEgn(t *testing.T) {
	egn1 := "o645196854"
	_, err := Parse(egn1)
	assert.Error(t, err)

	egn2 := os.Getenv("PIN1")
	res, err := Parse(egn2)
	require.NoError(t, err)
	assert.Equal(t, res, ddd.PinValidation{
		Dob:    ddd.Dob{Year: 1942, Month: 2, Day: 13},
		Gender: "male",
	})
	t.Logf("%+v", res)

	res, err = Parse(os.Getenv("PIN2"))
	require.NoError(t, err)
	assert.Equal(t, res, ddd.PinValidation{
		Dob:    ddd.Dob{Year: 1978, Month: 2, Day: 13},
		Gender: "male",
	})

}
