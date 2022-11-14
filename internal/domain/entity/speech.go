package entity

type Speech struct {
	Content     string `json:"content"`
	ParagraphID uint64 `json:"paragraph_id,omitempty"`
	OrderNum    uint64 `json:"order_num,omitempty"`
}
