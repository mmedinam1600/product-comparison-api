package domain

// Item representa un producto en el catálogo
type Item struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	ImageURL       string                 `json:"image_url"`
	Description    string                 `json:"description"`
	Price          float64                `json:"price"`
	Rating         float64                `json:"rating"`
	Specifications map[string]interface{} `json:"specifications"`
}
