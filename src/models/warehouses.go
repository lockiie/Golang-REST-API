package models

import (
	"eco/src/types"
	"errors"
)

type Warehouses struct {
	// ID          uint32 `json:"-"`
	Name    string `json:"name"`
	ZipCode string `json:"zipcode"`
	Status  bool   `json:"active"`
	//	Register    time.Time `json:"-"`
	//	Update      time.Time `json:"-"`
	Code  string `json:"id"`
	UsrID uint32 `json:"-"`
	ComID uint32 `json:"-"`
}

//Validators Valida a estrutura Brands
func (w *Warehouses) Validators() error {
	if w.Name == types.EmptyStr {
		return errors.New("name é requerido!")
	}
	if len(w.Name) > 30 {
		return errors.New("name não pode ter mais de 30 caracteres!")
	}

	if len(w.Code) > 12 {
		return errors.New("code não pode ter mais de 12 caracteres!")
	}
	return nil
}
