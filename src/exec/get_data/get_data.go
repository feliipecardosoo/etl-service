package getdata

import bancoinicial "etl-service/src/config/model/banco_inicial"

// GetDataBancoInicial define a interface do serviço responsável por operações
// relacionadas à obtenção de dados do banco inicial.
// Essa interface abstrai as operações para facilitar a testabilidade e a troca da implementação.
type GetDataBancoInicial interface {
	// GetAll retorna uma lista de todos os membros presentes no banco inicial.
	// Retorna um slice de Membro e um erro caso ocorra algum problema durante a busca.
	GetAll() ([]bancoinicial.Membro, error)
}
