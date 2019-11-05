package models

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
