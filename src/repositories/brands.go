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

//RepoBrands é uma estrutura para inserir um usuário ao banco de dados
type RepoBrands struct {
	tx   *sql.Tx
	ctx  *context.Context
	conn *sql.Conn
}

//NewRepoBrands Inicia um novo repositorio de marcas
func NewRepoBrands(tx *sql.Tx, ctx *context.Context, conn *sql.Conn) *RepoBrands {
	return &RepoBrands{tx, ctx, conn}
}

//Insert é para inserir um novo usuário ao banco de dados
func (db RepoBrands) Insert(b *models.Brands) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_BRANDS(BND_NAME,BND_DESCRIPTION, BND_STATUS, BND_REGISTER, BND_UPDATE, BND_CODE, USR_ID, COM_ID)
		VALUES(:BND_NAME, :BND_DESCRIPTION, :BND_STATUS, SYSDATE, SYSDATE, :BND_CODE, :USR_ID, :COM_ID) 
		RETURNING BND_ID INTO :BND_ID`,
		b.Name, b.Description, functions.BoolToByte(b.Status), b.Code, b.UsrID, b.ComID,
		sql.Out{Dest: &b.ID},
	)
	return err
}

//Update é para atulizar uma marca
func (db RepoBrands) Update(b *models.Brands, code string) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_BRANDS SET BND_NAME =:0, BND_DESCRIPTION = :1, BND_STATUS = :2, 
		                BND_UPDATE = SYSDATE, BND_CODE =:3, USR_ID =:4 WHERE BND_CODE = :5 AND COM_ID = :6`,
		b.Name, b.Description, functions.BoolToByte(b.Status), b.Code, b.UsrID, code, b.ComID,
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

//Delete é deletar uma marca
func (db RepoBrands) Delete(code string, ComID uint32) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`DELETE TB_BRANDS WHERE BND_CODE = :0 AND COM_ID = :1`,
		code, ComID,
	)
	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return nil
}

//Query retorna uma consulta do banco de dados
func (db RepoBrands) Query(name string, status string, code string, comID uint32, pag Pagination) (*[]models.Brands, error) {
	var args []interface{}

	args = append(args, comID)
	var where string
	if name != "" {
		where += " AND REGEXP_LIKE(BND_NAME, :BND_NAME, 'i')"
		args = append(args, name)
	}

	if bStatus := functions.BoolStrToByte(status); bStatus != 10 {
		where += " AND BND_STATUS = " + strconv.Itoa(int(bStatus))
	}

	if code != "" {
		where += " AND REGEXP_LIKE(BND_CODE, :BND_CODE, 'i')"
		args = append(args, code)
	}
	strpag := pag.pag()
	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT BND_NAME, BND_DESCRIPTION, BND_CODE, BND_STATUS 
		 FROM BRANDS WHERE COM_ID = :COM_ID`+where+strpag, args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var brands []models.Brands

	for rows.Next() {
		var brand models.Brands
		rows.Scan(&brand.Name, &brand.Description, &brand.Code, &brand.Status)
		brands = append(brands, brand)
	}
	return &brands, nil
}

//QueryByID retorna uma consulta do banco de dados que trans 1 registro no máximo
func (db RepoBrands) QueryByID(code string, comID uint32) (models.Brands, error) {

	row := db.conn.QueryRowContext(
		*db.ctx,
		`SELECT BND_NAME, BND_DESCRIPTION, BND_CODE, BND_STATUS 
		 FROM BRANDS WHERE COM_ID = :COM_ID AND BND_CODE = :BND_CODE`, comID, code,
	)
	var brand models.Brands
	err := row.Scan(&brand.Name, &brand.Description, &brand.Code, &brand.Status)
	if err != nil {
		return brand, errors.New(types.MsgNotFound)
	}
	return brand, err
}
