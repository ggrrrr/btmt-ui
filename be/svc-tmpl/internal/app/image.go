package app

import (
	"context"
	"fmt"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func headImage(name string) (blob.ImageInfo, error) {
	reader, err := os.Open(name)
	if err != nil {
		return blob.ImageInfo{}, fmt.Errorf("headImage[%s] %w", name, err)
	}
	m, _, err := image.Decode(reader)
	if err != nil {
		return blob.ImageInfo{}, fmt.Errorf("headImage[%s] %w", name, err)
	}
	defer reader.Close()

	bounds := m.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	return blob.ImageInfo{
		Width:  int64(w),
		Height: int64(h),
	}, nil
}

func (a *App) PutImage(ctx context.Context, tempFile blob.TempFile) error {
	defer func() {
		tempFile.Delete(ctx)
	}()

	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "PutImage", nil,
		logger.AttributeString("FileName", tempFile.FileName),
		logger.AttributeString("ContentType", tempFile.ContentType),
		logger.AttributeString("TempFileName", tempFile.TempFileName),
	)
	defer func() {
		span.End(err)
	}()

	logger.InfoCtx(ctx).
		Any("file", tempFile).
		Msg("PutImage")

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Tenant, "some", authInfo)
	if err != nil {
		return err
	}

	blobInfo := blob.BlobInfo{
		Type:        tempFile.ContentType,
		ContentType: tempFile.ContentType,
		Name:        tempFile.FileName,
	}

	imageInfo, err1 := headImage(tempFile.TempFileName)
	if err1 == nil {
		blobInfo.ImageInfo = imageInfo
	}

	reader, err := os.Open(tempFile.TempFileName)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("tempFile", tempFile).Msg("PutFile")
		return err
	}
	defer reader.Close()

	id, err := a.blobFetcher.Push(ctx, authInfo.Tenant, fmt.Sprintf("%s/%s", a.imagesFolder, tempFile.FileName), blobInfo, reader)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("tempFile", tempFile).Msg("PutFile")
		return err
	}
	logger.InfoCtx(ctx).Str("id", id.String()).Msg("PutFile")

	return nil
}

func (a *App) GetImage(ctx context.Context, fileId string) (*ddd.FileWriterTo, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "GetImage", nil, logger.AttributeString("fileId", fileId))
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Tenant, "some", authInfo)
	if err != nil {
		return nil, err
	}

	logger.InfoCtx(ctx).
		Any("info", fileId).
		Msg("GetFile")

	res, err := a.blobFetcher.Fetch(ctx, "localhost", fmt.Sprintf("%s/%s", a.imagesFolder, fileId))
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Fetch")
		return nil, err
	}

	logger.InfoCtx(ctx).
		Any("info", res).
		Any("id", fileId).
		Msg("Fetch")

	return &ddd.FileWriterTo{
			ContentType: res.Info.ContentType,
			Version:     res.Id.Version(),
			Name:        res.Info.Name,
			WriterTo: &filePipe{
				reader: res.ReadCloser,
			},
		},
		nil
}

func (a *App) ListImages(ctx context.Context) ([]string, error) {
	var err error
	ctx, span := logger.Span(ctx, "ListImages", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Tenant, "some", authInfo)
	if err != nil {
		return nil, err
	}

	logger.InfoCtx(ctx).
		Msg("ListImages")

	res, err := a.blobFetcher.List(ctx, "localhost", fmt.Sprintf("%s", a.imagesFolder))
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ListImages.List")
		return nil, err
	}
	out := []string{}
	for _, v := range res {
		out = append(out, v.Id.String())
	}
	logger.InfoCtx(ctx).
		Any("info", res).
		Msg("Fetch")

	return out, nil
}
