package models

import "time"

type Image struct {
	Id            string
	Format        string
	Width         int
	Height        int
	OriginalPath  string
	ThumbnailPath string
	UploadedAt    time.Time
}
