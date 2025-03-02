package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	appError "github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	appImage "github.com/ggrrrr/btmt-ui/be/common/image"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func (a *Application) GetResizedImage(ctx context.Context, fileId string) (*ddd.FileWriterTo, error) {
	var err error
	ctx, span := a.tracer.SpanWithAttributes(ctx, "GetResizedImage", slog.String("fileId", fileId))
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

	log.Log().InfoCtx(ctx, "GetResizedImage",
		slog.Any("info", fileId),
	)

	reader, err := a.blobStore.Fetch(ctx, "localhost", imageBlobId)
	if err != nil {
		return nil, err
	}

	res, err := appImage.ResizeImage(ctx, 0, reader)
	if err != nil {
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

func (a *Application) GetImage(ctx context.Context, fileId string, maxWight int) (*ddd.FileWriterTo, error) {
	var err error
	ctx, span := a.tracer.SpanWithAttributes(ctx, "GetImage", slog.String("fileId", fileId))
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

	res, err := a.blobStore.Fetch(ctx, "localhost", imageBlobId)
	if err != nil {
		return nil, err
	}

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

func (a *Application) ListImages(ctx context.Context) ([]ddd.ImageInfo, error) {
	var err error
	ctx, span := a.tracer.Span(ctx, "ListImages")
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return nil, err
	}

	res, err := a.blobStore.List(ctx, authInfo.Realm, a.imagesFolder)
	if err != nil {
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

	return out, nil
}

func (a *Application) PutImage(ctx context.Context, tempFile blob.TempFile) error {
	defer func() {
		tempFile.Delete(ctx)
	}()

	var err error
	ctx, span := a.tracer.SpanWithAttributes(ctx, "PutImage",
		slog.String("FileName", tempFile.FileName),
		slog.String("ContentType", tempFile.ContentType),
		slog.String("TempFileName", tempFile.TempFileName),
	)
	defer func() {
		span.End(err)
	}()

	imageBlobId, err := a.imagesFolder.SetIdVersionFromString(tempFile.FileName)
	if err != nil {
		return appError.BadRequestError("file name not supported", err)
	}

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
		return err
	}
	defer reader.Close()

	_, err = a.blobStore.Push(ctx, authInfo.Realm, imageBlobId, blobInfo, reader)
	if err != nil {
		return err
	}

	return nil
}
