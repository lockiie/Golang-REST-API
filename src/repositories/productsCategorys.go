package repositories

import (
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//InsertQryCategory insere uma categoria com um select
func (db RepoProducts) InsertQryCategory(c *models.ProductsCategorys, comID uint32) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_CATEGORYS_PRODUCTS (PRO_ID, CTT_ID) 
		SELECT :PRO_ID, CTT_ID FROM CATEGORYS_TITLES WHERE CTT_CODE = :CTT_CODE AND COM_ID = :COM_ID`,
		c.ProID, c.Code, comID,
	)
	return err
}

//InsertTwoQryCategory insere uma categoria com 2 select
func (db RepoProducts) InsertTwoQryCategory(SKU string, Code string, comID uint32) error {
	_, err := db.conn.ExecContext(
		*db.ctx,
		`INSERT INTO TB_CATEGORYS_PRODUCTS (PRO_ID, CTT_ID) 
		 VALUES((SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :1 AND COM_ID = :2),
		 (SELECT CTT_ID FROM CATEGORYS_TITLES WHERE CTT_CODE = :CTT_CODE AND COM_ID = :COM_ID))`,
		SKU, comID, Code, comID,
	)
	return err
}

//DeleteTwoQryCategory exclui uma categoria de um produto com dois selects
func (db RepoProducts) DeleteTwoQryCategory(SKU string, Code string, comID uint32) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`DELETE FROM TB_CATEGORYS_PRODUCTS 
		WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :0 AND COM_ID = :1)
		  AND CTT_ID = (SELECT CTT_ID FROM CATEGORYS_TITLES WHERE CTT_CODE = :2 AND COM_ID = :3)`,
		SKU, comID, Code, comID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

// //QueryBySKUCategorys retorna um array de categorias do produto especificado
// func (db RepoProducts) QueryBySKUCategorys(SKU string, comID uint32) (*[]models.Categorys, error) {

// 	rows, err := db.conn.QueryContext(
// 		*db.ctx,
// 		`SELECT C.CTT_ID, C.CTT_TITLE, C.CTT_CODE, C.CTT_FATHER
// 		FROM PRODUCTS P, CATEGORYS_PRODUCTS CP, CATEGORYS_TITLES C
// 		WHERE CP.PRO_ID = P.PRO_ID
// 		  AND C.CTT_ID = CP.CTT_ID
// 		  AND P.PRO_SKU = :PRO_SKU
// 		  AND P.COM_ID = :COM_ID`, SKU, comID,
// 	)
// 	defer rows.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var categorys []models.Categorys

// 	for rows.Next() {
// 		var category models.Categorys
// 		rows.Scan(&category.ID, &category.Title, &category.Code, &category.Father)
// 		categorys = append(categorys, category)
// 	}
// 	return &categorys, nil
// }
