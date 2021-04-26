package repositories

import (
	"database/sql"
	"eco/src/functions"
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//UpdateQryPrices insere o preço do produto com um select
func (db RepoProducts) UpdateQryPrices(p *models.ProductsPrices, SKU string, ComID uint32) error {
	var err error
	var res sql.Result
	if p.MkcID > 0 { //Quer dizer que não é do marketplace
		res, err = db.tx.ExecContext(
			*db.ctx,
			`UPDATE TB_PRODUCTS_PRICES SET PRP_PRICE = :0, PRP_PRICEPROMO = :1, PRP_UPDATE = SYSDATE, PRP_STATUS = :2, USR_ID = :3 
			 WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :4 AND COM_ID = :5) AND MKC_ID = :6`,
			p.Price, p.PricePromo, functions.BoolToByte(p.Status), p.UsrID, SKU, ComID, p.MkcID,
		)
	} else {
		res, err = db.tx.ExecContext(
			*db.ctx,
			`UPDATE TB_PRODUCTS_PRICES SET PRP_PRICE = :0, PRP_PRICEPROMO = :1, PRP_UPDATE = SYSDATE, PRP_STATUS = 1, USR_ID = :3 
			WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :4 AND COM_ID = :5)`,
			p.Price, p.PricePromo, p.UsrID, SKU, ComID,
		)
	}
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//UpdatePrices altera o preço do produto
func (db RepoProducts) UpdatePrices(p *models.ProductsPrices) error {
	var res sql.Result
	var err error
	if p.MkcID > 0 { //Quer dizer que não é do marketplace
		res, err = db.tx.ExecContext(
			*db.ctx,
			`UPDATE TB_PRODUCTS_PRICES SET PRP_PRICE = :0, PRP_PRICEPROMO = :1, PRP_UPDATE = SYSDATE, PRP_STATUS = :2, 
			USR_ID = :3 WHERE PRO_ID = :4 AND MKC_ID = :6`,
			p.Price, p.PricePromo, functions.BoolToByte(p.Status), p.UsrID, p.ProID, p.MkcID,
		)
	} else {
		res, err = db.tx.ExecContext(
			*db.ctx,
			`UPDATE TB_PRODUCTS SET PRO_PRICE = :0, PRO_PRICEPROMO = :1
		     WHERE PRO_ID = :2`,
			p.Price, p.PricePromo, p.ProID,
		)
	}
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//QueryBySKUPrices retorna um array de preços do produto
func (db RepoProducts) QueryBySKUPrices(SKU string, comID uint32) ([]models.ProductsPrices, error) {
	var prices []models.ProductsPrices

	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT PC.PRP_PRICE, PC.PRP_PRICEPROMO, PC.PRP_STATUS, PC.MKC_ID
		 FROM PRODUCTS_PRICES PC
		WHERE PC.PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :PRO_SKU AND COM_ID = :COM_ID)`, SKU, comID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pri models.ProductsPrices
		rows.Scan(&pri.Price, &pri.PricePromo, &pri.Status, &pri.MkcID)
		prices = append(prices, pri)
	}
	return prices, err
}
