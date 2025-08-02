package getdata

import (
	"etl-service/src/exec/domain"
	inicialrepository "etl-service/src/exec/repository/inicial_repository"
	"fmt"
	"os"
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
// GetAll busca todos os membros, processa e cria arquivo txt com duplicados sem parar a execução
func (g *getDataBancoInicial) GetAll() error {
	start := time.Now()

	membros, err := g.repo.GetAllMembrosRequisicao()
	if err != nil {
		return fmt.Errorf("erro ao obter membros: %w", err)
	}

	// Extrai todos os nomes para consulta em lote
	var nomes []string
	for _, m := range membros {
		domainMembro, err := domain.NewBancoFinalMembroDomain(m)
		if err != nil {
			return err
		}
		model := domainMembro.ToModel()
		nomes = append(nomes, model.Name)
	}

	// Consulta no banco os nomes já existentes em lote
	existingMap, err := g.repo.ExistsByNames(nomes)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência dos membros: %w", err)
	}

	// Armazena nomes duplicados para gerar arquivo no fim
	var duplicados []string

	for _, m := range membros {
		domainMembro, err := domain.NewBancoFinalMembroDomain(m)
		if err != nil {
			return err
		}
		model := domainMembro.ToModel()

		if existingMap[model.Name] {
			// Acumula o nome do membro duplicado
			duplicados = append(duplicados, model.Name)
			// Continua a execução sem parar
			continue
		}

		// Inserir no banco novo aqui
		// err = g.repoFinal.Insert(model)
		// if err != nil {
		// 	 return err
		// }

	}

	// Se houver duplicados, cria arquivo txt com a lista
	if len(duplicados) > 0 {
		file, err := os.Create("duplicados.txt")
		if err != nil {
			return fmt.Errorf("erro ao criar arquivo de duplicados: %w", err)
		}
		defer file.Close()

		for _, nome := range duplicados {
			_, err := file.WriteString(nome + "\n")
			if err != nil {
				return fmt.Errorf("erro ao escrever no arquivo de duplicados: %w", err)
			}
		}

		fmt.Printf("Arquivo 'duplicados.txt' criado com %d nomes duplicados\n", len(duplicados))
	} else {
		fmt.Println("Nenhum membro duplicado encontrado.")
	}

	fmt.Printf("Membros totais: %d\n", len(membros))
	fmt.Printf("Tempo de execução: %s\n", time.Since(start))

	return nil
}
