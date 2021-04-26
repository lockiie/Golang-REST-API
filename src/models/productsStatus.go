package models

//ProductsStatus é o modelo do status do produto
type ProductsStatus struct {
	ID          uint32 `json:"-"`
	ProID       uint32 `json:"-"`
	MkpID       uint32 `json:"marketplace_id"`
	UsrID       uint32 `json:"-"`
	Description string `json:"description"`
	Status      bool   `json:"active"`
}

//Validators Valida a estrutura do status do produto
func (pro *ProductsStatus) Validators() error {
	return translateError(validate.Struct(pro))

}
