package models

import (
	"eco/src/types"
	"errors"
)

type Brands struct {
	ID          uint32 `json:"-"`
	Name        string `json:"name" validate:"required" max:"100"`
	Description string `json:"description"`
	Status      bool   `json:"active"`
	//	Register    time.Time `json:"-"`
	//	Update      time.Time `json:"-"`
	Code  string `json:"id"`
	UsrID uint32 `json:"-"`
	ComID uint32 `json:"-"`
}

//Validators Valida a estrutura Brands
func (b *Brands) Validators() error {
	if b.Name == types.EmptyStr {
		return errors.New("name é requerido!")
	}
	if len(b.Name) > 100 {
		return errors.New("name não pode ter mais de 100 caracteres!")
	}

	if b.Description == types.EmptyStr {
		return errors.New("description não pode ter mais de 500 caracteres!")
	}

	if len(b.Code) > 12 {
		return errors.New("code não pode ter mais de 12 caracteres!")
	}

	return nil
}
