# Projeto ETL 

## Visão Geral

Este projeto é um sistema ETL, com foco em:

- **Importação e tratamento de dados** vindos de uma base inicial.
- **Transformação** para um modelo final estruturado.
- **Validação e persistência** no MongoDB com esquema de validação JSON Schema.
- Geração de arquivos de log para registros duplicados e erros de inserção.

## Tecnologias e Estrutura

- Linguagem: **Go (Golang)**
- Banco de Dados: **MongoDB** (com validação por JSON Schema)
- Organização:  
  - **Modelos** para as estruturas iniciais (`bancoinicial`) e finais (`bancofinal`).
  - **Domínio** para conversão e tratamento dos dados.
  - **Repositório** para acesso e persistência no banco.
  - Goroutines para inserção concorrente e controle de erros.

## Funcionalidades Principais (No exemplo utilizei uma estrutura de membro fictícia e segui algumas regras de negócio)

### 1. Conversão do modelo inicial para final

- Utiliza o domínio para converter campos, tratando:
  - Formatação de datas (ex: nascimento para "DD/MM")
  - Conversão de tipos (ex: ano de batismo string → int)
  - Normalização (ex: nomes para uppercase)
  - Controle de ponteiros para campos opcionais (ex: `NomeConjuge`, `Complemento`)

### 2. Inserção concorrente no banco

- Insere membros em MongoDB usando pool limitado de goroutines (`maxWorkers = 10`).
- Registra nomes duplicados antes da inserção (evitando repetir registros).
- Captura erros de inserção para posterior análise.

### 3. Geração de arquivos de log

- Arquivo `duplicados.txt` para nomes já existentes.
- Arquivo `erros_insercao.txt` para erros no momento da inserção.
- Logs no console para sucesso e contagem de registros processados.

## Modelo MongoDB com validação JSON Schema

O documento `Membro` possui campos essenciais como:

- `name`, `dataNascimento`, `anoBatismo`, `sexo`, `status`, `dataStatus`, `validado`, `dataAniversario`, entre outros.
- O campo `endereco` é um **subdocumento** com:
  - `cep`, `rua`, `numero`, `bairro` (requeridos)
  - `complemento` (opcional)
- Todos os campos possuem tipagem correta com `bsonType`.
- A validação foi configurada para garantir a integridade dos dados.

Exemplo de schema JSON Schema para o MongoDB:

```json
{
  "$jsonSchema": {
    "bsonType": "object",
    "required": [
      "name",
      "dataNascimento",
      "anoBatismo",
      "sexo",
      "status",
      "dataStatus",
      "validado",
      "dataAniversario",
      "endereco"
    ],
    "properties": {
      "name": {"bsonType": "string"},
      "dataNascimento": {"bsonType": "string"},
      "anoBatismo": {"bsonType": "int"},
      "sexo": {"bsonType": "string"},
      "estadoCivil": {"bsonType": "string"},
      "dataCasamento": {"bsonType": "string"},
      "nomeConjuge": {"bsonType": "string"},
      "filho": {"bsonType": "bool"},
      "email": {"bsonType": "string"},
      "telefone": {"bsonType": "string"},
      "status": {"bsonType": "string"},
      "dataStatus": {"bsonType": "string"},
      "validado": {"bsonType": "bool"},
      "dataAniversario": {"bsonType": "string"},
      "dataModificacao": {"bsonType": "string"},
      "endereco": {
        "bsonType": "object",
        "required": ["cep", "rua", "numero", "bairro"],
        "properties": {
          "cep": {"bsonType": "string"},
          "rua": {"bsonType": "string"},
          "numero": {"bsonType": "string"},
          "bairro": {"bsonType": "string"},
          "complemento": {"bsonType": "string"}
        }
      }
    }
  }
}
```

## Estrutura do Código

- **domain/membro.go**: Define o domínio e conversão de dados.
- **bancoinicial/model.go**: Modelos brutos iniciais.
- **bancofinal/model.go**: Modelos finais com tags BSON para o MongoDB.

### Função `GetAll()` para:
- Buscar membros.
- Converter e validar dados.
- Inserir registros concorrentes.
- Gerar logs de duplicados e erros.

## Como Rodar

1. Configure o MongoDB aplicando o schema de validação.
2. Configure a conexão com o MongoDB no repositório Go.
3. Execute a função `GetAll()` para processar os membros.
4. Verifique os arquivos `duplicados.txt` e `erros_insercao.txt` para auditoria.

## Sistema de backup
- Possuo um sistema de backup deste banco no repositório: `https://github.com/feliipecardosoo/backup`
