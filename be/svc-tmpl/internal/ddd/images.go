package ddd

import (
	"io"
	"time"
)

type (
	FileWriterTo struct {
		ContentType string
		Version     string
		Name        string
		WriterTo    io.WriterTo
	}

	ImageInfo struct {
		Id          string
		Version     string
		FileName    string
		CreatedAt   time.Time
		ContentType string
		Width       int64
		Height      int64
		Size        int64
		Versions    []string
	}
)
