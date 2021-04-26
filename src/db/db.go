package db

import (
	"database/sql"
	"fmt"
	"log"

	//driver para conectar no banco de dados da ORACLE
	_ "github.com/godror/godror"
)

//Pool of connections
var Pool *sql.DB

const (
	stringConnector = `user="ACADEMIA"
                         password="123456"
						 connectString="127.0.0.1:1522/lucas?connect_timeout=30"
						 libDir="/oracle19/client/instantclient_19_9"
						 poolMinSessions=5
						 poolMaxSessions=15
						 poolIncrement=0
						 `
)

func init() {
	fmt.Println("Carregando Banco de dados")
	connect()
}

//função apra conectar no banco de dados
func connect() {
	db, err := sql.Open("godror", stringConnector)
	if err != nil {
		log.Fatal("Falha ao conectar no banco de dados", err)
	}
	if err = db.Ping(); err != nil {
		db.Close()
		log.Fatal("Falha na respota do banco de dados", err)
	}
	Pool = db
	fmt.Println("Banco de dados iniciado")
}
