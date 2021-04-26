package models

//ProductsStocks é o modelo do estoque do produto
type ProductsStocks struct {
	//ID         uint32 `json:"-"`
	Qty        uint16 `json:"qty" validate:"required"`
	QtyBooking uint16 `json:"booking_qty"`
	Warehouse  string `json:"warehouse_id"`
	UsrID      uint32 `json:"-"`
	ProID      uint32 `json:"-"`
}

//Validators Valida a estrutura estoque
func (pro *ProductsStocks) Validators() error {
	return translateError(validate.Struct(pro))
}
