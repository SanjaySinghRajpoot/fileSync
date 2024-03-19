package models

type DownloadPayload struct {
	FileName string `json:"file_name"`
	Version  int    `json:"version"`
	UserID   int    `json:"user_id"`
}

type VersionPayload struct {
	FileName string `json:"file_name"`
	UserID   int    `json:"user_id"`
}

type Chunk struct {
	Chunk string `json:"chunk"`
	Order int    `json:"order"`
}

type RecordPayload struct {
	UserID   int     `json:"user_id"`
	FileName string  `json:"filename"`
	Version  int     `json:"version"`
	Chunks   []Chunk `json:"chunks"`
}
