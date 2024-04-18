package database

type RecordMap map[int]Record
type IndexMap map[string][]int

type Record struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

func NewRecord(url string, keywords []string) *Record {
	return &Record{Url: url, Keywords: keywords}
}
