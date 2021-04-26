package repositories

import (
	"database/sql"
	"eco/src/functions"
	"eco/src/models"
	"eco/src/types"
	"errors"
	"fmt"
)

//InsertVariation é para inserir um novo produto variado ao banco de dados
func (db RepoProducts) InsertVariation(p *models.ProductsVariations) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_PRODUCTS (PRO_SKU, PRO_TITLE, PRO_STATUS, COM_ID, PRO_DESCRIPTION, 
			PRO_BARCODE, PRO_VARIATION, PRO_UPDATE, PRO_CROSSDOCKING, USR_ID, PRO_REGISTER, PDT_ID) 
		 VALUES (:0, :1, :2, :3, :5,:6, :7, SYSDATE, :8, :9, SYSDATE, 2) RETURNING PRO_ID INTO :10`,
		p.SKU, p.Title, functions.BoolToByte(p.Status), p.ComID, p.Description,
		p.Barcode, p.Variation, p.Crossdocking, p.UsrID, sql.Out{Dest: &p.ID},
	)

	return err
}

//InsertQryVariation é para inserir um novo produto variado ao banco de dados Ccom um select
func (db RepoProducts) InsertQryVariation(p *models.ProductsVariations, SKU string) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_PRODUCTS (PRO_SKU, PRO_TITLE, PRO_STATUS, COM_ID, PRO_DESCRIPTION, 
			PRO_BARCODE, PRO_VARIATION, PRO_UPDATE, PRO_CROSSDOCKING, USR_ID, PRO_REGISTER, PDT_ID) 
		 VALUES (:0, :1, :2, :3, :4, :5, (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :6 AND COM_ID = :7 AND PDT_ID = 2), 
		 SYSDATE, :9, :10, SYSDATE, 2) RETURNING PRO_ID INTO :11`,
		p.SKU, p.Title, functions.BoolToByte(p.Status), p.ComID, p.Description,
		p.Barcode, SKU, p.ComID, p.Crossdocking, p.UsrID, sql.Out{Dest: &p.ID},
	)
	return err
}

//UpdateQryVariation é para alterar um produto variado com um select
func (db RepoProducts) UpdateQryVariation(p *models.ProductsVariations, SKU string) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_PRODUCTS SET PRO_TITLE = :0, PRO_STATUS = :1, PRO_DESCRIPTION = :2, PRO_BARCODE = :3,
		PRO_UPDATE = SYSDATE, PRO_REGISTER = SYSDATE, PRO_CROSSDOCKING = :4, USR_ID = :5
		WHERE PRO_SKU = :6 AND COM_ID = :7 AND 
		PRO_VARIATION = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU =:8 AND COM_ID = :9 AND PDT_ID = 2) RETURNING PRO_ID INTO :10`,
		p.Title, functions.BoolToByte(p.Status), p.Description,
		p.Barcode, p.Crossdocking, p.UsrID, p.SKU, p.ComID, SKU, p.ComID, sql.Out{Dest: &p.ID},
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//DeleteVariations deleta produtos com base no id
func (db RepoProducts) DeleteVariations(ProID uint32) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`DELETE TB_PRODUCTS WHERE PRO_VARIATION = :PRO_ID AND PRO_VARIATION <> PRO_ID AND PDT_ID = 2`,
		ProID,
	)

	return err
}

//DeleteQryVariation deleta uma variação do produto com um select
func (db RepoProducts) DeleteQryVariation(SKU string, variation string, comID uint32) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`DELETE FROM TB_PRODUCTS WHERE PRO_SKU = :0 AND COM_ID = :1 AND 
		PRO_VARIATION = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :3 AND COM_ID = :4)`,
		SKU, comID, variation, comID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//QueryBySKUVariations retorna um array da variação dos produtos
func (db RepoProducts) QueryBySKUVariations(SKU string, comID uint32) (*[]models.ProductsVariations, error) {

	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT P.PRO_SKU, P.PRO_TITLE, P.PRO_STATUS, P.PRO_DESCRIPTION, 
		P.PRO_BARCODE, P.PRO_CROSSDOCKING, S.PRO_QTY, S.PRO_QTYBOOKING,
	    CURSOR(SELECT SP.SPE_VALUE, ST.SPT_CODE, ST.SPT_TITLE
			 FROM SPECIFICATIONS_PRODUCTS SP, SPECIFICATIONS_TITLES ST
			 WHERE SP.SPT_ID = ST.SPT_ID
			   AND SP.PRO_ID = P.PRO_ID),
	    CURSOR(SELECT PC.PRP_PRICE, PC.PRP_PRICEPROMO, PC.PRP_STATUS, PC.MKC_ID
			 FROM PRODUCTS_PRICES PC
			 WHERE PC.PRO_ID = P.PRO_ID),
	    CURSOR(SELECT I.PRI_ORDER, I.PRI_URI
			 FROM PRODUCTS_IMAGES I
			 WHERE I.PRO_ID = P.PRO_ID) IMAGES         
        FROM PRODUCTS P, PRODUCTS_STOCKS S
        WHERE P.PRO_ID = S.PRO_ID
          AND P.PRO_ID <> P.PRO_VARIATION
          AND P.PDT_ID = 2
          AND P.PRO_VARIATION = (SELECT PRO_VARIATION FROM PRODUCTS WHERE PRO_SKU = :PRO_SKU AND COM_ID = :COM_ID)`, SKU, comID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	var products []models.ProductsVariations

	for rows.Next() {
		var pro models.ProductsVariations
		var spe *sql.Rows
		var pri *sql.Rows
		var ima *sql.Rows
		err = rows.Scan(&pro.SKU, &pro.Title, &pro.Status, &pro.Description, &pro.Barcode, &pro.Crossdocking, &pro.Stock.Qty,
			&pro.Stock.QtyBooking, &spe, &pri, &ima)

		for spe.Next() {
			var specification models.ProductsSpecifications
			spe.Scan(&specification.Value, &specification.Code, &specification.Title)
			pro.Specifications = append(pro.Specifications, specification)
		}
		spe.Close()
		for pri.Next() {
			var price models.ProductsPrices
			pri.Scan(&price.Price, &price.PricePromo, &price.Status, &price.MkcID)
			pro.Prices = append(pro.Prices, price)
		}
		pri.Close()

		for ima.Next() {
			var image models.ProductsImages
			ima.Scan(&image.Order, &image.URI)
			pro.Images = append(pro.Images, image)
		}
		ima.Close()
		products = append(products, pro)
	}
	return &products, nil
}
