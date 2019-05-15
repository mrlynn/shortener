package models

type Shortener struct {
	ID           int64  `json:"id" bson:"_id,omitempty"`
	OriginalURL  string `json:"originalurl"`
	GeneratedURL string `json:"generatedurl"`
	Visited      bool   `json:"visited"`
	Count        int64  `json:"count"`
}
