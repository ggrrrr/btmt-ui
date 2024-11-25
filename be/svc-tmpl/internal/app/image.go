package app

import (
	"context"
	"fmt"
	"os"

	appError "github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	appImage "github.com/ggrrrr/btmt-ui/be/common/image"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func (a *App) GetResizedImage(ctx context.Context, fileId string) (*ddd.FileWriterTo, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "GetResizedImage", nil, logger.TraceKVString("fileId", fileId))
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return nil, err
	}

	imageBlobId, err := a.imagesFolder.SetIdVersionFromString(fileId)
	if err != nil {
		return nil, appError.BadRequestError("file name not supported", err)
	}

	logger.InfoCtx(ctx).
		Any("info", fileId).
		Msg("GetResizedImage")

	reader, err := a.blobStore.Fetch(ctx, "localhost", imageBlobId)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Fetch")
		return nil, err
	}

	logger.InfoCtx(ctx).
		Any("info", reader).
		Any("id", fileId).
		Msg("Fetch")

	res, err := appImage.ResizeImage(ctx, 0, reader)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Fetch")
		return nil, err
	}

	// res.BlobMd.ContentLength
	return &ddd.FileWriterTo{
			ContentType: reader.Blob.MD.ContentType,
			Version:     reader.Blob.Id.Version(),
			Name:        reader.Blob.MD.Name,
			// : res.BlobMd.ContentLength,
			WriterTo: &filePipe{
				reader: &res.ReadCloser,
			},
		},
		nil

}

func (a *App) GetImage(ctx context.Context, fileId string, maxWight int) (*ddd.FileWriterTo, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "GetImage", nil, logger.TraceKVString("fileId", fileId))
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return nil, err
	}

	imageBlobId, err := a.imagesFolder.SetIdVersionFromString(fileId)
	if err != nil {
		return nil, appError.BadRequestError("file name not supported", err)
	}

	logger.InfoCtx(ctx).
		Any("info", fileId).
		Msg("GetFile")

	res, err := a.blobStore.Fetch(ctx, "localhost", imageBlobId)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Fetch")
		return nil, err
	}

	logger.InfoCtx(ctx).
		Any("info", res).
		Any("id", fileId).
		Msg("Fetch")

	return &ddd.FileWriterTo{
			ContentType: res.Blob.MD.ContentType,
			Version:     res.Blob.Id.Version(),
			Name:        res.Blob.MD.Name,
			WriterTo: &filePipe{
				reader: res.ReadCloser,
			},
		},
		nil
}

func (a *App) ListImages(ctx context.Context) ([]ddd.ImageInfo, error) {
	var err error
	ctx, span := logger.Span(ctx, "ListImages", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return nil, err
	}

	logger.InfoCtx(ctx).
		Msg("ListImages")

	res, err := a.blobStore.List(ctx, authInfo.Realm, a.imagesFolder)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ListImages.List")
		return nil, err
	}

	out := []ddd.ImageInfo{}
	for _, v := range res {
		fmt.Printf("\t\t %v %v \n", v.Id.String(), v.Id.IdVersion())
		i := ddd.ImageInfo{
			Id:          v.Id.IdVersion(),
			Version:     v.Id.Version(),
			FileName:    v.MD.Name,
			ContentType: v.MD.ContentType,
			Width:       v.MD.ImageInfo.Width,
			Height:      v.MD.ImageInfo.Height,
			Size:        v.MD.ContentLength,
			CreatedAt:   v.MD.CreatedAt,
			Versions:    []string{},
		}
		for _, version := range v.Versions {
			i.Versions = append(i.Versions, fmt.Sprintf("%s:%s", version.Id.Id(), version.Id.Version()))
		}

		out = append(out, i)
	}

	logger.InfoCtx(ctx).
		Any("info", res).
		Msg("Fetch")

	return out, nil
}

func (a *App) PutImage(ctx context.Context, tempFile blob.TempFile) error {
	defer func() {
		tempFile.Delete(ctx)
	}()

	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "PutImage", nil,
		logger.TraceKVString("FileName", tempFile.FileName),
		logger.TraceKVString("ContentType", tempFile.ContentType),
		logger.TraceKVString("TempFileName", tempFile.TempFileName),
	)
	defer func() {
		span.End(err)
	}()

	imageBlobId, err := a.imagesFolder.SetIdVersionFromString(tempFile.FileName)
	if err != nil {
		return appError.BadRequestError("file name not supported", err)
	}

	logger.InfoCtx(ctx).
		Any("file", tempFile).
		Msg("PutImage")

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return err
	}

	blobInfo := blob.BlobMD{
		Type:        blob.BlobTypeImage,
		ContentType: tempFile.ContentType,
		Name:        tempFile.FileName,
	}

	imageInfo, err1 := appImage.HeadImage(tempFile.TempFileName)
	if err1 == nil {
		blobInfo.ImageInfo = imageInfo
	}

	reader, err := os.Open(tempFile.TempFileName)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("tempFile", tempFile).Msg("PutFile")
		return err
	}
	defer reader.Close()

	id, err := a.blobStore.Push(ctx, authInfo.Realm, imageBlobId, blobInfo, reader)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("tempFile", tempFile).Msg("PutFile")
		return err
	}
	logger.InfoCtx(ctx).Str("id", id.String()).Msg("PutFile")

	return nil
}
