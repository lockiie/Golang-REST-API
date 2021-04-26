package models

//Products é o modelo para a especificação
type Products struct {
	ID             uint32                   `json:"-"`
	SKU            string                   `json:"sku" max:"50"`
	Title          string                   `json:"title" validate:"required" max:"200"`
	Status         bool                     `json:"active"`
	ComID          uint32                   `json:"-"`
	BndCode        string                   `json:"brand_id" validate:"required"`
	Description    string                   `json:"description" validate:"required" max:"1000"`
	Weigth         float32                  `json:"weigth"  max:"4"`
	Size           float32                  `json:"size"  max:"4"`
	Width          float32                  `json:"width"  max:"4"`
	Heigth         float32                  `json:"heigth"  max:"4"`
	Barcode        string                   `json:"barcode"  max:"13"`
	Crossdocking   uint16                   `json:"crossdocking"`
	UsrID          uint32                   `json:"-"`
	PdtID          uint8                    `json:"type_product"`
	Variation      *uint32                  `json:"-"`
	Stock          ProductsStocks           `json:"stock" validate:"required"`
	Category       uint32                   `json:"category_id,omitempty"`
	Specifications []ProductsSpecifications `json:"specifications,omitempty"`
	Prices         []ProductsPrices         `json:"prices,omitempty" validate:"required"`
	Kits           []Kits                   `json:"kits,omitempty"`
	Variations     []ProductsVariations     `json:"variations,omitempty"`
	Images         []ProductsImages         `json:"images,omitempty" validate:"required"`
}

//Validators Valida a estrutura de especificação
func (pro *Products) Validators() error {

	return translateError(validate.Struct(pro))
	//if pro.SKU == types.EmptyStr {
	// 	return errors.New("sku" + types.MsgRequerid)
	// }
	// if len(pro.SKU) > 50 {

	// }
	// if len(pro.Title) > 200 {

	// }
	// if pro.Size > 9999.9 {

	// }
	// if pro.Weigth > 9999.9 {

	// }
	// if pro.Heigth > 9999.9 {

	// }
	// if pro.Width > 9999.9 {

	// }
	// if len(pro.Barcode) > 13 {

	// }
	// if pro.Categorys == nil {

	// }
	// if pro.Specifications == nil {

	// }
	// if pro.Prices == nil {

	// }
	// if pro.Images == nil {

	// }
	// return nil
}
