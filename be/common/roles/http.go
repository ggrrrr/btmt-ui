package roles

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mileusna/useragent"
)

const (
	CookieName   string = "authorization"
	CookieSchema string = "Bearer"

	HttpForwardedFor  string = "X-Forwarded-For"
	HttpForwardedHost string = "X-Forwarded-Host"

	HttpAuthorization string = "Authorization"
	HttpUserAgent     string = "User-Agent"
)

func FromHttp(md http.Header, cookies []*http.Cookie, fullMethod string) UserRequest {
	out := UserRequest{
		FullMethod: fullMethod,
	}

	auth, ok := authorizationFromCokies(cookies)
	if !ok {
		auth, _ = authorizationFromHeaders(md)
	}

	out.Authorization = auth
	out.Device = extractHttpDevice(md)

	ua := useragent.Parse(out.Device.DeviceInfo)

	out.Device.DeviceInfo = fmt.Sprintf("%s/%s/%s", ua.OS, ua.OSVersionNoFull(), ua.Name)

	return out
}

func authorizationFromCokies(cookies []*http.Cookie) (Authorization, bool) {
	for _, v := range cookies {
		if v.Name == CookieName {
			return Authorization{
				AuthScheme:      CookieSchema,
				AuthCredentials: AuthCredentials(v.Value),
			}, true
		}
	}
	return Authorization{}, false
}

func authorizationFromHeaders(md http.Header) (Authorization, bool) {
	out := Authorization{}
	if len(md[HttpAuthorization]) == 1 {
		gwAuthorization := strings.Split(md[HttpAuthorization][0], " ")
		if len(gwAuthorization) == 2 {
			out.AuthScheme = gwAuthorization[0]
			out.AuthCredentials = AuthCredentials(gwAuthorization[1])
			return out, true
		}
	}
	return out, false
}

func extractHttpDevice(md http.Header) Device {
	out := Device{}
	if len(md[HttpUserAgent]) > 0 {
		out.DeviceInfo = strings.Join(md[HttpUserAgent], ",")
	}
	if len(md[HttpForwardedFor]) > 0 {
		out.RemoteAddr = strings.Join(md[HttpForwardedFor], ",")
	}
	return out
}
