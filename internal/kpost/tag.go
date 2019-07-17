package kpost

type Tag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	Ambiguous bool   `json:"ambiguous"`

	// 某个tag对post的重要性越高，它的TF-IDF值就越大
	TfIdf float64 `json:"tf_idf"`
}

type Tags []Tag

//func (p Tags) Len() int           { return len(p) }
//func (p Tags) Less(i, j int) bool { return p[i].TfIdf > p[j].TfIdf }
//func (p Tags) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
