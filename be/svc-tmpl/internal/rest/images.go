package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	appError "github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

/*
	curl -v -H'Authorization: mock admin' \
		-F "file=@glass-mug-variant.png" \
		http://localhost:8010/tmpl/image
*/
func (s *server) UploadImage(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.UploadImage", nil)
	defer func() {
		span.End(err)
	}()
	authInfo := roles.AuthInfoFromCtx(ctx)
	if authInfo.Subject == "" {
		web.SendError(ctx, w, appError.ErrAuthUnauthenticated)
	}

	tmpFiles, err := web.HandleFileUpload(ctx, r)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("hanlde uploads")
		web.SendError(ctx, w, err)
		return
	}
	imageForm, ok := tmpFiles["file"]
	if !ok {
		err := appError.BadRequestError("image form[file] empty", nil)
		logger.ErrorCtx(ctx, err).Send()
		web.SendError(ctx, w, err)
		return
	}
	if len(imageForm) == 0 {
		err := appError.BadRequestError("image form no uploads", nil)
		logger.ErrorCtx(ctx, err).Send()
		web.SendError(ctx, w, err)
		return
	}

	err = s.app.PutImage(ctx, imageForm[0])
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	// time.Sleep(5 * time.Second)
	web.SendPayload(ctx, w, "ok", nil)

}

/*
curl http://localhost:8010/tmpl/images
*/
func (s *server) GetImage(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.GetFile", nil)
	defer func() {
		span.End(err)
	}()

	fileId := chi.URLParam(r, "id")
	logger.DebugCtx(ctx).
		Str("path", r.URL.Path).
		Str("fileId", fileId).
		Msg("rest.GetImage")

	// fileId := chi.URLParam(r, "id")
	// fileIds := strings.Split(r.URL.Path, "/")
	// fileId := fileIds[len(fileIds)-1]
	download := r.URL.Query().Get("download")

	attch, err := s.app.GetImage(ctx, fileId, 0)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("rest.GetImage")
		web.SendError(ctx, w, err)
		return
	}

	w.Header().Add("Content-Type", attch.ContentType)

	if download != "" {
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", attch.Name))
	}

	w.WriteHeader(http.StatusOK)

	_, err = attch.WriterTo.WriteTo(w)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("unable to write file")
	}
}

/*
curl http://localhost:8010/tmpl/image/resized/
*/
func (s *server) GetResizedImage(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.GetResizedImage", nil)
	defer func() {
		span.End(err)
	}()

	download := r.URL.Query().Get("download")

	fileId := chi.URLParam(r, "id")
	logger.DebugCtx(ctx).
		Str("path", r.URL.Path).
		Str("fileId", fileId).
		Msg("rest.GetResizedImage")

	attch, err := s.app.GetResizedImage(ctx, fileId)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	w.Header().Add("Content-Type", attch.ContentType)

	if download != "" {
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", attch.Name))
	}

	w.WriteHeader(http.StatusOK)

	_, err = attch.WriterTo.WriteTo(w)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("unable to write file")
	}

}

/*
	curl -v -H'Authorization: mock admin' \
		http://localhost:8010/tmpl/images
*/
func (s *server) ListImages(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.ListImages", nil)
	defer func() {
		span.End(err)
	}()

	images, err := s.app.ListImages(ctx)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	var response = struct {
		List []ddd.ImageInfo `json:"list"`
	}{
		List: make([]ddd.ImageInfo, 0, len(images)),
	}

	response.List = append(response.List, images...)
	web.SendPayload(ctx, w, "ok", response)
}
