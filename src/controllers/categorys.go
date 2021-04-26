package controllers

import (
	"context"
	"database/sql"
	"eco/src/db"
	"eco/src/models"
	"eco/src/repositories"
	"eco/src/services"
	"eco/src/types"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

/*
{
    "title":"Eletronico",
    "father_id": null,
    "id": "014"
}
*/

//CreateCategorys cria um apelido a uma categoria
func CreateCategorys(c *fiber.Ctx) error {
	ctt := models.Categorys{}

	err := json.Unmarshal(c.Body(), &ctt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = ctt.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoCategorys := repositories.NewRepoCategorys(tx, &ctx, nil)

	ctt.UsrID, ctt.ComID = services.ExtractUserAndCompany(c)
	if err := repoCategorys.Insert(&ctt); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateCategorys altera o apelido da categoria
func UpdateCategorys(c *fiber.Ctx) error {
	ctt := models.Categorys{}

	err := json.Unmarshal(c.Body(), &ctt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = ctt.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoCategorys := repositories.NewRepoCategorys(tx, &ctx, nil)
	ctt.UsrID, ctt.ComID = services.ExtractUserAndCompany(c)

	if err := repoCategorys.Update(&ctt, c.Params(paramID)); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteCategorys deleta um apelido de uma categoria
func DeleteCategorys(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	_, ComID := services.ExtractUserAndCompany(c)

	repoCategorys := repositories.NewRepoCategorys(tx, &ctx, nil)
	if err := repoCategorys.Delete(c.Params(paramID), ComID); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//QueryCategorysNiks faz consulta de Categorys e retorna os nicks da empresa
func QueryCategorysNiks(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoCategorys := repositories.NewRepoCategorys(nil, &ctx, conn)

	_, ComID := services.ExtractUserAndCompany(c)

	var pag repositories.Pagination
	pag.Limit = c.Query(types.Limit)
	pag.OffSet = c.Query(types.OffSet)

	categorys, err := repoCategorys.QueryNiks(c.Query(title), c.Query(paramID), ComID, pag)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(categorys)
}

//QueryCategorys faz consulta de Categorys e retorna um array de Categorys
func QueryCategorys(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoCategorys := repositories.NewRepoCategorys(nil, &ctx, conn)

	_, ComID := services.ExtractUserAndCompany(c)

	var pag repositories.Pagination
	pag.Limit = c.Query(types.Limit)
	pag.OffSet = c.Query(types.OffSet)

	categorys, err := repoCategorys.Query(c.Query(title), c.Query(paramID), ComID, pag)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(categorys)
}
