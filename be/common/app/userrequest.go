package app

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
