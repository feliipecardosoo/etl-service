package getdata

import (
	"encoding/json"
	bancoinicial "etl-service/src/config/model/banco_inicial"
	inicialrepository "etl-service/src/exec/repository/inicial_repository"
	"fmt"
	"log"
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
func (g *getDataBancoInicial) GetAll() ([]bancoinicial.Membro, error) {
	membros, err := g.repo.GetAllMembrosRequisicao()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter membros: %w", err)
	}

	// Verificar como esta vindo essa variavel membros
	// antes do insert, verificar se ja tem o membro com este nome.

	// Imprime cada membro como JSON formatado
	for i, m := range membros {
		jsonData, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			log.Printf("Erro ao converter membro %d para JSON: %v", i, err)
			continue
		}
		log.Printf("📌 Membro %d:\n%s\n", i+1, string(jsonData))
	}

	return membros, nil
}
