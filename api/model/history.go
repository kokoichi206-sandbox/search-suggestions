package model

// OpenSearch にリクエストする JSON の構造体。
type HistorySearchRequest struct {
	ID     string `json:"id"`
	Params Params `json:"params"`
}
type Params struct {
	Query string `json:"query"`
}

// OpenSearch から受け取る JSON の構造体。
// 適当に json to go で生成したものを元にしている。
type HistorySearchResult struct {
	Took         int          `json:"took"`
	TimedOut     bool         `json:"timed_out"`
	Shards       Shards       `json:"_shards"`
	Hits         Hits         `json:"hits"`
	Aggregations Aggregations `json:"aggregations"`
}
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}
type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
type Hits struct {
	Total    Total         `json:"total"`
	MaxScore interface{}   `json:"max_score"`
	Hits     []interface{} `json:"hits"`
}
type Buckets struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}
type Keywords struct {
	DocCountErrorUpperBound int       `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int       `json:"sum_other_doc_count"`
	Buckets                 []Buckets `json:"buckets"`
}
type Aggregations struct {
	Keywords Keywords `json:"keywords"`
}

// サジェストエンドポイント用 HTTP レスポンスの構造体。
type SuggestResponse struct {
	Suggestions []string `json:"suggestions"`
}
