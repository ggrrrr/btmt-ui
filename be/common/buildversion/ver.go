/*

This folder is to enable setting app version during build
Check Dockerfile and be/common/system/system.go

**/

package buildversion

import (
	"strings"
)

/*

go run -ldflags "-X github.com/ggrrrr/btmt-ui/be/common/buildversion.version=SOME_VER" be/monolith/main.go

*/

var version string = "dev"

func BuildVersion() string {
	return strings.TrimSpace(version)
}
