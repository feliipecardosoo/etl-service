package inicialrepository

import bancoinicial "etl-service/src/config/model/banco_inicial"

// InicialRepository define a interface para o repositório que gerencia o acesso aos dados
// do banco inicial. Ela especifica os métodos necessários para buscar os membros da coleção.
type InicialRepository interface {
	// GetAllMembrosRequisicao retorna todos os membros encontrados na coleção do banco inicial.
	// Retorna um slice de Membro e um erro caso a operação falhe.
	GetAllMembrosRequisicao() ([]bancoinicial.Membro, error)
}
