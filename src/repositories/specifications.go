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

//RepoSpecifications é uma estrutura para inserir um usuário ao banco de dados
type RepoSpecifications struct {
	tx   *sql.Tx
	ctx  *context.Context
	conn *sql.Conn
}

//NewRepoSpecifications Inicia um novo repositorio de especificação
func NewRepoSpecifications(tx *sql.Tx, ctx *context.Context, conn *sql.Conn) *RepoSpecifications {
	return &RepoSpecifications{tx, ctx, conn}
}

//Insert é para inserir uma nova especificação
func (db RepoSpecifications) Insert(spt *models.Specifications) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_SPECIFICATIONS_TITLES(SPT_TITLE, SPT_CODE, SPT_STATUS, USR_ID, COM_ID, SPT_REGISTER, SPT_UPDATE)
		VALUES(:SPT_TITLE, :SPT_CODE, :SPT_STATUS, :USR_ID, :COM_ID, SYSDATE, SYSDATE)
		RETURNING SPT_ID INTO :SPT_ID`,
		spt.Title, spt.Code, functions.BoolToByte(spt.Status), spt.UsrID, spt.ComID,
		sql.Out{Dest: &spt.ID},
	)
	return err
}

//Update é para atulizar uma especificação
func (db RepoSpecifications) Update(spt *models.Specifications, code string) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_SPECIFICATIONS_TITLES SET SPT_TITLE = :SPT_TITLE, SPT_CODE = :SPT_CODE, SPT_STATUS = :SPT_STATUS,
		SPT_UPDATE = SYSDATE, USR_ID = :USR_ID WHERE COM_ID = :COM_ID AND SPT_CODE = :CODE`,
		spt.Title, spt.Code, functions.BoolToByte(spt.Status), spt.UsrID, spt.ComID, code,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//Delete é deletar uma especificação
func (db RepoSpecifications) Delete(code string, ComID uint32) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`DELETE TB_SPECIFICATIONS_TITLES WHERE SPT_CODE = :0 AND COM_ID = :1`,
		code, ComID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//Query retorna uma consulta do banco de dados
func (db RepoSpecifications) Query(title string, status string, code string, comID uint32, pag Pagination) (*[]models.Specifications, error) {
	var args []interface{}

	args = append(args, comID)
	var where string
	if title != "" {
		where += " AND REGEXP_LIKE(SPT_TITLE, :SPT_TITLE, 'i')"
		args = append(args, title)
	}

	if code != "" {
		where += " AND REGEXP_LIKE(SPT_CODE, :SPT_CODE, 'i')"
		args = append(args, code)
	}
	if bStatus := functions.BoolStrToByte(status); bStatus != 10 {
		where += " AND SPT_STATUS = " + strconv.Itoa(int(bStatus))
	}

	strpag := pag.pag()

	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT SPT_TITLE, SPT_CODE, SPT_STATUS
		 FROM SPECIFICATIONS_TITLES WHERE COM_ID = :COM_ID`+where+strpag, args...,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var Specifications []models.Specifications

	for rows.Next() {
		var category models.Specifications
		rows.Scan(&category.Title, &category.Code, &category.Status)
		Specifications = append(Specifications, category)
	}
	return &Specifications, nil
}

//QueryByID retorna uma consulta do banco de dados que traz 1 registro no máximo
func (db RepoSpecifications) QueryByID(code string, comID uint32) (models.Specifications, error) {
	row := db.conn.QueryRowContext(
		*db.ctx,
		`SELECT SPT_TITLE, SPT_CODE, SPT_STATUS
		FROM SPECIFICATIONS_TITLES WHERE COM_ID = :COM_ID AND SPT_CODE = :SPT_CODE`, comID, code,
	)
	var speicifications models.Specifications
	err := row.Scan(&speicifications.Title, &speicifications.Code, &speicifications.Status)
	if err != nil {
		return speicifications, errors.New(types.MsgNotFound)
	}
	return speicifications, err
}
