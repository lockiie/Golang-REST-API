package repositories

import (
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//InsertQryImage insere uma imagem com um select
func (db RepoProducts) InsertQryImage(i models.ProductsImages, SKU string, ComID uint32) error {
	_, err := db.conn.ExecContext(
		*db.ctx,
		`INSERT INTO TB_PRODUCTS_IMAGES (PRI_ORDER, PRI_URI, PRO_ID) 
		VALUES(:0, :1, (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :2 AND COM_ID = :3))`,
		i.Order, i.URI, SKU, ComID,
	)
	return err
}

//InsertImage insere uma imagem a um produto
func (db RepoProducts) InsertImage(i *models.ProductsImages) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_PRODUCTS_IMAGES (PRI_ORDER, PRI_URI, PRO_ID) 
		VALUES (:PRI_ORDER, :PRI_URI, :PRO_ID)`,
		i.Order, i.URI, i.ProID,
	)
	return err
}

//UpdateQryImage altera uma imagem com um select
func (db RepoProducts) UpdateQryImage(i *models.ProductsImages, SKU string, comID uint32, order uint64) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`UPDATE TB_PRODUCTS_IMAGES SET PRI_ORDER = :0, PRI_URI = :1
		WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :2 AND COM_ID = :3)
		  AND PRI_ORDER = :4`,
		i.Order, i.URI, SKU, comID, order,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//DeleteQryImage exclui uma categoria de um produto
func (db RepoProducts) DeleteQryImage(SKU string, comID uint32, Order uint64) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`DELETE FROM TB_PRODUCTS_IMAGES
		WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :0 AND COM_ID = :1)
		  AND PRI_ORDER = :2`,
		SKU, comID, Order,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//QueryBySKUImages retorna um array de imagens
func (db RepoProducts) QueryBySKUImages(SKU string, comID uint32) (*[]models.ProductsImages, error) {
	var images []models.ProductsImages

	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT I.PRI_ORDER, I.PRI_URI
		 FROM PRODUCTS_IMAGES I
		 WHERE PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :PRO_SKU AND COM_ID = :COM_ID)`, SKU, comID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ima models.ProductsImages
		rows.Scan(&ima.Order, &ima.URI)
		images = append(images, ima)
	}
	return &images, err
}
