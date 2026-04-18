package dto

type ProductData struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Categorie *CategorieData `json:"categorie"`
	Merek     string         `json:"merek"`
}

type ProductCreated struct {
	Name        string `json:"name" validate:"required"`
	CategorieId string `json:"categorie_id" validate:"required"`
	Merek       string `json:"merek" validate:"required"`
}

type ProductUpdated struct {
	ID          string `json:"-"`
	Name        string `json:"name" validate:"required"`
	CategorieId string `json:"categorie_id" validate:"required"`
	Merek       string `json:"merek" validate:"required"`
}
