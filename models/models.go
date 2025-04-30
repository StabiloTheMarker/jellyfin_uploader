package models

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	gorm.Model
	Filepath        string
	Uploaded        bool
	UploadProcessID uint
	UploadedAt      time.Time
}
type UploadProcess struct {
	gorm.Model
	Files []File
}
