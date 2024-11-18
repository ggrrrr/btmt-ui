package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthData(t *testing.T) {

	assert.True(t, AuthData{AuthScheme: "", AuthToken: ""}.IsZero() == true)
	assert.True(t, AuthData{AuthScheme: "asd", AuthToken: ""}.IsZero() == false)
	assert.True(t, AuthData{AuthScheme: "", AuthToken: "asd"}.IsZero() == false)
	assert.True(t, AuthData{AuthScheme: "asd", AuthToken: "asd"}.IsZero() == false)

	assert.Equal(t, AuthData{AuthScheme: "s", AuthToken: "t"}, AuthDataFromValue("s t"))

	assert.Equal(t, AuthData{AuthScheme: "", AuthToken: ""}, AuthDataFromValue("t"))

}
