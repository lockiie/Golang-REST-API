package controllers

import (
	"context"
	"eco/src/db"
	"eco/src/models"
	"eco/src/repositories"
	"eco/src/services"
	"eco/src/types"

	"github.com/gofiber/fiber/v2"
)

//CreateProductsCategorys insere um categoria no produto
func CreateProductsCategorys(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.InsertTwoQryCategory(c.Params(types.ParamSKUValue),
		c.Params(types.ParamIDValue), comID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteProductsCategorys deleta uma categoria do produto
func DeleteProductsCategorys(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.DeleteTwoQryCategory(c.Params(types.ParamSKUValue),
		c.Params(types.ParamIDValue), comID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// func insertCategorys(p string, proID uint32, comID uint32, repo *repositories.RepoProducts) error {
// 	for _, code := range p {
// 		var category models.ProductsCategorys
// 		category.ProID = proID
// 		category.Code = code
// 		category.ProID = proID
// 		if err := repo.InsertQryCategory(&category, comID); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// //QueryCategorysProducts retorna as categorias do produto
// func QueryProductsCategorys(c *fiber.Ctx) error {
// 	var ctx = context.Background()
// 	conn, err := db.Pool.Conn(ctx)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	defer conn.Close()
// 	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
// 	_, comID := services.ExtractUserAndCompany(c)

// 	categorysPro, err := repoProducts.QueryBySKUCategorys(c.Params(types.ParamSKUValue), comID)

// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	return c.Status(fiber.StatusOK).JSON(categorysPro)
// }
