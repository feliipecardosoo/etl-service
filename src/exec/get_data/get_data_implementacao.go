package getdata

import (
	"etl-service/src/exec/domain"
	inicialrepository "etl-service/src/exec/repository/inicial_repository"
	"fmt"
	"time"
)

// getDataBancoInicial é a implementação da interface GetDataBancoInicial.
// Ela encapsula a lógica para acessar dados do banco inicial por meio do repositório.
type getDataBancoInicial struct {
	repo inicialrepository.InicialRepository
}

// NewGetDataBancoInicial cria uma nova instância de getDataBancoInicial
// recebendo uma implementação da interface InicialRepository.
// Isso permite a injeção de dependência para maior flexibilidade e testabilidade.
func NewGetDataBancoInicial(repo inicialrepository.InicialRepository) GetDataBancoInicial {
	return &getDataBancoInicial{
		repo: repo,
	}
}

// GetAll busca todos os membros na fonte de dados usando o repositório.
// Retorna uma lista de membros e um erro caso a operação falhe.
func (g *getDataBancoInicial) GetAll() error {
	start := time.Now() // captura o tempo no começo

	membros, err := g.repo.GetAllMembrosRequisicao()
	if err != nil {
		return fmt.Errorf("erro ao obter membros: %w", err)
	}

	for _, m := range membros {
		domainMembro, err := domain.NewBancoFinalMembroDomain(m) // aqui é 1 a 1
		if err != nil {
			return err
		}
		model := domainMembro.ToModel()

		// Verificar se membro já existe no BD

		fmt.Println("Membro pronto pra inserir:", model.Name)
	}

	fmt.Printf("Membros totais: %d\n", len(membros))

	duration := time.Since(start) // calcula o tempo decorrido desde o start
	fmt.Printf("Tempo de execução: %s\n", duration)

	return nil
}
