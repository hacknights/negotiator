package negotiator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Factory func(w http.ResponseWriter, r *http.Request) NegotiatorFunc

type NegotiatorFunc func(interface{}, int, error)

func NewNegotiator(w http.ResponseWriter, r *http.Request) NegotiatorFunc {
	type response struct {
		Ok      bool        `json:"ok"`
		Errors  []string    `json:"errors,omitempty"`
		Content interface{} `json:"content,omitempty"`
	}
	newResponse := func(value interface{}, status int, err error) response {
		res := response{
			Ok:      err == nil,
			Content: value,
		}
		if err != nil {
			res.Errors = append(res.Errors, err.Error())
		}
		return res
	}

	return func(value interface{}, status int, err error) {

		res := newResponse(value, status, err)

		json, e := json.Marshal(res)
		if e != nil {
			http.Error(w, "unable to negotiate response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(json)
	}
}

func (n NegotiatorFunc) OK(value interface{}) {
	n(value, http.StatusOK, nil)
}

func (n NegotiatorFunc) NotFound() {
	n(nil, http.StatusNotFound, fmt.Errorf("not found"))
}

func (n NegotiatorFunc) InternalServer(err string) {
	n.InternalServerError(fmt.Errorf(err))
}

func (n NegotiatorFunc) InternalServerError(err error) {
	n(nil, http.StatusInternalServerError, err)
}

func (n NegotiatorFunc) Unauthorized(err string) {
	n.UnauthorizedError(fmt.Errorf(err))
}

func (n NegotiatorFunc) UnauthorizedError(err error) {
	n(nil, http.StatusUnauthorized, err)
}

func (n NegotiatorFunc) BadRequest(err string) {
	n.BadRequestError(fmt.Errorf(err))
}

func (n NegotiatorFunc) BadRequestError(err error) {
	n(nil, http.StatusBadRequest, err)
}
