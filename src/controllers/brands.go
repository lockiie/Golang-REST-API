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
    "name":"Tilibra",
    "description": "<p>ótima marca de cadernos</p>",
    "active": true,
    "id": "A1"
}

*/

//CreateBrands cria uma nova marca
func CreateBrands(c *fiber.Ctx) error {
	b := models.Brands{}

	err := json.Unmarshal(c.Body(), &b)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = b.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoBrands := repositories.NewRepoBrands(tx, &ctx, nil)

	b.UsrID, b.ComID = services.ExtractUserAndCompany(c)

	if err := repoBrands.Insert(&b); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateBrands altera uma marca
func UpdateBrands(c *fiber.Ctx) error {
	b := models.Brands{}

	err := json.Unmarshal(c.Body(), &b)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = b.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoBrands := repositories.NewRepoBrands(tx, &ctx, nil)
	b.UsrID, b.ComID = services.ExtractUserAndCompany(c)

	if err := repoBrands.Update(&b, c.Params(paramID)); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteBrands Excluir uma marca
func DeleteBrands(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	_, ComID := services.ExtractUserAndCompany(c)

	repoBrands := repositories.NewRepoBrands(tx, &ctx, nil)
	if err := repoBrands.Delete(c.Params(paramID), ComID); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//QueryBrands faz consulta de brands e retorna um array de brands
func QueryBrands(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoBrands := repositories.NewRepoBrands(nil, &ctx, conn)

	_, ComID := services.ExtractUserAndCompany(c)

	var pag repositories.Pagination
	pag.Limit = c.Query(types.Limit)
	pag.OffSet = c.Query(types.OffSet)

	brands, err := repoBrands.Query(c.Query(name), c.Query(status), c.Query(paramID), ComID, pag)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(brands)
}

//QueryBrandsByID faz consulta de brands e retorna um item de brands
func QueryBrandsByID(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer conn.Close()
	_, comID := services.ExtractUserAndCompany(c)

	repoBrands := repositories.NewRepoBrands(nil, &ctx, conn)

	brand, err := repoBrands.QueryByID(c.Params(paramID), comID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(brand)
}
