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

// {
//     "price": 71.99,
//     "promotional_price": 90.99,
//     "active": true,
//     "marketplace_id": 1
// }

//UpdateProductsPrice atualiza o preço de um produto
func UpdateProductsPrice(c *fiber.Ctx) error {
	var price models.ProductsPrices

	err := json.Unmarshal(c.Body(), &price)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	if err = price.Validators(); err != nil {
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
	price.UsrID = usrID

	if err := repoProducts.UpdateQryPrices(&price, sku, comID); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

func insertPrices(p *[]models.ProductsPrices, proID uint32, usrID uint32, repo *repositories.RepoProducts) error {
	for _, price := range *p {
		price.ProID = proID
		price.UsrID = usrID
		if err := repo.UpdatePrices(&price); err != nil {
			return err
		}
	}
	return nil
}

//QueryProductsPrices retorna os preços do produto
func QueryProductsPrices(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
	_, comID := services.ExtractUserAndCompany(c)

	pricesPro, err := repoProducts.QueryBySKUPrices(c.Params(types.ParamSKUValue), comID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(pricesPro)
}
