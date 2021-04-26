package controllers

import (
	"context"
	"eco/src/db"
	"eco/src/models"
	"eco/src/repositories"
	"eco/src/services"
	"eco/src/types"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// {
//     "id": "01",
//     "value": "Preta"
// }

//CreateProductsSpecifications inserir um especificação do produto
func CreateProductsSpecifications(c *fiber.Ctx) error {
	var specification models.ProductsSpecifications

	err := json.Unmarshal(c.Body(), &specification)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = specification.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.InsertTwoQrySpecification(&specification, comID, c.Params(types.ParamSKUValue)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateProductsSpecifications inserir um especificação do produto
func UpdateProductsSpecifications(c *fiber.Ctx) error {
	var specification models.ProductsSpecifications

	err := json.Unmarshal(c.Body(), &specification)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	specification.Code = c.Params(types.ParamIDValue)
	if err = specification.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.UpdateTwoQrySpecification(&specification, comID, c.Params(types.ParamSKUValue)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteProductsSpecifications excluir um especificação do produto
func DeleteProductsSpecifications(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.DeleteTwoQrySpecification(c.Params(types.ParamSKUValue),
		c.Params(types.ParamIDValue), comID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func insertSpecifications(p *[]models.ProductsSpecifications, proID uint32, comID uint32, repo *repositories.RepoProducts) error {
	for _, specification := range *p {
		specification.ProID = proID
		if err := repo.InsertQrySpecification(&specification, comID); err != nil {
			return err
		}
	}
	return nil
}

//QueryProductsSpecifications retorna as especificações do produto
func QueryProductsSpecifications(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
	_, comID := services.ExtractUserAndCompany(c)

	specificationsPro, err := repoProducts.QueryBySKUSpecifications(c.Params(types.ParamSKUValue), comID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(specificationsPro)
}
