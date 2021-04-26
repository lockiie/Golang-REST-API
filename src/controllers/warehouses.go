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
    "name":"PADRAO",
    "ZipCode": "",
    "active": true,
	"id": ""
}
*/

//CreateWarehouses cria uma warehouse
func CreateWarehouses(c *fiber.Ctx) error {
	w := models.Warehouses{}

	err := json.Unmarshal(c.Body(), &w)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = w.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoWarehouses := repositories.NewRepoWarehouses(tx, &ctx, nil)

	w.UsrID, w.ComID = services.ExtractUserAndCompany(c)
	if err := repoWarehouses.Insert(&w); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateWarehouses altera uma warehouses
func UpdateWarehouses(c *fiber.Ctx) error {
	w := models.Warehouses{}

	err := json.Unmarshal(c.Body(), &w)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = w.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoWarehouses := repositories.NewRepoWarehouses(tx, &ctx, nil)
	w.UsrID, w.ComID = services.ExtractUserAndCompany(c)

	if err := repoWarehouses.Update(&w, c.Params(paramID)); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

// //DeleteWarehouses deleta um apelido de uma categoria
// func DeleteWarehouses(c *fiber.Ctx) error {
// 	var ctx = context.Background()
// 	conn, err := db.Pool.Conn(ctx)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	defer conn.Close()
// 	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

// 	_, ComID := services.ExtractUserAndCompany(c)

// 	repoWarehouses := repositories.NewRepoWarehouses(tx, &ctx, nil)
// 	if err := repoWarehouses.Delete(c.Params(paramID), ComID); err != nil {
// 		tx.Rollback()
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	tx.Commit()
// 	return c.SendStatus(fiber.StatusNoContent)
// }

//QueryWarehousesNiks faz consulta de Warehouses e retorna os nicks da empresa
func QueryWarehouses(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoWarehouses := repositories.NewRepoWarehouses(nil, &ctx, conn)

	_, ComID := services.ExtractUserAndCompany(c)

	var pag repositories.Pagination
	pag.Limit = c.Query(types.Limit)
	pag.OffSet = c.Query(types.OffSet)

	warehouses, err := repoWarehouses.Query(c.Query(title), c.Query(paramID), c.Query(status), ComID, pag)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(warehouses)
}

// //QueryWarehouses faz consulta de Warehouses e retorna um array de Warehouses
// func QueryWarehouses(c *fiber.Ctx) error {
// 	var ctx = context.Background()
// 	conn, err := db.Pool.Conn(ctx)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	defer conn.Close()
// 	repoWarehouses := repositories.NewRepoWarehouses(nil, &ctx, conn)

// 	_, ComID := services.ExtractUserAndCompany(c)

// 	var pag repositories.Pagination
// 	pag.Limit = c.Query(types.Limit)
// 	pag.OffSet = c.Query(types.OffSet)

// 	Warehouses, err := repoWarehouses.Query(c.Query(title), c.Query(paramID), ComID, pag)

// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}

// 	return c.Status(fiber.StatusOK).JSON(Warehouses)
// }
