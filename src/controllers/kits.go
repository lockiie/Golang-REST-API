package controllers

import (
	"eco/src/models"
	"eco/src/repositories"
)

// //CreateKits liga um novo produto ao kit
// func CreateKits(c *fiber.Ctx) error {
// 	kit := models.Kits{}

// 	err := json.Unmarshal(c.Body(), &kit)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error(), types.CodeErrJSONDecode))
// 	}

// 	if err = kit.Validators(); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error(), types.CodeErrValidate))
// 	}

// 	var ctx = context.Background()
// 	conn, err := db.Pool.Conn(ctx)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error(), types.CodeErrPoolDB))
// 	}
// 	defer conn.Close()

// 	tx, _ := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

// 	repo := repositories.NewRepoProducts(tx, &ctx, nil)
// 	var comID uint32
// 	kit.UsrID, comID = services.ExtractUserAndCompany(c)
// 	if err := repo.InsertTwoQryKit(&kit, c.Params(types.ParamSKUValue), comID); err != nil {
// 		tx.Rollback()
// 		return c.Status(fiber.StatusBadRequest).JSON(models.SendError(err.Error(), types.CodeErrInsertDB))
// 	}
// 	tx.Commit()
// 	return c.SendStatus(fiber.StatusNoContent)
// }

func insertKits(k *[]models.Kits, proID uint32, comID uint32, usrID uint32, repo *repositories.RepoProducts) error {
	for _, kit := range *k {
		kit.ProID = proID
		kit.UsrID = usrID
		if err := repo.InsertQryKit(&kit, comID); err != nil {
			return err
		}
	}
	return nil
}
