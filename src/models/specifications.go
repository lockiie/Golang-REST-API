package models

//Specifications é o modelo para a especificação
type Specifications struct {
	ID     uint32 `json:"-"`
	Title  string `json:"title" validate:"required" max:"50"`
	Code   string `json:"id" validate:"required" max:"50"`
	Value  string `json:"value,omitempty"`
	Status bool   `json:"active"`
	UsrID  uint32 `json:"-"`
	ComID  uint32 `json:"-"`
}

//Validators Valida a estrutura de especificação
func (spt *Specifications) Validators() error {
	return translateError(validate.Struct(spt))
}
