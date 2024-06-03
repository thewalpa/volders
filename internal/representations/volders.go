package representations

type Folder struct {
	Common
	ParentID *string `json:"parent_id,omitempty"` // Use pointer for nullable foreign key
	Name     string  `json:"name"`
}

type File struct {
	Common
	FolderID    string `json:"folder_id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	Data        []byte `json:"data,omitempty"` // Optional: include only when needed
}
