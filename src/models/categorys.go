package models

import (
	"eco/src/types"
	"errors"
)

//Categorys é o modelo para a categoria
type Categorys struct {
	CTT_ID uint32 `json:"internal_id"`
	ComID  uint32 `json:"-"`
	UsrID  uint32 `json:"-"`
	Code   string `json:"nick"`
	Title  string `json:"title"`
	Father uint32 `json:"father_id,omitempty" `
}

//Validators Valida a estrutura de categoria
func (c *Categorys) Validators() error {
	if c.Code == types.EmptyStr {
		return errors.New("nick é requerido!")
	}

	if len(c.Code) > 12 {
		return errors.New("nick não pode ter mais de 12 caracteres!")
	}

	if c.CTT_ID == 0 {
		return errors.New("internal_id é requerido!")
	}
	return nil
}
