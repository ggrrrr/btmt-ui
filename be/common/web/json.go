package web

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
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
		log.Log().Error(err, "decode", slog.String("body", string(b)))
		return app.BadRequestError("bad json", err)
	}
	return nil
}
