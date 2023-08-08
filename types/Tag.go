package types

type Tag struct {
	Id      int64  `json:"id" db:"id"`
	TagText string `json:"tagText" db:"tag_text"`
}
