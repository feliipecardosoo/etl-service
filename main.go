package main

import (
	"context"
	"log"
	"os"

	"etl-service/src/config/database"
	"etl-service/src/config/env"
	getdata "etl-service/src/exec/get_data"
	inicialrepository "etl-service/src/exec/repository/inicial_repository"
)

// main é o ponto de entrada da aplicação.
// Ele carrega as variáveis de ambiente, conecta ao banco MongoDB,
// cria as camadas de repositório e serviço, e busca todos os membros,
// imprimindo seus nomes no log.
func main() {
	// Carrega as variáveis do arquivo .env para o ambiente
	env.LoadEnv()

	// Lê a variável de ambiente com a URI do MongoDB
	bancoInicial := os.Getenv("BANCO_INICIAL")
	if bancoInicial == "" {
		log.Fatal("❌ Variável de ambiente BANCO_INICIAL não configurada.")
	}

	// Cria a conexão com o MongoDB via interface MongoConnection
	conn := database.NewMongoConnection()
	// Tenta conectar ao banco com a URI obtida
	if err := conn.Connect(bancoInicial); err != nil {
		log.Fatalf("❌ Erro ao conectar ao MongoDB: %v", err)
	}
	// Garante o fechamento da conexão ao final da execução
	defer func() {
		if err := conn.Disconnect(context.Background()); err != nil {
			log.Printf("Erro ao desconectar: %v", err)
		}
	}()

	// Inicializa o repositório de dados, passando a conexão Mongo
	repo := inicialrepository.NewDataInicialRepository(conn)
	// Inicializa o serviço de acesso a dados, injetando o repositório
	service := getdata.NewGetDataBancoInicial(repo)

	// Chama o método GetAll para buscar todos os membros no banco
	err := service.GetAll()
	if err != nil {
		log.Fatalf("Erro ao buscar membros: %v", err)
	}
}
