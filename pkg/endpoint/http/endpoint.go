package http

import (
	"net/http"

	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/keyauth/pkg"
	"github.com/infraboard/keyauth/pkg/endpoint"
)

// CreateApplication 创建自定义角色
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	tk, err := pkg.GetTokenFromContext(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	req := endpoint.NewDefaultRegistryRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.WithToken(tk)

	err = h.endpoint.Registry(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, req)
	return
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	page := request.NewPageRequestFromHTTP(r)
	req := endpoint.NewQueryEndpointRequest(page)

	set, err := h.endpoint.QueryEndpoints(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
	return
}