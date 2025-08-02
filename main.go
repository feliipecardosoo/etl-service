package main

import (
	"log"
	"os"

	"etl-service/src/config/database"
	"etl-service/src/config/env"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	env.LoadEnv()

	// Obtém a URI do MongoDB a partir da variável de ambiente BANCO_INICIAL
	bancoInicial := os.Getenv("BANCO_INICIAL")
	if bancoInicial == "" {
		log.Fatal("❌ Variável de ambiente BANCO_INICIAL não configurada.")
	}

	// Cria uma instância da interface MongoConnection
	var conn database.MongoConnection = database.NewMongoConnection()

	// Conecta ao banco de dados usando a URI fornecida
	if err := conn.Connect(bancoInicial); err != nil {
		log.Fatalf("❌ Erro ao conectar ao MongoDB: %v", err)
	}
}
