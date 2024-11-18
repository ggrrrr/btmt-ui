package app

import (
	"strings"
)

type (
	Device struct {
		RemoteAddr string
		DeviceInfo string
	}

	AuthData struct {
		// Basic, Bearer, ...,
		AuthScheme string
		// JWT TOKEN OR OTHER secret data
		AuthToken string
	}

	// This is request commint in this service
	RequestIn struct {
		FullMethod string
		// Requester device/client/useragent  info/
		Device Device
		// From header/metadata
		AuthData AuthData
	}
)

func (a AuthData) IsZero() bool {
	if a.AuthScheme == "" && a.AuthToken == "" {
		return true
	}
	return false
}

func AuthDataFromValue(val string) AuthData {
	valSlice := strings.Split(val, " ")
	if len(valSlice) != 2 {
		return AuthData{}
	}
	return AuthData{
		AuthScheme: valSlice[0],
		AuthToken:  valSlice[1],
	}
}
