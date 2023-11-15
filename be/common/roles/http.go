package roles

import (
	"net/http"
	"strings"
)

const (
	HttpForwardedFor  string = "X-Forwarded-For"
	HttpForwardedHost string = "X-Forwarded-Host"

	HttpAuthorization string = "Authorization"
	HttpUserAgent     string = "User-Agent"
)

func FromHttpMetadata(md http.Header, fullMethod string) UserRequest {
	out := UserRequest{
		FullMethod: fullMethod,
	}
	out.Authorization = extractHttpAuthorization(md)
	out.Device = extractHttpDevice(md)
	return out
}

func extractHttpAuthorization(md http.Header) Authorization {
	out := Authorization{}
	// We check first for grpc specific header 'authorization'
	if len(md[HttpAuthorization]) == 1 {
		gwAuthorization := strings.Split(md[HttpAuthorization][0], " ")
		if len(gwAuthorization) == 2 {
			out.AuthScheme = gwAuthorization[0]
			out.AuthCredentials = AuthCredentials(gwAuthorization[1])
			return out
		}
	}
	return out
}

func extractHttpDevice(md http.Header) Device {
	out := Device{}
	if len(md[lower(HttpUserAgent)]) > 0 {
		out.DeviceInfo = strings.Join(md[HttpUserAgent], ",")
	}
	if len(md[lower(HttpForwardedFor)]) > 0 {
		out.RemoteAddr = strings.Join(md[HttpForwardedFor], ",")
	}
	return out
}
