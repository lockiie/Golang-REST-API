package repositories

import (
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//InsertQrySpecification insere uma especificação com um select
func (db RepoProducts) InsertQrySpecification(s *models.ProductsSpecifications, ComID uint32) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_SPECIFICATIONS_PRODUCTS (SPE_VALUE, SPT_ID,  PRO_ID)
		 SELECT :SPE_VALUE, SPT_ID, :PRO_ID FROM SPECIFICATIONS_TITLES WHERE SPT_CODE = :SPT_CODE AND COM_ID = :COM_ID`,
		s.Value, s.ProID, s.Code, ComID,
	)
	return err
}

//InsertTwoQrySpecification insere uma especificação com um 2 selects
func (db RepoProducts) InsertTwoQrySpecification(s *models.ProductsSpecifications, ComID uint32, SKU string) error {
	_, err := db.conn.ExecContext(
		*db.ctx,
		`INSERT INTO TB_SPECIFICATIONS_PRODUCTS (SPE_VALUE, SPT_ID,  PRO_ID)
		VALUES(:0, (SELECT SPT_ID FROM SPECIFICATIONS_TITLES WHERE SPT_CODE = :1 AND COM_ID = :2),
	        (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :3 AND COM_ID = :4))`,
		s.Value, s.Code, ComID, SKU, ComID,
	)

	return err
}

//UpdateTwoQrySpecification altera uma especificação com um 2 selects
func (db RepoProducts) UpdateTwoQrySpecification(s *models.ProductsSpecifications, ComID uint32, SKU string) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`UPDATE TB_SPECIFICATIONS_PRODUCTS SET SPE_VALUE = :0  
		WHERE SPT_ID = (SELECT SPT_ID FROM SPECIFICATIONS_TITLES WHERE SPT_CODE = :1 AND COM_ID = :2)
		AND PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :3 AND COM_ID = :4)`,
		s.Value, s.Code, ComID, SKU, ComID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//DeleteTwoQrySpecification exclui uma especificação de um produto com dois selects
func (db RepoProducts) DeleteTwoQrySpecification(SKU string, Code string, comID uint32) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`DELETE FROM TB_SPECIFICATIONS_PRODUCTS 
		WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :0 AND COM_ID = :1)
		  AND SPT_ID = (SELECT SPT_ID FROM SPECIFICATIONS_TITLES WHERE SPT_CODE = :2 AND COM_ID = :3)`,
		SKU, comID, Code, comID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//QueryBySKUSpecifications retorna um array de specifications com base no sku do produto
func (db RepoProducts) QueryBySKUSpecifications(SKU string, comID uint32) (*[]models.Specifications, error) {
	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT S.SPT_TITLE, S.SPT_CODE, S.SPT_STATUS, SP.SPE_VALUE
		 FROM  PRODUCTS P, SPECIFICATIONS_TITLES S, SPECIFICATIONS_PRODUCTS SP
	     WHERE SP.SPT_ID = S.SPT_ID
		   AND SP.PRO_ID = P.PRO_ID
		   AND P.PRO_SKU = :PRO_SKU
		   AND P.COM_ID = :COM_ID`, SKU, comID,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var specifications []models.Specifications

	for rows.Next() {
		var specification models.Specifications
		rows.Scan(&specification.Title, &specification.Code, &specification.Status, &specification.Value)
		specifications = append(specifications, specification)
	}
	return &specifications, nil
}
