package getdata

import (
	bancofinal "etl-service/src/config/model/banco_final"
	"etl-service/src/exec/domain"
	inicialrepository "etl-service/src/exec/repository/inicial_repository"
	"fmt"
	"os"
	"sync"
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

	var nomes []string
	for _, m := range membros {
		domainMembro, err := domain.NewBancoFinalMembroDomain(m)
		if err != nil {
			return err
		}
		model := domainMembro.ToModel()
		nomes = append(nomes, model.Name)
	}

	existingMap, err := g.repo.ExistsByNames(nomes)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência dos membros: %w", err)
	}

	var duplicados []string
	type insertResult struct {
		name string
		err  error
	}

	// Canal para resultados das inserções
	resultCh := make(chan insertResult, len(membros))
	var wg sync.WaitGroup

	// Limitar número de goroutines concorrentes (pool)
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)

	for _, m := range membros {
		domainMembro, err := domain.NewBancoFinalMembroDomain(m)
		if err != nil {
			return err
		}
		model := domainMembro.ToModel()

		if existingMap[model.Name] {
			duplicados = append(duplicados, model.Name)
			continue
		}

		wg.Add(1)
		sem <- struct{}{} // bloqueia se maxWorkers estiver cheio
		go func(mod bancofinal.Membro) {
			defer wg.Done()

			err := g.repo.Insert(mod)
			resultCh <- insertResult{name: mod.Name, err: err}
			<-sem // libera vaga no pool
		}(model)
	}

	// Aguarda todas as goroutines terminarem
	wg.Wait()
	close(resultCh)

	// Coleta erros de inserção para arquivo
	var errosInsercao []string
	for res := range resultCh {
		if res.err != nil {
			errosInsercao = append(errosInsercao, fmt.Sprintf("%s: %v", res.name, res.err))
		}
	}

	// Grava duplicados num arquivo txt
	if len(duplicados) > 0 {
		err := writeLinesToFile("duplicados.txt", duplicados)
		if err != nil {
			return fmt.Errorf("erro ao criar arquivo de duplicados: %w", err)
		}
		fmt.Printf("Arquivo 'duplicados.txt' criado com %d nomes duplicados\n", len(duplicados))
	} else {
		fmt.Println("Nenhum membro duplicado encontrado.")
	}

	// Grava erros de inserção num arquivo txt
	if len(errosInsercao) > 0 {
		err := writeLinesToFile("erros_insercao.txt", errosInsercao)
		if err != nil {
			return fmt.Errorf("erro ao criar arquivo de erros de inserção: %w", err)
		}
		fmt.Printf("Arquivo 'erros_insercao.txt' criado com %d erros de inserção\n", len(errosInsercao))
	} else {
		fmt.Println("Nenhum erro de inserção encontrado.")
	}

	fmt.Printf("Membros totais processados: %d\n", len(membros))
	fmt.Printf("Tempo de execução: %s\n", time.Since(start))

	return nil
}

// writeLinesToFile grava uma slice de strings em arquivo, uma linha por string
func writeLinesToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
