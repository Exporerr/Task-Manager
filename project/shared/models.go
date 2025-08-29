package shared

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      bool
	Created_at  time.Time
}
type IDResponse struct {
	ID int64 `json:"id"`
}
type PostResponse struct {
	Message string `json:"message"`
	ID      int64  `json:"id"`
}
type DeleteOrUpdateResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}
type Config struct {
	DBService struct {
		URL string `yaml:"url"`
	} `yaml:"db_service"`
}
