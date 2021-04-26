package main

import (
	"fmt"

	//iniciar o servidor e o banco de dados
	_ "eco/src/db"
	//Inicia a conexão com as rotas
	_ "eco/src/routers"
)

func main() {
	fmt.Println("Iniciando app")
}
