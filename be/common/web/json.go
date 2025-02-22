package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

// return system error on ReadAll
// return BadRequestError on Decode
func DecodeJsonRequest(r *http.Request, payload any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return app.SystemError("unable to load http body", err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&payload)
	if err != nil {
		logger.Error(err).Str("body", string(b)).Send()
		return app.BadRequestError("bad json", err)
	}
	return nil
}
