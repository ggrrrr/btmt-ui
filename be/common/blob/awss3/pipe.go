package awss3

import "io"

type blobCopier struct {
	reader io.ReadCloser
}

var _ (io.WriterTo) = (*blobCopier)(nil)

func (b *blobCopier) WriteTo(w io.Writer) (n int64, err error) {
	defer b.reader.Close()
	return io.Copy(w, b.reader)
}
