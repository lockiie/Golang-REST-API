package models

//ProductsImages é o modelo para a images do produto
type ProductsImages struct {
	ID    uint32 `json:"-"`
	Order uint8  `json:"order" validate:"required" max:"1"`
	URI   string `json:"uri" validate:"required" max:"200"`
	ProID uint32 `json:"-"`
}

//Validators Valida a estrutura Brands
func (pri *ProductsImages) Validators() error {
	return translateError(validate.Struct(pri))
}
