package models

//Errors é uma estrutura personalizada para ser retornada na api
type Errors struct {
	// Code    uint16 `json:"code"`
	Message string `json:"message"`
}

//SendError cria uma mensagem de erro
func SendError(message string) Errors {
	return Errors{message}
}

// func (err *Errors) SentError(c *fiber.Ctx) error {
// 	return
// }
