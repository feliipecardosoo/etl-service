package bancoinicial

// Endereco representa os dados de endereço de um membro,
// contendo informações como CEP, rua, número, bairro e complemento.
//
// O campo Complemento é opcional e pode estar ausente na estrutura BSON.
type Endereco struct {
	Cep         string  `bson:"cep"`                   // Código postal
	Rua         string  `bson:"rua"`                   // Nome da rua
	Numero      string  `bson:"numero"`                // Número da residência
	Bairro      string  `bson:"bairro"`                // Nome do bairro
	Complemento *string `bson:"complemento,omitempty"` // Complemento do endereço (opcional)
}

// Membro representa as informações pessoais e de status de um membro da igreja.
//
// Os campos possuem tags BSON para mapear corretamente ao banco MongoDB.
//
// Alguns campos são opcionais (como DataCasamento e NomeConjuge) e podem estar ausentes.
type Membro struct {
	Name           string   `bson:"name"`                     // Nome completo do membro
	DataNascimento string   `bson:"data_nascimento"`          // Data de nascimento (formato string)
	AnoBatismo     string   `bson:"ano_batismo"`              // Ano em que foi batizado
	Sexo           string   `bson:"sexo"`                     // Sexo do membro
	EstadoCivil    string   `bson:"estado_civil"`             // Estado civil atual
	DataCasamento  string   `bson:"data_casamento,omitempty"` // Data do casamento (opcional)
	NomeConjuge    *string  `bson:"nome_conjuge,omitempty"`   // Nome do cônjuge (opcional)
	Filho          string   `bson:"filho"`                    // Indica se possui filhos
	Email          string   `bson:"email"`                    // E-mail de contato
	Telefone       string   `bson:"telefone"`                 // Telefone de contato
	Status         string   `bson:"status"`                   // Status do membro (ativo, inativo, etc)
	DataStatus     string   `bson:"data_status"`              // Data da última alteração de status
	Validado       bool     `bson:"validado"`                 // Indica se o cadastro foi validado
	Endereco       Endereco `bson:"endereco"`                 // Endereço completo do membro
}
