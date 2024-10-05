package rest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	appError "github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/app"
	"github.com/go-chi/chi/v5"
)

type server struct {
	app *app.App
}

func New(a *app.App) *server {
	return &server{
		app: a,
	}
}

func (s *server) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/render", s.Render)
	router.Get("/image/*", s.GetImage)
	router.Post("/image", s.UploadImage)
	router.Get("/images", s.ListImages)
	router.Get("/get", s.GetTmpl)

	return router
}

func (s *server) GetTmpl(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context()).Msg("rest.GetTmpl")

	web.SendError(r.Context(), w, fmt.Errorf("asdasd"))
}

func (s *server) Render(w http.ResponseWriter, r *http.Request) {
}

/*
curl -v -F "image=@glass-mug-variant.png"  http://localhost:8010/tmpl/file
*/
func (s *server) UploadImage(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.UploadImage", nil)
	defer func() {
		span.End(err)
	}()
	authInfo := roles.AuthInfoFromCtx(ctx)
	if authInfo.User == "" {
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
		err := appError.BadRequestError("image form empty", nil)
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
	time.Sleep(5 * time.Second)
	web.SendPayload(ctx, w, "ok", nil)

}

func (s *server) GetImage(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.GetFile", nil)
	defer func() {
		span.End(err)
	}()

	// fileId := chi.URLParam(r, "id")
	fileIds := strings.Split(r.URL.Path, "tmpl/image/")
	fileId := fileIds[len(fileIds)-1]
	download := r.URL.Query().Get("download")
	logger.DebugCtx(ctx).
		Str("id", fileId).
		Msg("rest.GetFile")

	attch, err := s.app.GetImage(ctx, fileId)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	w.Header().Add("Content-Type", attch.ContentType)

	if download != "" {
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", attch.Name))
	}

	w.WriteHeader(http.StatusOK)

	attch.WriterTo.WriteTo(w)

}

func (s *server) ListImages(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.ListImages", nil)
	defer func() {
		span.End(err)
	}()

	attch, err := s.app.ListImages(ctx)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	var response = struct {
		List []string
	}{
		List: attch,
	}
	web.SendPayload(ctx, w, "ok", response)
}
