package repositories

import (
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//UpdateQryStock altera o estoque de um produto com um select
func (db RepoProducts) UpdateQryStock(s *models.ProductsStocks, SKU string, ComID uint32) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_PRODUCTS_STOCKS SET PRO_QTY = :0, PRO_UPDATE = SYSDATE, PRO_QTYBOOKING = :1, USR_ID = :2
		WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :3 AND COM_ID = :4)`,
		s.Qty, s.QtyBooking, s.UsrID, SKU, ComID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//UpdateStock altera o estoque de um produto
func (db RepoProducts) UpdateStock(s *models.ProductsStocks) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_PRODUCTS_STOCKS SET PRO_QTY = :0, PRO_UPDATE = SYSDATE, PRO_QTYBOOKING = :1, USR_ID = :2
		WHERE PRO_ID = :3`,
		s.Qty, s.QtyBooking, s.UsrID, s.ProID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//QueryBySKUStock retorna o stock do produto
func (db RepoProducts) QueryBySKUStock(SKU string, comID uint32) (models.ProductsStocks, error) {
	var stock models.ProductsStocks

	row := db.conn.QueryRowContext(
		*db.ctx,
		`SELECT S.PRO_QTY, S.PRO_QTYBOOKING
		 FROM PRODUCTS P, PRODUCTS_STOCKS S
		 WHERE P.PRO_ID = S.PRO_ID
		   AND P.PRO_SKU = :PRO_SKU
		   AND P.COM_ID = :COM_ID`, SKU, comID,
	)
	err := row.Scan(&stock.Qty, &stock.QtyBooking)
	return stock, err
}
