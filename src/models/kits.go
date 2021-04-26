package models

//Kits é o modelo do preço por marketplace do produto
type Kits struct {
	ID    uint32 `json:"-"`
	ProID uint32 `json:"-"` //Código do produto que é o KIT
	//KitProID uint32 `json:"-"`
	UsrID  uint32 `json:"-"`
	KitSKU string `json:"sku"`
	Qty    uint16 `json:"qty"`
}

//Validators Valida a estrutura  Kits
func (kit *Kits) Validators() error {
	return translateError(validate.Struct(kit))
}
