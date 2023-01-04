package entity

type Absent struct {
	ID          uint64 `json:"id,omitempty"`
	Pseudo      string `json:"pseudo,omitempty"`
	Done        bool   `json:"done,omitempty"`
	ChapterID   uint64 `json:"c_id,omitempty"`
	ParagraphID uint64 `json:"p_id,omitempty"`
}

func (a *Absent) Validate() bool {
	if a.ChapterID <= 0 || a.Pseudo == "" || a.ParagraphID <= 0 {
		return false
	}
	return true
}
