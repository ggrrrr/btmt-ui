package roles

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mileusna/useragent"

	"github.com/ggrrrr/btmt-ui/be/common/app"
)

const (
	CookieName   string = "authorization"
	CookieSchema string = "Bearer"

	HttpForwardedFor  string = "X-Forwarded-For"
	HttpForwardedHost string = "X-Forwarded-Host"

	HttpAuthorization string = "Authorization"
	HttpUserAgent     string = "User-Agent"
)

func FromHttpRequest(md http.Header, cookies []*http.Cookie, r *http.Request) app.RequestIn {
	out := app.RequestIn{
		FullMethod: r.RequestURI,
	}

	auth, ok := authorizationFromCokies(cookies)
	if !ok {
		auth, _ = authorizationFromHeaders(md)
	}

	out.AuthData = auth
	out.Device = extractHttpDevice(md)
	remoteAddr := strings.LastIndex(r.RemoteAddr, ":")

	out.Device.RemoteAddr = r.RemoteAddr[0:remoteAddr]

	ua := useragent.Parse(out.Device.DeviceInfo)

	out.Device.DeviceInfo = fmt.Sprintf("%s/%s/%s", ua.OS, ua.OSVersionNoFull(), ua.Name)

	return out
}

func authorizationFromCokies(cookies []*http.Cookie) (app.AuthData, bool) {
	for _, v := range cookies {
		if v.Name == CookieName {
			return app.AuthData{
				AuthScheme: CookieSchema,
				AuthToken:  v.Value,
			}, true
		}
	}
	return app.AuthData{}, false
}

func authorizationFromHeaders(md http.Header) (app.AuthData, bool) {
	out := app.AuthData{}
	if len(md[HttpAuthorization]) == 1 {
		// return
		out = app.AuthDataFromValue(md[HttpAuthorization][0])
		if !out.IsZero() {
			return out, true
		}
		// gwAuthorization := strings.Split(md[HttpAuthorization][0], " ")
		// if len(gwAuthorization) == 2 {
		// 	out.AuthScheme = gwAuthorization[0]
		// 	out.AuthToken = gwAuthorization[1]
		// 	return out, true
		// }
	}
	return out, false
}

func extractHttpDevice(md http.Header) app.Device {
	out := app.Device{}
	if len(md[HttpUserAgent]) > 0 {
		out.DeviceInfo = strings.Join(md[HttpUserAgent], ",")
	}
	if len(md[HttpForwardedFor]) > 0 {
		out.RemoteAddr = strings.Join(md[HttpForwardedFor], ",")
	}
	return out
}
