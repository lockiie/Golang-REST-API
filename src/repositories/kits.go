package repositories

import (
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//InsertQryKit insere um produto a um kit com um select
func (db RepoProducts) InsertQryKit(k *models.Kits, ComID uint32) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_KITS (PRO_ID, KIT_PRO_ID, KIT_QTY, KIT_UPDATE, USR_ID) 
		VALUES(:PRO_ID, (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :PRO_SKU AND COM_ID = :COM_ID), :KIT_QTD, SYSDATE, :USR_ID)`,
		k.ProID, k.KitSKU, ComID, k.Qty, k.UsrID,
	)
	return err
}

//InsertTwoQryKit insere um produto a um kit com dois selcts ||||| Não quero usar por enquanto
func (db RepoProducts) InsertTwoQryKit(k *models.Kits, SKU string, ComID uint32) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_KITS (PRO_ID, KIT_PRO_ID, KIT_QTY, KIT_UPDATE, USR_ID) 
		VALUES((SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :0 AND COM_ID = :1), 
		(SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :2 AND COM_ID = :3), :4, SYSDATE, :5)`,
		SKU, ComID, k.KitSKU, ComID, k.Qty, k.UsrID,
	)
	return err
}

//UpdateQryKit altera um kit com base no produto ||||| Não quero usar por enquanto
func (db RepoProducts) UpdateQryKit(k *models.Kits, SKU string, ComID uint32) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_KITS SET KIT_QTY = :0, USR_ID = :1, KIT_UPDATE = SYSDATE
		WHERE KIT_PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :2 AND COM_ID = :3)
		  AND PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :4 AND COM_ID = :5)`, //esse é para o produto principal
		k.Qty, k.UsrID, k.KitSKU, ComID, SKU, ComID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//DeleteQryKit desvincula um produto do kit ||||| Não quero usar por enquanto
func (db RepoProducts) DeleteQryKit(SKU string, KitSKU string, comID uint32) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`DELETE FROM TB_KITS
		WHERE KIT_PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :0 AND COM_ID = :1)
		  AND PRO_ID = (SELECT PRO_ID FROM PRODUCTS WHERE PRO_SKU = :2 AND COM_ID = :3)`, //código do produto que é o kit
		KitSKU, SKU, comID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}
