package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Validator interface {
	Validate() error
}

type Wrapper[Req Validator, Res any] struct {
	fn func(context.Context, Req) (Res, error)
}

func New[Req Validator, Res any](fn func(ctx context.Context, req Req) (Res, error)) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res]{fn: fn}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (wrapper *Wrapper[Req, Res]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req Req
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, "decoding JSON", err, http.StatusBadRequest)
		return
	}

	err = req.Validate()
	if err != nil {
		writeError(w, "validating request", err, http.StatusBadRequest)
		return
	}

	res, err := wrapper.fn(ctx, req)
	if err != nil {
		writeError(w, "running handler", err, http.StatusInternalServerError)
		return
	}

	err = writeJSON(w, res, http.StatusOK)
	if err != nil {
		writeError(w, "encoding JSON", err, http.StatusInternalServerError)
		return
	}
}

func writeJSON[Res any](w http.ResponseWriter, res Res, statusCode int) error {
	rawJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "applicaton/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(rawJSON)
	return err
}

func writeError(w http.ResponseWriter, text string, err error, statusCode int) error {
	buf := bytes.NewBufferString(text)
	buf.WriteString(": ")
	buf.WriteString(err.Error())

	res := ErrorResponse{Error: buf.String()}
	return writeJSON(w, res, statusCode)
}
