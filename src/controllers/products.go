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

//CreateProducts cria um produto
func CreateProducts(c *fiber.Ctx) error {
	return CreateOrUpdatesProducts(c, types.Insert)
}

//UpdateProducts altera um produto
func UpdateProducts(c *fiber.Ctx) error {
	return CreateOrUpdatesProducts(c, types.Update)
}

//CreateOrUpdatesProducts cria ou altera um produto
func CreateOrUpdatesProducts(c *fiber.Ctx, action types.ProAction) error {
	pro := models.Products{}

	//validacoes
	err := json.Unmarshal(c.Body(), &pro)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	if err = pro.Validators(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	//seta a empresa e o usuario
	pro.UsrID, pro.ComID = services.ExtractUserAndCompany(c)

	isVariation := len(pro.Variations)
	isKit := len(pro.Kits)
	if isKit > 0 && isVariation > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	//fim das validacoes

	//ABRIR CONEXAO
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	repoProducts := repositories.NewRepoProducts(tx, &ctx, nil)
	//FIM DAS CONFIGURAÇÕES DE CONEXAO

	if isVariation > 0 {
		pro.PdtID = uint8(types.Variation)
	} else if isKit > 0 {
		pro.PdtID = uint8(types.Kit)
	} else {
		pro.PdtID = uint8(types.Normal)
	}
	switch action {
	case types.Insert:
		if pro.SKU == "" {
			return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
		}
		if err = repoProducts.Insert(&pro); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
		}
	case types.Update:
		pro.SKU = c.Params(types.ParamSKUValue)
		if err = repoProducts.Update(&pro); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
		}
		if err = repoProducts.DeleteVariations(pro.ID); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
		}

	}
	//tem variação
	if isVariation > 0 {
		for _, variation := range pro.Variations {
			variation.Variation = pro.ID
			variation.ComID = pro.ComID
			variation.UsrID = pro.UsrID
			if err = insertVarition(&variation, c.Params(types.ParamSKUValue), repoProducts, types.Insert); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
			}
		}
	}

	pro.Stock.UsrID = pro.UsrID
	pro.Stock.ProID = pro.ID
	if err := repoProducts.UpdateStock(&pro.Stock); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	// if err = insertCategorys(pro.Category, pro.ID, pro.ComID, repoProducts); err != nil {
	// 	tx.Rollback()
	// 	return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	// }
	if err = insertSpecifications(&pro.Specifications, pro.ID, pro.ComID, repoProducts); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	if err = insertImages(&pro.Images, pro.ID, repoProducts); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	if err = insertPrices(&pro.Prices, pro.ID, pro.UsrID, repoProducts); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	} //insertKits

	if err = insertKits(&pro.Kits, pro.ID, pro.ComID, pro.UsrID, repoProducts); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	tx.Commit()
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateProductsActivate ativa um produto
func UpdateProductsActivate(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
	usrID, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.ProUpdateActivate(c.Params(types.ParamSKUValue), comID, usrID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//UpdateProductsDeactivate desativa um produto
func UpdateProductsDeactivate(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()
	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)
	usrID, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.ProUpdateDeactivate(c.Params(types.ParamSKUValue), comID, usrID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//DeleteProducts deleta um produto
func DeleteProducts(c *fiber.Ctx) error {
	var ctx = context.Background()
	conn, err := db.Pool.Conn(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	defer conn.Close()

	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

	_, comID := services.ExtractUserAndCompany(c)

	if err = repoProducts.Delete(c.Params(types.ParamSKUValue), comID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//QueryProducts busca um produto pelo ID
// func QueryProducts(c *fiber.Ctx) error {
// 	var ctx = context.Background()
// 	conn, err := db.Pool.Conn(ctx)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	defer conn.Close()

// 	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

// 	_, comID := services.ExtractUserAndCompany(c)

// 	var pag repositories.Pagination
// 	pag.Limit = c.Query(types.Limit)
// 	pag.OffSet = c.Query(types.OffSet)

// 	products, err := repoProducts.Query(c.Query(title), c.Query(status), c.Query(types.ParamSKUValue),
// 		c.Query(proType), comID, pag)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	return c.Status(fiber.StatusOK).JSON(products)
// }

// //QueryProductsByID retorna um produto pelo ID
// func QueryProductsByID(c *fiber.Ctx) error {
// 	var ctx = context.Background()
// 	conn, err := db.Pool.Conn(ctx)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	defer conn.Close()

// 	repoProducts := repositories.NewRepoProducts(nil, &ctx, conn)

// 	_, comID := services.ExtractUserAndCompany(c)

// 	product, err := repoProducts.QueryByID(c.Params(types.ParamSKUValue), comID)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error()))
// 	}
// 	return c.Status(fiber.StatusOK).JSON(product)
// }
