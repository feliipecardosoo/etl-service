package finalrepository

import bancofinal "etl-service/src/config/model/banco_final"

type FinalRepository interface {
	Insert(membro bancofinal.Membro) error
}
