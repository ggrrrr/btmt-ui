package app

import (
	"context"
	"fmt"
	"image"
	"os"

	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	appError "github.com/ggrrrr/btmt-ui/be/common/app"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func headImage(name string) (blob.MDImageInfo, error) {
	reader, err := os.Open(name)
	if err != nil {
		return blob.MDImageInfo{}, fmt.Errorf("headImage[%s] %w", name, err)
	}
	m, _, err := image.Decode(reader)
	if err != nil {
		return blob.MDImageInfo{}, fmt.Errorf("headImage[%s] %w", name, err)
	}
	defer reader.Close()

	bounds := m.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	return blob.MDImageInfo{
		Width:  int64(w),
		Height: int64(h),
	}, nil
}

// go get -u github.com/disintegration/imaging
// go get -u github.com/rwcarlsen/goexif/exif
func resizeImage(from blob.BlobReader) (string, error) {

	defer from.ReadCloser.Close()

	img, imgFormat, err := image.Decode(from.ReadCloser)
	fmt.Printf("image.Decode %s \n", imgFormat)
	if err != nil {
		return "", fmt.Errorf("resizeImage.Decode[%s] %w", from.Blob.MD.Name, err)
	}
	toHeight := 256 / 2

	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	resizeFactor := float32(imgHeight) / float32(toHeight)
	ratio := float32(imgWidth) / float32(imgHeight)
	toWidth := int(float32(toHeight) * ratio)

	resizedImage := image.NewNRGBA(image.Rect(0, 0, toWidth, toHeight))
	var imgX, imgY int
	var imgColor color.Color
	for x := 0; x < toWidth; x++ {
		for y := 0; y < toHeight; y++ {
			imgX = int(resizeFactor*float32(x) + 0.5)
			imgY = int(resizeFactor*float32(y) + 0.5)
			imgColor = img.At(imgX, imgY)
			resizedImage.Set(x, y, imgColor)
		}
	}

	fmt.Printf("from H / W %5d / %5d \n", imgHeight, imgWidth)
	fmt.Printf("to         %5d / %5d \n", toHeight, toWidth)

	newImageFile, err := os.CreateTemp("", ".resize.bin")
	if err != nil {
		return "", fmt.Errorf("resizeImage.CreateTemp[%s] %w", from.Blob.MD.Name, err)
	}

	switch imgFormat {
	case "jpeg":
		err = jpeg.Encode(newImageFile, resizedImage, nil)
	case "png":
		err = png.Encode(newImageFile, resizedImage)
	default:
		os.Remove(newImageFile.Name())
		return "", fmt.Errorf("unsupported file format: %s", imgFormat)
	}

	if err != nil {
		os.Remove(newImageFile.Name())
		return "", fmt.Errorf("%s.Encode error %w", imgFormat, err)
	}

	defer newImageFile.Close()
	return newImageFile.Name(), err
}

func (a *App) GetResizedImage(ctx context.Context, fileId string) (*ddd.FileWriterTo, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "GetResizedImage", nil, logger.AttributeString("fileId", fileId))
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

	res, err := a.blobStore.Fetch(ctx, "localhost", imageBlobId)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Fetch")
		return nil, err
	}

	logger.InfoCtx(ctx).
		Any("info", res).
		Any("id", fileId).
		Msg("Fetch")

	resizedFile, err := resizeImage(res)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("resizeImage")
		return nil, err
	}

	tmpFile, err := os.Open(resizedFile)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Open")
		return nil, err
	}

	logger.InfoCtx(ctx).
		Str("resizedFile", resizedFile).
		Msg("Open")

	defer func() {
		tmpFile.Close()
		// os.Remove(tmpFile.Name())
	}()

	return &ddd.FileWriterTo{
			ContentType: res.Blob.MD.ContentType,
			Version:     res.Blob.Id.Version(),
			Name:        res.Blob.MD.Name,
			WriterTo: &filePipe{
				reader: tmpFile,
			},
		},
		nil

}

func (a *App) GetImage(ctx context.Context, fileId string, maxWight int) (*ddd.FileWriterTo, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "GetImage", nil, logger.AttributeString("fileId", fileId))
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
		logger.AttributeString("FileName", tempFile.FileName),
		logger.AttributeString("ContentType", tempFile.ContentType),
		logger.AttributeString("TempFileName", tempFile.TempFileName),
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

	id, err := a.blobStore.Push(ctx, authInfo.Realm, imageBlobId, blobInfo, reader)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("tempFile", tempFile).Msg("PutFile")
		return err
	}
	logger.InfoCtx(ctx).Str("id", id.String()).Msg("PutFile")

	return nil
}
