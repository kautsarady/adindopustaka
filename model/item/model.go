package item

type Author struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books []int  `json:"books,omitempty"`
}

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books []int  `json:"books,omitempty"`
}

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books []int  `json:"books,omitempty"`
}
