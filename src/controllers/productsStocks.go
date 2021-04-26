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

//UpdateProductsStock atualiza o estoque de um produto
func UpdateProductsStock(c *fiber.Ctx) error {
	var stock models.ProductsStocks

	err := json.Unmarshal(c.Body(), &stock)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = stock.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	repoProducts := repositories.NewRepoProducts(tx, &ctx, nil)
	usrID, comID := services.ExtractUserAndCompany(c)
	sku := c.Params(types.ParamSKUValue)
	stock.UsrID = usrID

	if err := repoProducts.UpdateQryStock(&stock, sku, comID); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//QueryProductsStocks retorna o estoque do produto
func QueryProductsStocks(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
	_, comID := services.ExtractUserAndCompany(c)

	stockPro, err := repoProducts.QueryBySKUStock(c.Params(types.ParamSKUValue), comID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(stockPro)
}
