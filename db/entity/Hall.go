package entity

type Hall struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	Tools       string `json:"tools"`
	ImageLink   string `json:"image_link"`
}
