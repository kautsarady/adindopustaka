package model

// Book .
type Book struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	GramedURL   string `json:"gramed_url,omitempty"`
	Description string `json:"description"`
	Authors     []Item `json:"authors,omitempty"`
	Categories  []Item `json:"categories,omitempty"`
	Tags        []Item `json:"tags,omitempty"`
}

// Item .
type Item struct {
	ID     int    `json:"id,omitempty"`
	BookID int    `json:"-"`
	Name   string `json:"name,omitempty"`
	Books  []Book `json:"books,omitempty"`
}
