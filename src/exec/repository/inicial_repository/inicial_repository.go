package inicialrepository

import (
	bancoinicial "etl-service/src/config/model/banco_inicial"
)

// InicialRepository define a interface para o repositório que gerencia o acesso aos dados
// do banco inicial.
//
// Essa interface abstrai as operações de leitura necessárias para manipular membros da coleção,
// permitindo a implementação flexível da persistência, como MongoDB, PostgreSQL, entre outros.
type InicialRepository interface {
	// GetAllMembrosRequisicao busca e retorna todos os membros existentes na coleção do banco inicial.
	//
	// Retorna:
	// - Um slice contendo todos os membros encontrados (bancoinicial.Membro).
	// - Um erro caso a operação falhe, seja por problemas de conexão, timeout ou falha na consulta.
	GetAllMembrosRequisicao() ([]bancoinicial.Membro, error)

	// ExistsByNames verifica quais nomes dentre os passados existem atualmente no banco.
	//
	// Parâmetro:
	// - names: slice de strings contendo os nomes a serem consultados.
	//
	// Retorna:
	// - Um mapa string->bool indicando quais nomes existem no banco (true significa que o nome existe).
	// - Um erro caso ocorra falha durante a consulta.
	ExistsByNames(names []string) (map[string]bool, error)
}
