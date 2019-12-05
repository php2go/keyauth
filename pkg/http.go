package pkg

import (
	"strings"

	"github.com/infraboard/mcube/http/router"
)

var (
	v1httpAPIs = make(map[string]HTTPAPI)
)

// HTTPAPI restful 服务
type HTTPAPI interface {
	Registry(router router.SubRouter)
	Config() error
}

// RegistryHTTPV1 注册HTTP服务
func RegistryHTTPV1(name string, api HTTPAPI) {
	if _, ok := v1httpAPIs[name]; ok {
		panic("http api " + name + " has registry")
	}
	v1httpAPIs[name] = api
}

// InitV1HTTPAPI 初始化API服务
func InitV1HTTPAPI(pathPrefix string, root router.Router) error {
	for name, api := range v1httpAPIs {
		if err := api.Config(); err != nil {
			return err
		}

		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.SubRouter(pathPrefix + "/v1/" + name))
	}

	return nil
}
