package api

type (
	Response struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	TagContent struct {
		URL     string `json:"url"`
		TagName string `json:"tag_name"`
	}
)
