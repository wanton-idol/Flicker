package elasticsearchPkg

type ElasticResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	HitsObject struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore *float64 `json:"max_score"`
		Hits     []Hits   `json:"hits"`
	} `json:"hits"`
}

type Hits struct {
	Id     string      `json:"id"`
	Index  string      `json:"_index"`
	Score  string      `json:"_score"`
	Source UserProfile `json:"_source"`
	// Sort   []float64   `json:"sort"`
}
