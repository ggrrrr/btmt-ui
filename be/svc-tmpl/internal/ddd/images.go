package ddd

import "time"

type (
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
