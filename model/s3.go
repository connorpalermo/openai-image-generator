package model

type ImageRequest struct {
	Prompt     string `json:"prompt" binding:"required"`
	FileName   string `json:"fileName" binding:"required"`
	BucketName string `json:"bucketName" binding:"required"`
}

type DownloadImage struct {
	Item     string `json:"item" binding:"required"`
	Bucket   string `json:"bucket" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}

type ImageRequestLocal struct {
	Prompt   string `json:"prompt" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}
