package domain

type Permissions struct {
	Comment                 bool          `json:"comment"`
	CommentGroups           CommentGroups `json:"commentGroups"`
	Copy                    bool          `json:"copy"`
	DeleteCommentAuthorOnly bool          `json:"deleteCommentAuthorOnly"`
	Download                bool          `json:"download"`
	Edit                    bool          `json:"edit"`
	EditCommentAuthorOnly   bool          `json:"editCommentAuthorOnly"`
	FillForms               bool          `json:"fillForms"`
	ModifyContentControl    bool          `json:"modifyContentControl"`
	ModifyFilter            bool          `json:"modifyFilter"`
	Print                   bool          `json:"print"`
	Review                  bool          `json:"review"`
	ReviewGroups            []string      `json:"reviewGroups"`
}

type CommentGroups struct {
	Edit   []string `json:"edit"`
	Remove []string `json:"remove"`
	View   string   `json:"view"`
}

type User struct {
	Group string `json:"group"`
	ID    uint64 `json:"id,string"`
	Name  string `json:"name"`
}

// InsertImage
type InsertImage struct {
	C      string  `json:"c"` // add
	Images []Image `json:"images"`
}
type Image struct {
	FileType string `json:"fileType"`
	Url      string `json:"url"`
}

// SetHistoryData
type SetHistoryData struct {
	ChangesUrl string   `json:"changesUrl"`
	Key        string   `json:"key"`
	Url        string   `json:"url"`
	Version    int      `json:"version"`
	Previous   Previous `json:"previous"`
}
type Previous struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

// SetMailMergeRecipients
type SetMailMergeRecipients struct {
	FileType string `json:"fileType"`
	Url      string `json:"url"`
}

// SetReferenceData
type SetReferenceData struct {
	FileType      string        `json:"fileType"`
	Path          string        `json:"path"`
	ReferenceData ReferenceData `json:"referenceData"`
	Url           string        `json:"url"`
}
type ReferenceData struct {
	FileKey    string `json:"fileKey"`
	InstanceID string `json:"instanceId"`
}

// SetRevisedFile
type SetRevisedFile struct {
	FileType string `json:"fileType"`
	Url      string `json:"url"`
}
