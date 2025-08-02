package inicialrepository

import (
	"context"
	"etl-service/src/config/database"
	bancoinicial "etl-service/src/config/model/banco_inicial"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

// dataInicialRepository é a implementação concreta da interface InicialRepository.
// Responsável por executar operações de leitura na base de dados MongoDB para a entidade Membro.
type dataInicialRepository struct {
	conn database.MongoConnection // Interface que gerencia a conexão com o MongoDB.
}

// NewDataInicialRepository cria e retorna uma nova instância de dataInicialRepository
// recebendo uma conexão MongoConnection para interação com o banco.
func NewDataInicialRepository(conn database.MongoConnection) InicialRepository {
	return &dataInicialRepository{
		conn: conn,
	}
}

// GetAllMembrosRequisicao busca todos os documentos da coleção membros no banco definido nas variáveis de ambiente.
//
// Essa função:
// - Obtém o contexto com timeout da conexão Mongo para evitar consultas muito longas.
// - Lê as variáveis de ambiente MONGO_DB_NAME e MONGO_COLLECTION_MEMBRO para determinar o banco e coleção.
// - Executa uma consulta Find com filtro vazio (bson.D{}), ou seja, retorna todos os documentos da coleção.
// - Retorna um slice com todos os membros encontrados ou erro, que pode ser por timeout ou falha na consulta.
//
// Caso o timeout do contexto seja excedido durante a busca ou decodificação, retorna erro específico de timeout.
func (d *dataInicialRepository) GetAllMembrosRequisicao() ([]bancoinicial.Membro, error) {
	// Obtém contexto com timeout via conexão
	ctx, cancel := d.conn.ContextWithTimeout()
	defer cancel()

	// Obtém nome do banco via variável de ambiente
	MONGO_DB_NAME := os.Getenv("MONGO_DB_NAME")
	if MONGO_DB_NAME == "" {
		log.Fatal("❌ Variável de ambiente MONGO_DB_NAME não configurada.")
	}

	// Obtém nome da coleção via variável de ambiente
	MONGO_COLLECTION_MEMBRO := os.Getenv("MONGO_COLLECTION_MEMBRO")
	if MONGO_COLLECTION_MEMBRO == "" {
		log.Fatal("❌ Variável de ambiente MONGO_COLLECTION_MEMBRO não configurada.")
	}

	// Obtém referência à coleção do MongoDB para realizar operações
	collection := d.conn.Collection(MONGO_DB_NAME, MONGO_COLLECTION_MEMBRO)

	// Executa a consulta Find para obter todos os documentos
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("tempo limite excedido para buscar membros")
		}
		return nil, fmt.Errorf("erro ao buscar membros: %w", err)
	}
	defer cursor.Close(ctx)

	var membros []bancoinicial.Membro
	// Decodifica todos os documentos retornados para o slice membros
	if err := cursor.All(ctx, &membros); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("tempo limite excedido ao decodificar os membros")
		}
		return nil, fmt.Errorf("erro ao decodificar os membros: %w", err)
	}

	// Retorna o slice preenchido com os membros encontrados
	return membros, nil
}
