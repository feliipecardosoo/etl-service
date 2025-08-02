# Etapa 1: Build da aplicação
FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod  ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o etl-service ./src

FROM alpine:3.20

# Adiciona certificado raiz (caso o app faça chamadas HTTPS)
RUN apk add --no-cache ca-certificates

# Define diretório de trabalho
WORKDIR /app

# Copia apenas o binário compilado
COPY --from=builder /app/etl-service .

# Comando padrão para rodar o app
CMD ["./etl-service"]
