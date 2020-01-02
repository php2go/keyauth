package http

import (
	"net/http"

	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/keyauth/pkg/application"
)

func (h *handler) QueryUserApplication(w http.ResponseWriter, r *http.Request) {
	page := request.LoadPagginFromReq(r)
	req := application.NewQueryApplicationRequest(page)

	apps, total, err := h.service.QueryApplication(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	data := response.PageData{
		PageRequest: page,
		TotalCount:  uint(total),
		List:        apps,
	}
	response.Success(w, data)
	return
}

// CreateApplication 创建主账号
func (h *handler) CreateUserApplication(w http.ResponseWriter, r *http.Request) {
	req := application.NewCreateApplicatonRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	d, err := h.service.CreateUserApplication("xxx", req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, d)
	return
}

func (h *handler) GetApplication(w http.ResponseWriter, r *http.Request) {
	rctx := context.GetContext(r)

	req := application.NewDescriptApplicationRequest()
	req.ID = rctx.PS.ByName("id")
	d, err := h.service.DescriptionApplication(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, d)
	return
}

// DestroyPrimaryAccount 注销账号
func (h *handler) DestroyApplication(w http.ResponseWriter, r *http.Request) {
	rctx := context.GetContext(r)

	if err := h.service.DeleteApplication(rctx.PS.ByName("id")); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, "delete ok")
	return
}
