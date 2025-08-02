package domain

import (
	"errors"
	bancofinal "etl-service/src/config/model/banco_final"
	bancoinicial "etl-service/src/config/model/banco_inicial"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// membroDomain representa o domínio de um membro pronto para ser convertido.
// Contém os dados necessários para converter do modelo inicial para o modelo final.
type membroDomain struct {
	name            string          // Nome do membro
	dataNascimento  string          // Data de nascimento (string no formato original)
	anoBatismo      int             // Ano do batismo (convertido para int)
	sexo            string          // Sexo do membro
	estadoCivil     string          // Estado civil do membro
	dataCasamento   string          // Data do casamento (string)
	nomeConjuge     string          // Nome do cônjuge (string, vazio se não houver)
	filho           bool            // Indica se o membro tem filhos (true/false)
	email           string          // E-mail do membro
	telefone        string          // Telefone do membro
	status          string          // Status atual do membro
	dataStatus      string          // Data da última alteração do status
	endereco        enderecoRequest // Endereço do membro no formato interno do domínio
	validado        bool            // Indica se o cadastro foi validado
	dataAniversario string          // Data do aniversário no formato dia/mês (ex: "31/03")
}

// enderecoRequest representa o endereço usado internamente no domínio,
// armazenando os dados do endereço do membro.
type enderecoRequest struct {
	cep         string // Código postal
	rua         string // Nome da rua
	numero      string // Número da residência
	bairro      string // Bairro
	complemento string // Complemento do endereço (opcional)
}

// NewBancoFinalMembroDomain cria uma instância de membroDomain a partir de um membro do modelo inicial.
// Converte e trata campos específicos, como data de nascimento, ano de batismo e complementos.
// Retorna erro caso algum campo esteja em formato inválido.
func NewBancoFinalMembroDomain(m bancoinicial.Membro) (BancoFinalMembroDomain, error) {
	end := enderecoRequest{
		cep:         m.Endereco.Cep,
		rua:         m.Endereco.Rua,
		numero:      m.Endereco.Numero,
		bairro:      m.Endereco.Bairro,
		complemento: "",
	}

	if m.Endereco.Complemento != nil {
		end.complemento = *m.Endereco.Complemento
	}

	nomeConjuge := ""
	if m.NomeConjuge != nil {
		nomeConjuge = *m.NomeConjuge
	}

	// Formata a data de nascimento para o padrão dia/mês
	dataNascimentoFormatada, err := getAniversario(m.DataNascimento)
	if err != nil {
		return nil, fmt.Errorf("erro ao formatar dataNascimento '%s': %w", m.DataNascimento, err)
	}

	// Converte o ano de batismo para inteiro, se informado
	var dataBatismoFormatada int
	if m.AnoBatismo != "" {
		dataBatismoFormatada, err = getBatismo(m.AnoBatismo)
		if err != nil {
			return nil, fmt.Errorf("erro ao formatar dataBatismo '%s': %w", m.AnoBatismo, err)
		}
	} else {
		dataBatismoFormatada = 0
	}

	nameFormatado, err := getName(m.Name)
	if err != nil {
		return nil, fmt.Errorf("erro ao formatar name '%s': %w", m.Name, err)
	}

	return &membroDomain{
		name:            nameFormatado,
		dataNascimento:  m.DataNascimento,
		anoBatismo:      dataBatismoFormatada,
		sexo:            m.Sexo,
		estadoCivil:     m.EstadoCivil,
		dataCasamento:   m.DataCasamento,
		nomeConjuge:     nomeConjuge,
		filho:           m.Filho == "Sim",
		email:           m.Email,
		telefone:        m.Telefone,
		status:          m.Status,
		dataStatus:      m.DataStatus,
		endereco:        end,
		validado:        m.Validado,
		dataAniversario: dataNascimentoFormatada,
	}, nil
}

// ToModel converte o membroDomain para o modelo final bancofinal.Membro,
// pronto para ser utilizado na camada de repositório ou persistência.
func (m *membroDomain) ToModel() bancofinal.Membro {
	return bancofinal.Membro{
		Name:           m.name,
		DataNascimento: m.dataNascimento,
		AnoBatismo:     m.anoBatismo,
		Sexo:           m.sexo,
		EstadoCivil:    m.estadoCivil,
		DataCasamento:  m.dataCasamento,
		NomeConjuge:    &m.nomeConjuge,
		Filho:          m.filho,
		Email:          m.email,
		Telefone:       m.telefone,
		Status:         m.status,
		DataStatus:     m.dataStatus,
		Validado:       m.validado,
		Endereco: bancofinal.Endereco{
			Cep:         m.endereco.cep,
			Rua:         m.endereco.rua,
			Numero:      m.endereco.numero,
			Bairro:      m.endereco.bairro,
			Complemento: &m.endereco.complemento,
		},
		DataAniversario: m.dataAniversario,
	}
}

// getAniversario converte a string da data de nascimento no formato "YYYY-MM-DD"
// para uma string formatada "DD/MM". Retorna erro se o formato for inválido.
func getAniversario(dataNascimento string) (string, error) {
	t, err := time.Parse("2006-01-02", dataNascimento)
	if err != nil {
		return "", errors.New("Erro ao formatar dataNascimento (dentro da func getAniversario)")
	}
	return t.Format("02/01"), nil
}

// getBatismo converte o ano de batismo de string para inteiro.
// Retorna erro caso a conversão falhe.
func getBatismo(dataBatismo string) (int, error) {
	ano, err := strconv.Atoi(dataBatismo)
	if err != nil {
		return 0, errors.New("erro ao converter ano de batismo para inteiro")
	}
	return ano, nil
}

// getName converte o name para uppercase.
// Retorna erro caso a conversão falhe.
func getName(name string) (string, error) {
	if name == "" {
		return "", errors.New("nome vazio")
	}
	return strings.ToUpper(name), nil
}
