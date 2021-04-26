package repositories

import (
	"context"
	"database/sql"
	"eco/src/functions"
	"eco/src/models"
	"eco/src/types"
	"errors"
	"strconv"
)

//RepoWarehouses é uma estrutura para inserir um warehouses ao banco de dados
type RepoWarehouses struct {
	tx   *sql.Tx
	ctx  *context.Context
	conn *sql.Conn
}

//NewRepoWarehouses Inicia um novo repositorio de categoria
func NewRepoWarehouses(tx *sql.Tx, ctx *context.Context, conn *sql.Conn) *RepoWarehouses {
	return &RepoWarehouses{tx, ctx, conn}
}

//Insert é para inserir um apelido a categoria
func (db RepoWarehouses) Insert(w *models.Warehouses) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_WAREHOUSES (WHS_NAME, WHS_STATUS, WHS_ZIPCODE, WHS_CODE, COM_ID, USR_ID, WHS_REGISTER, WHS_UPDATE) 
		 VALUES (:0, :1, :2, :3, :4, :5, SYSDATE, SYSDATE)`,
		w.Name, w.Status, w.ZipCode, w.Code, w.ComID, w.UsrID,
	)
	return err
}

//Update é para atulizar uma categoria
func (db RepoWarehouses) Update(w *models.Warehouses, code string) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_WAREHOUSES SET WHS_NAME = :0, WHS_STATUS = :1,
		WHS_UPDATE = SYSDATE, WHS_ZIPCODE = :2, WHS_CODE = :3, USR_ID = :4 WHERE COM_ID = :5 AND WHS_CODE = :6`,
		w.Name, w.Status, w.ZipCode, w.Code, w.UsrID, w.ComID, code,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}

	return nil
}

//Query retorna as WAREHOUSES
func (db RepoWarehouses) Query(title string, code string, status string, ComId uint32, pag Pagination) (*[]models.Warehouses, error) {
	var args []interface{}

	var where string
	args = append(args, ComId)
	if title != "" {
		where += " AND REGEXP_LIKE(W.WHS_NAME, :WHS_NAME, 'i')"
		args = append(args, title)
	}

	if code != "" {
		where += " AND REGEXP_LIKE(W.WHS_CODE, :WHS_CODE, 'i')"
		args = append(args, code)
	}

	if bStatus := functions.BoolStrToByte(status); bStatus != 10 {
		where += " AND W.WHS_STATUS = " + strconv.Itoa(int(bStatus))
	}

	strpag := pag.pag()

	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT W.WHS_NAME, W.WHS_STATUS, W.WHS_ZIPCODE, W.WHS_CODE
		 FROM WAREHOUSES W
		 WHERE W.COM_ID = :0 `+where+strpag, args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var warehouses []models.Warehouses
	for rows.Next() {
		var warehouse models.Warehouses
		err = rows.Scan(&warehouse.Name, &warehouse.Status, &warehouse.ZipCode, &warehouse.Code)
		warehouses = append(warehouses, warehouse)
	}
	return &warehouses, nil
}

// //QueryByID retorna uma consulta do banco de dados que trans 1 registro no máximo
// func (db RepoCategorys) QueryByID(code string, comID uint32) (models.Categorys, error) {
// 	row := db.conn.QueryRowContext(
// 		*db.ctx,
// 		`SELECT CTT_ID, CTT_TITLE, CTT_CODE, CTT_FATHER
// 		FROM CATEGORYS_TITLES WHERE COM_ID = :COM_ID AND CTT_CODE = :CTT_CODE`, comID, code,
// 	)
// 	var category models.Categorys
// 	err := row.Scan(&category.ID, &category.Title, &category.Code, &category.Father)
// 	if err != nil {
// 		return category, errors.New(types.MsgNotFound)
// 	}
// 	return category, err
// }
