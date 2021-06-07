package models

import "time"

type PRLists struct {
	Page    int `json:"page"`
	Pagelen int `json:"pagelen"`
	Size    int `json:"size"`
	Values  []struct {
		ID    int `json:"id"`
		CreatedOn         time.Time   `json:"created_on"`
		State   string `json:"state"`
		Title     string    `json:"title"`
		UpdatedOn time.Time `json:"updated_on"`
		Author struct {
			DisplayName string `json:"display_name"`
		} `json:"author"`
		Destination       struct {
			Branch struct {
				Name string `json:"name"`
			} `json:"branch"`
		} `json:"destination"`
		Source      struct {
			Branch struct {
				Name string `json:"name"`
			} `json:"branch"`
		} `json:"source"`

	} `json:"values"`
}
