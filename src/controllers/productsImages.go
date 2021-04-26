package controllers

import (
	"context"
	"eco/src/db"
	"eco/src/models"
	"eco/src/repositories"
	"eco/src/services"
	"eco/src/types"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

//CreateProductsImages insere uma imagem par ao produto
func CreateProductsImages(c *fiber.Ctx) error {
	var imagem models.ProductsImages

	err := json.Unmarshal(c.Body(), &imagem)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	if err = imagem.Validators(); err != nil {
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

	if err = repoProducts.InsertQryImage(imagem,
		c.Params(types.ParamSKUValue), comID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateProductsImages altera uma imagem do produto
func UpdateProductsImages(c *fiber.Ctx) error {
	var imagem models.ProductsImages

	err := json.Unmarshal(c.Body(), &imagem)

	if err = imagem.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	order, err := strconv.ParseUint(c.Params(types.ParamIDValue), 10, 64)
	if err != nil {
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

	if err = repoProducts.UpdateQryImage(&imagem,
		c.Params(types.ParamSKUValue), comID, order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteProductsImages deleta uma imagem do produto
func DeleteProductsImages(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}

	order, err := strconv.ParseUint(c.Params(types.ParamIDValue), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.DeleteQryImage(c.Params(types.ParamSKUValue), comID, order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func insertImages(p *[]models.ProductsImages, proID uint32, repo *repositories.RepoProducts) error {
	for _, image := range *p {
		image.ProID = proID
		if err := repo.InsertImage(&image); err != nil {
			return err
		}
	}
	return nil
}

//QueryProductsImages retorna as imagens do produto
func QueryProductsImages(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
	_, comID := services.ExtractUserAndCompany(c)

	imagesPro, err := repoProducts.QueryBySKUImages(c.Params(types.ParamSKUValue), comID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(imagesPro)
}
