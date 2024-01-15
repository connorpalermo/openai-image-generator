package model

type ImageRequestLocal struct {
	Prompt   string `json:"prompt" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}

type ImageRequestS3 struct {
	Prompt     string `json:"prompt" binding:"required"`
	FileName   string `json:"fileName" binding:"required"`
	BucketName string `json:"bucketName" binding:"required"`
}

type DownloadImageS3 struct {
	Item     string `json:"item" binding:"required"`
	Bucket   string `json:"bucket" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}
