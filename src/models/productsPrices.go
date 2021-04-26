package models

//ProductsPrices é o modelo do preço por marketplace do produto
type ProductsPrices struct {
	ID         uint32  `json:"-"`
	ProID      uint32  `json:"-"`
	Price      float32 `json:"price" validate:"required"`
	PricePromo float32 `json:"promotional_price" validate:"required"`
	Status     bool    `json:"active"`
	MkcID      uint32  `json:"marketplace_id,omitempty"`
	UsrID      uint32  `json:"-"`
}

//Validators Valida a estrutura  ProductsPrices
func (prp *ProductsPrices) Validators() error {
	return translateError(validate.Struct(prp))
}
