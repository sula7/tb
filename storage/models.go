package storage

type (
	NewsItem struct {
		URL        string
		TagContent string
	}

	TagContent struct {
		ID         int64  `json:"id"`
		TagContent string `json:"tag_content"`
		URL        string `json:"url"`
		CreatedAt  string `json:"created_at"`
	}

	TagContentFilter struct {
		ContentKeyword string
		URL            string
		CreatedAtFrom  string
		CreatedAtTo    string
	}
)
