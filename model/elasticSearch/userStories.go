package elasticsearchPkg

import "time"

type Values struct {
	Text     string    `json:"text"`
	Location []float64 `json:"location"`
}

type UserStories struct {
	ID            string    `json:"id"`
	UserProfileID int       `json:"user_profile_id"`
	Text          string    `json:"text"`
	MediaURL      string    `json:"media_url"`
	MediaType     string    `json:"media_type"`
	Location      []float64 `json:"location"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type StoriesResponse struct {
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
		MaxScore *float64      `json:"max_score"`
		Hits     []StoriesHits `json:"hits"`
	} `json:"hits"`
}

type StoriesHits struct {
	Id     string      `json:"id"`
	Index  string      `json:"_index"`
	Score  string      `json:"_score"`
	Source UserStories `json:"_source"`
	// Sort   []float64   `json:"sort"`
}
