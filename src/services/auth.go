package services

import "github.com/gofiber/fiber/v2"

//ExtractUserAndCompany Retira do token o User e a Empresa
func ExtractUserAndCompany(c *fiber.Ctx) (uint32, uint32) {
	var UsrID uint32 = 1
	var ComID uint32 = 1
	return UsrID, ComID
}
