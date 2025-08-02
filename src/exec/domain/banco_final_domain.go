package domain

import bancofinal "etl-service/src/config/model/banco_final"

// BancoFinalMembroDomain define a interface para entidades de membro
// que podem ser convertidas para o modelo bancofinal.Membro.
type BancoFinalMembroDomain interface {
	// ToModel converte a entidade de domínio para o modelo bancofinal.Membro,
	// que pode ser utilizado na camada de persistência ou outras camadas.
	ToModel() bancofinal.Membro
}
