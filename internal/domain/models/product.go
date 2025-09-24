package models

type Product struct {
	ID          string
	Name        string
	ImageURL    string
	Description string
	Price       float64
	Rating      float64
	Specs       map[string]string
}
