package search

type Document struct {
	ID           uint64   `json:"id,string"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Route        string   `json:"route"`
	RelatedUsers []string `json:"relatedUsers"`
	CreateTime   int64    `json:"createTime"`
	UpdateTime   int64    `json:"updateTime"`
}
