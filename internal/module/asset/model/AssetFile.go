package model

type AssetFile struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	CategoryID uint64 `json:"categoryId,omitempty,string" gorm:"index;comment:分类ID"`
}
