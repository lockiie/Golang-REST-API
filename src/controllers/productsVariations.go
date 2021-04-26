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

//CreateProductsVariations altera a variação do produto
func CreateProductsVariations(c *fiber.Ctx) error {
	var variation models.ProductsVariations

	err := json.Unmarshal(c.Body(), &variation)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = variation.Validators(); err != nil {
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

	variation.UsrID, variation.ComID = services.ExtractUserAndCompany(c)
	if err = insertVarition(&variation, c.Params(types.ParamSKUValue), repoProducts, types.Insert); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateProductsVariations altera a variação do produto
func UpdateProductsVariations(c *fiber.Ctx) error {
	var variation models.ProductsVariations
	variation.SKU = c.Params(types.ParamIDValue)

	err := json.Unmarshal(c.Body(), &variation)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	if err = variation.Validators(); err != nil {
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

	variation.UsrID, variation.ComID = services.ExtractUserAndCompany(c)
	if err = insertVarition(&variation, c.Params(types.ParamSKUValue), repoProducts, types.Update); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteProductsVariations deleta uma variação do produto
func DeleteProductsVariations(c *fiber.Ctx) error {

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.DeleteQryVariation(c.Params(types.ParamIDValue), c.Params(types.ParamSKUValue),
		comID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func insertVarition(variation *models.ProductsVariations, SKU string, repoProducts *repositories.RepoProducts, action types.ProAction) error {
	var err error
	switch action {
	case types.Insert:
		if variation.Variation == 0 {
			if err = repoProducts.InsertQryVariation(variation, SKU); err != nil {
				return err
			}
		} else {
			if err = repoProducts.InsertVariation(variation); err != nil {
				return err
			}
		}
	case types.Update:
		if err = repoProducts.UpdateQryVariation(variation, SKU); err != nil {
			return err
		}
	}
	if err = insertSpecifications(&variation.Specifications, variation.ID, variation.ComID, repoProducts); err != nil {
		return err
	}
	if err = insertImages(&variation.Images, variation.ID, repoProducts); err != nil {
		return err
	}
	if err = insertPrices(&variation.Prices, variation.ID, variation.UsrID, repoProducts); err != nil {
		return err
	}
	variation.Stock.ProID = variation.ID
	variation.Stock.UsrID = variation.UsrID
	if err := repoProducts.UpdateStock(&variation.Stock); err != nil {
		return err
	}
	return nil
}

//QueryProductsVariations retorna as variações do produto
func QueryProductsVariations(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
	_, comID := services.ExtractUserAndCompany(c)

	variationsPro, err := repoProducts.QueryBySKUVariations(c.Params(types.ParamSKUValue), comID)

	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(variationsPro)
}
