package inicialrepository

import (
	"context"
	"etl-service/src/config/database"
	bancofinal "etl-service/src/config/model/banco_final"
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

// NewDataInicialRepository cria e retorna uma nova instância de dataInicialRepository,
// recebendo uma conexão MongoConnection para interação com o banco.
func NewDataInicialRepository(conn database.MongoConnection) InicialRepository {
	return &dataInicialRepository{
		conn: conn,
	}
}

// GetAllMembrosRequisicao busca todos os documentos da coleção membros no banco definido nas variáveis de ambiente.
//
// Fluxo da função:
// - Obtém contexto com timeout via conexão para limitar tempo de consulta.
// - Lê as variáveis de ambiente MONGO_DB_NAME e MONGO_COLLECTION_MEMBRO para definir banco e coleção.
// - Realiza consulta Find com filtro vazio (bson.D{}), retornando todos os documentos da coleção.
// - Decodifica os documentos retornados para um slice de bancoinicial.Membro.
// - Retorna o slice de membros ou erro caso a consulta ou decodificação falhe.
//
// Tratamento especial para erros de timeout do contexto, retornando mensagens específicas.
func (d *dataInicialRepository) GetAllMembrosRequisicao() ([]bancoinicial.Membro, error) {
	ctx, cancel := d.conn.ContextWithTimeout()
	defer cancel()

	MONGO_DB_NAME := os.Getenv("MONGO_DB_NAME")
	if MONGO_DB_NAME == "" {
		log.Fatal("❌ Variável de ambiente MONGO_DB_NAME não configurada.")
	}

	MONGO_COLLECTION_MEMBRO := os.Getenv("MONGO_COLLECTION_MEMBRO")
	if MONGO_COLLECTION_MEMBRO == "" {
		log.Fatal("❌ Variável de ambiente MONGO_COLLECTION_MEMBRO não configurada.")
	}

	collection := d.conn.Collection(MONGO_DB_NAME, MONGO_COLLECTION_MEMBRO)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("tempo limite excedido para buscar membros")
		}
		return nil, fmt.Errorf("erro ao buscar membros: %w", err)
	}
	defer cursor.Close(ctx)

	var membros []bancoinicial.Membro
	if err := cursor.All(ctx, &membros); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("tempo limite excedido ao decodificar os membros")
		}
		return nil, fmt.Errorf("erro ao decodificar os membros: %w", err)
	}

	return membros, nil
}

// ExistsByNames verifica a existência de múltiplos nomes na coleção do banco final.
//
// Parâmetros:
// - names: slice de strings com os nomes a serem verificados.
//
// Fluxo da função:
// - Cria contexto com timeout.
// - Obtém nomes do banco e da coleção via variáveis de ambiente.
// - Executa uma consulta usando filtro {$in: names} para encontrar membros com esses nomes.
// - Itera sobre os resultados e preenche um mapa string->bool indicando quais nomes existem.
//
// Retorna:
// - Um mapa em que a chave é o nome e o valor bool indica se existe no banco.
// - Erro em caso de falha na consulta ou decodificação.
func (d *dataInicialRepository) ExistsByNames(names []string) (map[string]bool, error) {
	ctx, cancel := d.conn.ContextWithTimeout()
	defer cancel()

	dbName := os.Getenv("MONGO_DB_BANCO_FINAL")
	collectionName := os.Getenv("MONGO_COLLECTION_BANCO_FINAL")

	collection := d.conn.Collection(dbName, collectionName)

	filter := bson.M{"name": bson.M{"$in": names}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	existing := make(map[string]bool)
	for cursor.Next(ctx) {
		var m struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		existing[m.Name] = true
	}

	return existing, nil
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
func (d *dataInicialRepository) Insert(membro bancofinal.Membro) error {
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
