package domain

const (
	DocumentTypeWord  = "word"
	DocumentTypeCell  = "cell"
	DocumentTypeSlide = "slide"
)
const (
	EditConfigTypeDesktop  = "desktop"
	EditConfigTypeMobile   = "mobile"
	EditConfigTypeEmbedded = "embedded"
)
const (
	EditConfigModeEdit = "edit"
	EditConfigModeView = "view"
)
const (
	EditConfigLangZhCn = "zh-CN"
)

type EditConfig struct {
	Document     Document `json:"document"`
	DocumentType string   `json:"documentType"`
	Token        string   `json:"token"`
	Type         string   `json:"type"`
	EditorConfig Editor   `json:"editorConfig"`
}

func (e EditConfig) Valid() error {
	return nil
}

type Document struct {
	FileType    string      `json:"fileType"`
	Key         string      `json:"key"`
	Title       string      `json:"title"`
	Url         string      `json:"url"`
	Permissions Permissions `json:"permissions"`
}

type Editor struct {
	CallbackUrl string `json:"callbackUrl"`
	Mode        string `json:"mode"`
	Lang        string `json:"lang"`
	User        User   `json:"user"`
}
