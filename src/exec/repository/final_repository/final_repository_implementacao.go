package finalrepository

import (
	"etl-service/src/config/database"
	bancofinal "etl-service/src/config/model/banco_final"
	"fmt"
	"log"
	"os"
)

// dataInicialRepository é a implementação concreta da interface InicialRepository.
// Responsável por executar operações de leitura na base de dados MongoDB para a entidade Membro.
type dataFinalRepository struct {
	conn database.MongoConnection // Interface que gerencia a conexão com o MongoDB.
}

// NewDataInicialRepository cria e retorna uma nova instância de dataInicialRepository,
// recebendo uma conexão MongoConnection para interação com o banco.
func NewDataFinalRepository(conn database.MongoConnection) FinalRepository {
	return &dataFinalRepository{
		conn: conn,
	}
}

// Insert insere um novo membro na coleção do banco inicial.
//
// Parâmetros:
// - membro: objeto do tipo bancoinicial.Membro contendo os dados a serem inseridos.
//
// Fluxo da função:
// - Obtém contexto com timeout da conexão para evitar operações longas.
// - Lê as variáveis de ambiente MONGO_DB_BANCO_FINAL e MONGO_COLLECTION_BANCO_FINAL para determinar banco e coleção.
// - Insere o documento na coleção usando InsertOne.
// - Retorna erro em caso de falha na inserção ou no contexto.
//
// Uso:
// err := repo.Insert(novoMembro)
//
//	if err != nil {
//	    // Tratar erro
//	}
func (d *dataFinalRepository) Insert(membro bancofinal.Membro) error {
	ctx, cancel := d.conn.ContextWithTimeout()
	defer cancel()

	MONGO_DB_BANCO_FINAL := os.Getenv("MONGO_DB_BANCO_FINAL")
	if MONGO_DB_BANCO_FINAL == "" {
		log.Fatal("❌ Variável de ambiente MONGO_DB_BANCO_FINAL não configurada.")
	}

	MONGO_COLLECTION_BANCO_FINAL := os.Getenv("MONGO_COLLECTION_BANCO_FINAL")
	if MONGO_COLLECTION_BANCO_FINAL == "" {
		log.Fatal("❌ Variável de ambiente MONGO_COLLECTION_BANCO_FINAL não configurada.")
	}

	collection := d.conn.Collection(MONGO_DB_BANCO_FINAL, MONGO_COLLECTION_BANCO_FINAL)

	_, err := collection.InsertOne(ctx, membro)
	if err != nil {
		return fmt.Errorf("erro ao inserir membro: %w", err)
	}

	return nil
}
