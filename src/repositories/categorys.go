package repositories

import (
	"context"
	"database/sql"
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//RepoCategorys é uma estrutura para inserir um usuário ao banco de dados
type RepoCategorys struct {
	tx   *sql.Tx
	ctx  *context.Context
	conn *sql.Conn
}

//NewRepoCategorys Inicia um novo repositorio de categoria
func NewRepoCategorys(tx *sql.Tx, ctx *context.Context, conn *sql.Conn) *RepoCategorys {
	return &RepoCategorys{tx, ctx, conn}
}

//Insert é para inserir um apelido a categoria
func (db RepoCategorys) Insert(ctt *models.Categorys) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_CATEGORYS_NICKNAMES(COM_ID, CTT_ID, USR_ID, CNN_REGISTER, CNN_UPDATE, CNN_CODE)
		VALUES(:0, :1, :2, SYSDATE, SYSDATE, :3)`,
		ctt.ComID, ctt.CTT_ID, ctt.UsrID, ctt.Code,
	)
	return err
}

//Update é para atulizar uma categoria
func (db RepoCategorys) Update(ctt *models.Categorys, code string) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_CATEGORYS_NICKNAMES SET CTT_ID = :0, CNN_CODE = :1,
	 	CNN_UPDATE = SYSDATE, USR_ID = :2 WHERE COM_ID = :3 AND CNN_CODE = :4`,
		ctt.CTT_ID, ctt.Code, ctt.UsrID, ctt.ComID, code,
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

//Delete é deletar uma categoria
func (db RepoCategorys) Delete(code string, ComID uint32) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`DELETE TB_CATEGORYS_NICKNAMES WHERE CNN_CODE = :0 AND COM_ID = :1`,
		code, ComID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//QueryNiks retorna uma consulta do banco de dados dos nicks da empresa
func (db RepoCategorys) QueryNiks(title string, code string, comID uint32, pag Pagination) (*[]models.Categorys, error) {
	var args []interface{}

	args = append(args, comID)
	var where string
	if title != "" {
		where += " AND REGEXP_LIKE(T.CTT_TITLE, :CTT_TITLE, 'i')"
		args = append(args, title)
	}

	if code != "" {
		where += " AND REGEXP_LIKE(N.CTT_CODE, :CTT_CODE, 'i')"
		args = append(args, code)
	}

	strpag := pag.pag()

	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT T.CTT_ID, T.CTT_TITLE, N.CNN_CODE
		FROM CATEGORYS_TITLES T, CATEGORYS_NICKNAMES N
	   WHERE T.CTT_ID = N.CTT_ID 
	     AND COM_ID = :COM_ID`+where+strpag, args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categorys []models.Categorys

	for rows.Next() {
		var category models.Categorys
		rows.Scan(&category.CTT_ID, &category.Title, &category.Code)
		categorys = append(categorys, category)
	}
	return &categorys, nil
}

//Query retorna as categorias
func (db RepoCategorys) Query(title string, code string, ComId uint32, pag Pagination) (*[]models.Categorys, error) {
	var args []interface{}

	var where string
	args = append(args, ComId)
	if title != "" {
		where += " AND REGEXP_LIKE(T.CTT_TITLE, :CTT_TITLE, 'i')"
		args = append(args, title)
	}

	if code != "" {
		where += " AND REGEXP_LIKE(N.CTT_CODE, :CTT_CODE, 'i')"
		args = append(args, code)
	}

	strpag := pag.pag()

	rows, err := db.conn.QueryContext(
		*db.ctx,
		`SELECT T.CTT_ID, T.CTT_TITLE, NVL(T.CTT_FATHER,0), LISTAGG(N.CNN_CODE, ', ')
		 FROM CATEGORYS_TITLES T, CATEGORYS_NICKNAMES N
		 WHERE T.CTT_ID = N.CTT_ID(+)
		   AND N.COM_ID IS NULL OR N.COM_ID = :COM_ID`+where+
			` GROUP BY T.CTT_ID, T.CTT_TITLE, T.CTT_FATHER`+strpag, args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categorys []models.Categorys
	for rows.Next() {
		var category models.Categorys
		err = rows.Scan(&category.CTT_ID, &category.Title, &category.Father, &category.Code)
		categorys = append(categorys, category)
	}
	return &categorys, nil
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
