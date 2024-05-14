package http

import (
	"encoding/json"
	"net/http"
)

type UpdateResponse struct {
	Total int `json:"total"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type errResp struct {
	Error string `json:"error"`
}

func responseError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_ = responseJson(w, errResp{Error: msg})
}

func responseJson(w http.ResponseWriter, v any) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(res)
	return err
}
