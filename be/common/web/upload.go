package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
)

// TODO
// Please delete the files after work
// blob.TempFile.Delete(ctx)
func HandleFileUpload(ctx context.Context, r *http.Request) (out map[string][]blob.TempFile, err error) {

	// ctx, span := logger.Span(ctx, "web.HandleFileUpload", nil)
	// defer func() {
	// 	span.End(err)
	// }()

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		err = fmt.Errorf("upload error: %w", err)
		// logger.ErrorCtx(ctx, err).Msg("web.ParseMultipartForm")
		return nil, fmt.Errorf("web.ParseMultipartForm %w", err)
	}
	if log.Log().IsTrace() {
		fmt.Printf("Form: %+v \n", r.Form)
		fmt.Printf("Header: %+v \n", r.Header)
		fmt.Printf("MultipartForm.File: %+v \n", r.MultipartForm.File)
		fmt.Printf("MultipartForm.Value: %+v \n", r.MultipartForm.Value)

		for v := range r.MultipartForm.File {
			putValue := r.MultipartForm.File[v]
			fmt.Printf("\tFile[%s]: %+v \n", v, putValue)
			fmt.Printf("\tFile[%s]: %#v \n", v, putValue)
		}

		for v := range r.MultipartForm.Value {
			putValue := r.MultipartForm.Value[v]
			fmt.Printf("\tValue[%s]: %+v \n", v, putValue)
			fmt.Printf("\tValue[%s]: %#v \n", v, putValue)
		}

	}

	out = make(map[string][]blob.TempFile)

	for v := range r.MultipartForm.File {
		fileHeaders := r.MultipartForm.File[v]
		out[v] = []blob.TempFile{}
		// if logger.IsTrace() {
		// fmt.Printf("MultipartForm.File:[%s] \n", v)
		// }

		for i := range fileHeaders {
			fileHeader := fileHeaders[i]

			if log.Log().IsTrace() {
				fmt.Printf("\t[%s]:    Filename: %+v \n", v, fileHeader.Filename)
				fmt.Printf("\t[%s]:     Headers: %+v \n", v, fileHeader.Header)
				fmt.Printf("\t[%s]:Content-Type: %+v \n", v, fileHeader.Header.Get("Content-Type"))
				fmt.Printf("\t[%s]         Size: %+v \n", v, fileHeader.Size)
			}

			fileReader, err := fileHeader.Open()
			if err != nil {
				err = fmt.Errorf("upload error: %w", err)
				// logger.ErrorCtx(ctx, err).Msg("rest.FormFile")
				return nil, fmt.Errorf("web.fileHeader.Open[%s] %w", fileHeader.Filename, err)
			}
			defer fileReader.Close()

			tmpFile, err := os.CreateTemp("", ".bin")
			if err != nil {
				// logger.ErrorCtx(ctx, err).Msg("tmp file")
				return nil, fmt.Errorf("web.CreateTemp[%s] %w", fileHeader.Filename, err)
			}
			defer tmpFile.Close()
			// logger.DebugCtx(ctx).Str("tmpFile", tmpFile.Name()).Send()

			_, err = io.Copy(tmpFile, fileReader)
			if err != nil {
				// logger.ErrorCtx(ctx, err).Msg("copy to tmp file")
				return nil, fmt.Errorf("web.TempFile.Copy[%s] %w", fileHeader.Filename, err)
			}

			formFile := blob.TempFile{
				FileName:     blob.FileNameFilter(fileHeader.Filename),
				ContentType:  fileHeader.Header.Get("Content-Type"),
				TempFileName: tmpFile.Name(),
			}

			// logger.DebugCtx(ctx).
			// 	Any("info", tmpFile).
			// 	Msg("Upload")
			out[v] = append(out[v], formFile)
		}
	}
	return out, err

}
