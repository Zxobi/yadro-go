package protocol

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

func ResponseError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_ = ResponseJson(w, errResp{Error: msg})
}

func ResponseJson(w http.ResponseWriter, v any) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(res)
	return err
}
