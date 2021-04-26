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
    "title":"Cor",
    "active": true,
    "id": "01"
}

*/

//CreateSpecifications cria uma nova marca
func CreateSpecifications(c *fiber.Ctx) error {
	spt := models.Specifications{}
	err := json.Unmarshal(c.Body(), &spt)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))

	}

	if err = spt.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoSpecifications := repositories.NewRepoSpecifications(tx, &ctx, nil)

	spt.UsrID, spt.ComID = services.ExtractUserAndCompany(c)
	if err := repoSpecifications.Insert(&spt); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateSpecifications altera uma marca
func UpdateSpecifications(c *fiber.Ctx) error {
	b := models.Specifications{}
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

	repoSpecifications := repositories.NewRepoSpecifications(tx, &ctx, nil)
	b.UsrID, b.ComID = services.ExtractUserAndCompany(c)

	if err := repoSpecifications.Update(&b, c.Params(paramID)); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteSpecifications Excluir uma marca
func DeleteSpecifications(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	_, ComID := services.ExtractUserAndCompany(c)

	repoSpecifications := repositories.NewRepoSpecifications(tx, &ctx, nil)
	if err := repoSpecifications.Delete(c.Params(paramID), ComID); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//QuerySpecifications faz consulta de Specifications e retorna um array de Specifications
func QuerySpecifications(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoSpecifications := repositories.NewRepoSpecifications(nil, &ctx, conn)

	_, ComID := services.ExtractUserAndCompany(c)

	var pag repositories.Pagination
	pag.Limit = c.Query(types.Limit)
	pag.OffSet = c.Query(types.OffSet)

	specifications, err := repoSpecifications.Query(c.Query("title"), c.Query(status), c.Query(paramID), ComID, pag)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	// if len(*Specifications) == 0 {
	// 	return c.Status(fiber.StatusOK).JSON("[]")
	// }
	return c.Status(fiber.StatusOK).JSON(specifications)
}

//QuerySpecificationsByID faz consulta de Specifications e retorna um item de Specifications
func QuerySpecificationsByID(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoSpecifications := repositories.NewRepoSpecifications(nil, &ctx, conn)

	brand, err := repoSpecifications.QueryByID(c.Params(paramID), 1)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(brand)
}
