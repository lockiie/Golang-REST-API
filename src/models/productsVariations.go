package models

//ProductsVariations estrutura de um produto variado
type ProductsVariations struct {
	ID             uint32                   `json:"-"`
	SKU            string                   `json:"sku" validate:"required" max:"50"`
	Title          string                   `json:"title" validate:"required" max:"200"`
	Status         bool                     `json:"active"`
	ComID          uint32                   `json:"-"`
	Description    string                   `json:"description" validate:"required" max:"1000"`
	Barcode        string                   `json:"barcode"  max:"13"`
	Crossdocking   uint16                   `json:"crossdocking"`
	UsrID          uint32                   `json:"-"`
	Variation      uint32                   `json:"-"`
	BndID          uint32                   `json:"-"`
	Stock          ProductsStocks           `json:"stock"`
	Specifications []ProductsSpecifications `json:"specifications,omitempty"`
	Prices         []ProductsPrices         `json:"prices,omitempty" validate:"required"`
	Images         []ProductsImages         `json:"images,omitempty" validate:"required"`
}

//Validators Valida a estrutura produto variado
func (pro *ProductsVariations) Validators() error {
	return translateError(validate.Struct(pro))
}
