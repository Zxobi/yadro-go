package database

type RecordMap map[int]Entity

type Entity struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

func NewEntity(url string, keywords []string) Entity {
	return Entity{Url: url, Keywords: keywords}
}
