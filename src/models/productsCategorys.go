package models

//ProductsCategorys é um modelo de categorias produtos
type ProductsCategorys struct {
	ID    uint32 `json:"-"`
	Code  string `json:"id" validate:"required" max:"50"`
	ProID uint32 `json:"-"`
}

//Validators Valida a estrutura de uma categoria com produto
func (ctp *ProductsCategorys) Validators() error {
	return translateError(validate.Struct(ctp))
}
