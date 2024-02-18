package roles

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	GrpcAuthorization string = "grpcgateway-authorization"
	GrpcUserAgent     string = "grpcgateway-user-agent"
)

func FromGrpcMetadata(md metadata.MD, fullMethod string) UserRequest {
	out := UserRequest{
		FullMethod: fullMethod,
	}
	out.Authorization = extractGrpcAuthorization(md)
	out.Device = extractGrpcDevice(md)
	return out
}

func lower(from string) string {
	return strings.ToLower(from)
}

func extractGrpcAuthorization(md metadata.MD) Authorization {
	out := Authorization{}
	// We check first for grpc specific header 'authorization'
	if len(md[lower(HttpAuthorization)]) == 1 {
		gwAuthorization := strings.Split(md[lower(HttpAuthorization)][0], " ")
		if len(gwAuthorization) == 2 {
			out.AuthScheme = gwAuthorization[0]
			out.AuthCredentials = AuthCredentials(gwAuthorization[1])
			return out
		}
	}
	// We check first for http specific forwarded header (part of grpc-gateway)
	if len(md[GrpcAuthorization]) == 1 {
		gwAuthorization := strings.Split(md[GrpcAuthorization][0], " ")
		if len(gwAuthorization) == 2 {
			out.AuthScheme = gwAuthorization[0]
			out.AuthCredentials = AuthCredentials(gwAuthorization[1])
			return out
		}
	}
	return out
}

func extractGrpcDevice(md metadata.MD) Device {
	out := Device{}
	if len(md[lower(HttpUserAgent)]) > 0 {
		out.DeviceInfo = strings.Join(md[lower(HttpUserAgent)], ",")
	}
	// if we have grpc gateway set header this
	if len(md[GrpcUserAgent]) > 0 {
		out.DeviceInfo = strings.Join(md[GrpcUserAgent], ",")
	}
	if len(md[lower(HttpForwardedFor)]) > 0 {
		out.RemoteAddr = strings.Join(md[lower(HttpForwardedFor)], ",")
	}
	return out
}
