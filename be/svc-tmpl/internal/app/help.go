package app

import "io"

type filePipe struct {
	reader io.ReadCloser
}

func (f *filePipe) WriteTo(w io.Writer) (int64, error) {
	defer f.reader.Close()
	return io.Copy(w, f.reader)
}
