package endpoint

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/types/ftime"

	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/keyauth/pkg/user/types"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

// Service token管理服务
type Service interface {
	QueryEndpoints(req *QueryEndpointRequest) (*Set, error)
	Registry(req *RegistryRequest) error
}

// NewRegistryRequest 注册请求
func NewRegistryRequest(name, version string, entries []*router.Entry) *RegistryRequest {
	return &RegistryRequest{
		Service: name,
		Version: version,
		Entries: entries,
	}
}

// NewDefaultRegistryRequest todo
func NewDefaultRegistryRequest() *RegistryRequest {
	return &RegistryRequest{
		Session: token.NewSession(),
		Entries: []*router.Entry{},
	}
}

// RegistryRequest 服务注册请求
type RegistryRequest struct {
	*token.Session
	Service string          `json:"service" validate:"required,lte=64"`
	Version string          `json:"version" validate:"lte=32"`
	Entries []*router.Entry `json:"entries"`
}

// Validate 校验注册请求合法性
func (req *RegistryRequest) Validate() error {
	tk := req.GetToken()
	if !tk.UserType.Is(types.ServiceAccount) {
		return fmt.Errorf("only service account can registry endpoints")
	}

	return validate.Struct(req)
}

// Endpoints 功能列表
func (req *RegistryRequest) Endpoints() []*Endpoint {
	eps := make([]*Endpoint, 0, len(req.Entries))
	for i := range req.Entries {
		et := req.Entries[i]
		eps = append(eps, &Endpoint{
			CreateAt: ftime.Now(),
			UpdateAt: ftime.Now(),
			Service:  req.Service,
			Version:  req.Version,
			Entry:    *et,
		})
	}
	return eps
}

// NewQueryEndpointRequest 列表查询请求
func NewQueryEndpointRequest(pageReq *request.PageRequest) *QueryEndpointRequest {
	return &QueryEndpointRequest{
		PageRequest: pageReq,
	}
}

// QueryEndpointRequest 查询应用列表
type QueryEndpointRequest struct {
	*request.PageRequest
}

// DescribeEndpointRequest todo
type DescribeEndpointRequest struct {
	Name string
}