package models

//ProductsSpecifications é o modelo para a images do produto
type ProductsSpecifications struct {
	Value string `json:"value" validate:"required" max:"50"`
	Code  string `json:"id" validate:"required"  max:"50"`
	Title string `json:"title,omitempty"`
	//SptID uint32 `json:"-"`
	ProID uint32 `json:"-"`
}

//Validators Valida a estrutura Brands
func (spe *ProductsSpecifications) Validators() error {
	return translateError(validate.Struct(spe))
}
