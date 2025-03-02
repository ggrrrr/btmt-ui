package image

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
)

type (
	ResizedImage struct {
		BlobMd     blob.BlobMD
		ReadCloser TmpReadCloser
	}

	TmpReadCloser struct {
		tmpFile *os.File
	}
)

var _ (io.ReadCloser) = (*TmpReadCloser)(nil)

func (tmp *TmpReadCloser) Read(p []byte) (n int, err error) {
	return tmp.tmpFile.Read(p)
}

func (tmp *TmpReadCloser) Close() error {
	defer func() {
		err := os.Remove(tmp.tmpFile.Name())
		if err != nil {

			log.Log().Warn(err, "unable to remove tmp file",
				log.WithString("fmpFile", tmp.tmpFile.Name()))

		}

	}()

	err := tmp.tmpFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func HeadImage(name string) (blob.MDImageInfo, error) {
	reader, err := os.Open(name)
	if err != nil {
		return blob.MDImageInfo{}, fmt.Errorf("headImage[%s] %w", name, err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return blob.MDImageInfo{}, fmt.Errorf("headImage[%s] %w", name, err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	return blob.MDImageInfo{
		Width:  int64(w),
		Height: int64(h),
	}, nil
}

// go get -u github.com/disintegration/imaging
// go get -u github.com/rwcarlsen/goexif/exif
func ResizeImage(ctx context.Context, toHeight int, from blob.BlobReader) (*ResizedImage, error) {
	var err error

	if toHeight < 50 {
		toHeight = 256 / 2
	}

	defer from.ReadCloser.Close()

	img, imgFormat, err := image.Decode(from.ReadCloser)
	if err != nil {
		return nil, fmt.Errorf("ResizeImage.Decode[%s] %w", from.Blob.MD.Name, err)
	}

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

	newImageFile, err := os.CreateTemp("", ".resize.bin")
	if err != nil {
		return nil, fmt.Errorf("ResizeImage.CreateTemp[%s] %w", from.Blob.MD.Name, err)
	}

	switch imgFormat {
	case "jpeg", "jpg":
		err = jpeg.Encode(newImageFile, resizedImage, nil)
	case "png":
		err = png.Encode(newImageFile, resizedImage)
	default:
		os.Remove(newImageFile.Name())
		return nil, fmt.Errorf("unsupported file format: %s", imgFormat)
	}

	if err != nil {
		os.Remove(newImageFile.Name())
		return nil, fmt.Errorf("img.Encode[%s/%s] %w", imgFormat, from.Blob.MD.Name, err)
	}

	stat, err := newImageFile.Stat()
	if err != nil {
		os.Remove(newImageFile.Name())
		return nil, fmt.Errorf("%s.Encode error %w", imgFormat, err)
	}

	tmpFile, err := os.Open(newImageFile.Name())
	if err != nil {
		os.Remove(newImageFile.Name())
		return nil, fmt.Errorf("cant open tmp file[%s] %w", newImageFile.Name(), err)
	}

	out := ResizedImage{
		BlobMd: from.Blob.MD,
		// TmpFile: newImageFile.Name(),
		ReadCloser: TmpReadCloser{
			tmpFile: tmpFile,
		},
	}
	out.BlobMd.ContentType = fmt.Sprintf("%s/%s", "image", imgFormat)
	out.BlobMd.ContentLength = stat.Size()
	out.BlobMd.ImageInfo.Height = int64(toHeight)
	out.BlobMd.ImageInfo.Width = int64(toWidth)

	return &out, nil
}
