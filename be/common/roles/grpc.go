package roles

import (
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/ggrrrr/btmt-ui/be/common/app"
)

const (
	GrpcAuthorization string = "grpcgateway-authorization"
	GrpcUserAgent     string = "grpcgateway-user-agent"
)

func FromGrpcMetadata(md metadata.MD, fullMethod string) app.RequestIn {
	out := app.RequestIn{
		FullMethod: fullMethod,
	}
	out.AuthData = extractGrpcAuthorization(md)
	out.Device = extractGrpcDevice(md)
	return out
}

func extractGrpcAuthorization(md metadata.MD) app.AuthData {
	out := app.AuthData{}
	// We check first for http specific header 'authorization'
	if len(md[strings.ToLower(HttpAuthorization)]) == 1 {
		gwAuthorization := strings.Split(md[strings.ToLower(HttpAuthorization)][0], " ")
		if len(gwAuthorization) == 2 {
			out.AuthScheme = gwAuthorization[0]
			out.AuthToken = gwAuthorization[1]
			return out
		}
	}
	// We check first for http specific forwarded header (part of grpc-gateway)
	if len(md[GrpcAuthorization]) == 1 {
		gwAuthorization := strings.Split(md[GrpcAuthorization][0], " ")
		if len(gwAuthorization) == 2 {
			out.AuthScheme = gwAuthorization[0]
			out.AuthToken = gwAuthorization[1]
			return out
		}
	}
	return out
}

func extractGrpcDevice(md metadata.MD) app.Device {
	out := app.Device{}
	if len(md[strings.ToLower(HttpUserAgent)]) > 0 {
		out.DeviceInfo = strings.Join(md[strings.ToLower(HttpUserAgent)], ",")
	}
	// if we have grpc gateway set header this
	if len(md[GrpcUserAgent]) > 0 {
		out.DeviceInfo = strings.Join(md[GrpcUserAgent], ",")
	}
	if len(md[strings.ToLower(HttpForwardedFor)]) > 0 {
		out.RemoteAddr = strings.Join(md[strings.ToLower(HttpForwardedFor)], ",")
	}
	return out
}
