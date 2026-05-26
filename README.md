# 🎬 Movies API — Microserviços em Go

API REST de filmes construída com arquitetura de microserviços em **Go**, usando **Gin**, **gRPC + Protobuf**, **MongoDB**, **DynamoDB (LocalStack)**, **RabbitMQ** e **Docker**.

---

## 📋 Visão Geral

O projeto implementa um CRUD de filmes distribuído em dois microserviços que se comunicam via gRPC. Além dos requisitos obrigatórios, foram implementados os três extras opcionais do desafio:

- ✅ **Arquitetura Hexagonal** — domain, ports e adapters isolados
- ✅ **Microserviços** — API Gateway (HTTP) + Movies Service (gRPC)
- ✅ **gRPC + Protobuf** — comunicação tipada entre serviços
- ✅ **MongoDB** — persistência principal
- ✅ **Docker Compose** — inicialização com um único comando
- ✅ **Swagger** — documentação interativa em `/swagger/index.html`
- ✅ **Testes Unitários** — 20 testes (domain sem mock, service com mock, gRPC table-driven)
- ✅ **Event Driven (RabbitMQ)** — `POST /movies` assíncrono com retorno `202 Accepted`
- ✅ **Kubernetes** — 7 manifestos prontos em `k8s/`
- ✅ **LocalStack + DynamoDB** — substituto emulado do MongoDB via AWS SDK

---

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────┐
│                  Cliente (HTTP)                 │
└────────────────────┬────────────────────────────┘
                     │ HTTP :8080
┌────────────────────▼────────────────────────────┐
│              API Gateway (Gin)                  │
│  - REST endpoints                               │
│  - gRPC client                                  │
│  - RabbitMQ producer (POST assíncrono)          │
└──────────┬──────────────────────┬───────────────┘
           │ gRPC :50051          │ AMQP :5672
           │                     │ (apenas POST)
┌──────────▼──────────┐  ┌───────▼───────────────┐
│   Movies Service    │  │       RabbitMQ         │
│   (gRPC Server)     │◄─┤  fila: movies.create  │
│   - regras negócio  │  └───────────────────────┘
│   - RabbitMQ consumer                           │
└──────────┬──────────┘
           │ Driver :27017 / :4566
┌──────────▼──────────┐
│  MongoDB / DynamoDB │
│  (LocalStack)       │
└─────────────────────┘
```

### Como o fluxo assíncrono funciona

Os endpoints `GET`, `GET /{id}` e `DELETE` são **síncronos** — o api-gateway chama o movies-service via gRPC e aguarda a resposta.

O `POST /movies` é **assíncrono**:
1. O api-gateway publica a mensagem na fila `movies.create` do RabbitMQ
2. Retorna `202 Accepted` imediatamente (sem esperar o banco)
3. O movies-service consome a fila em background (goroutine) e persiste o filme

---

## 🗂️ Estrutura do Projeto

```
movies/
├── go.work                          # Workspace Go — conecta os módulos localmente
├── go.work.sum
├── docker-compose.yml               # Sobe tudo com um comando
├── .dockerignore
│
├── proto/                           # Módulo compartilhado de contrato gRPC
│   ├── movies.proto                 # Definição do serviço e mensagens
│   ├── movies.pb.go                 # Gerado pelo protoc
│   ├── movies_grpc.pb.go            # Gerado pelo protoc
│   └── go.mod
│
├── shared/                          # DTO compartilhado para mensageria
│   ├── movie_message.go             # MoviePublisherMessage (RabbitMQ)
│   └── go.mod
│
├── api-gateway/                     # Serviço HTTP — porta de entrada
│   ├── cmd/main.go
│   ├── internal/
│   │   └── adapters/
│   │       ├── handlers/
│   │       │   ├── movie_handler.go
│   │       │   └── health.go
│   │       ├── errors.go            # Mapeia gRPC codes → HTTP status
│   │       └── rabbitmq/
│   │           └── publisher.go     # Publica mensagens na fila
│   ├── docs/                        # Swagger gerado pelo swag
│   ├── go.mod
│   └── Dockerfile
│
├── movies-service/                  # Serviço de negócio — gRPC + banco
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── core/
│   │   │   ├── domain/
│   │   │   │   ├── movie.go         # Entidade Movie + validação
│   │   │   │   └── errors.go        # ErrMovieNotFound, ErrInvalidMovieData...
│   │   │   ├── ports/
│   │   │   │   ├── input/           # Interface MovieService (o que o service expõe)
│   │   │   │   └── output/          # Interface MovieRepository (o que o service precisa)
│   │   │   └── service/
│   │   │       └── movie_service.go # Lógica de negócio (não conhece MongoDB nem gRPC)
│   │   └── adapters/
│   │       ├── grpc_server/
│   │       │   ├── server.go        # Implementa o servidor gRPC
│   │       │   ├── errors.go        # Mapeia domain errors → gRPC codes
│   │       │   └── errors_test.go   # Testes table-driven do mapeamento de erros
│   │       ├── mongodb/
│   │       │   └── movie_repo.go    # Implementa MovieRepository com MongoDB
│   │       ├── dynamodb/
│   │       │   └── movie_repo.go    # Implementa MovieRepository com DynamoDB
│   │       ├── seed/
│   │       │   ├── seed.go          # Seed MongoDB (movies.json → banco)
│   │       │   └── dynamodb_seed.go # Seed DynamoDB (movies.json → banco)
│   │       └── rabbitmq/
│   │           └── consumer.go      # Consome fila movies.create em background
│   ├── movies.json                  # Dados iniciais (250 filmes)
│   ├── go.mod
│   └── Dockerfile
│
└── k8s/                             # Manifestos Kubernetes
    ├── api-gateway-deployment.yaml
    ├── api-gateway-service.yaml
    ├── movies-service-deployment.yaml
    ├── movies-service-service.yaml
    ├── mongodb-deployment.yaml
    ├── mongodb-service.yaml
    └── mongodb-pvc.yaml
```

---

## 🚀 Inicialização

### Pré-requisitos

- **Docker** e **Docker Compose**
- **Go 1.25+** (apenas se quiser rodar localmente sem Docker)

### Subir com Docker Compose (recomendado)

```bash
git clone https://github.com/FranciscoHonorat/movies.git
cd movies

# Build dos serviços (necessário na primeira vez)
docker compose build movies-service
docker compose build api-gateway

# Subir tudo
docker compose up
```

> **Por que buildar separadamente?** O contexto de build do movies-service inclui o diretório `vendor/` (~84MB). Buildar em paralelo pode causar EOF no Docker Desktop ao transferir dois contextos grandes ao mesmo tempo.

Após subir, os serviços ficam disponíveis em:
- API REST: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- MongoDB: `localhost:27017`
- RabbitMQ Management: `http://localhost:15672` (guest/guest)

O banco é populado automaticamente com os filmes do `movies.json` na primeira inicialização.

### Parar os serviços

```bash
docker compose down
```

---

## 🔌 Endpoints

### Health Check

```bash
curl http://localhost:8080/health
```

```json
{ "status": "ok", "service": "api-gateway" }
```

### GET /api/v1/movies

Lista filmes com suporte a filtros, paginação e ordenação.

**Parâmetros (todos opcionais):**

| Parâmetro | Tipo   | Padrão  | Descrição                          |
|-----------|--------|---------|------------------------------------|
| `title`   | string | —       | Filtro por título (busca parcial)  |
| `year`    | string | —       | Filtro por ano                     |
| `page`    | int    | 1       | Página                             |
| `limit`   | int    | 10      | Itens por página (máx: 100)        |
| `sort`    | string | `title` | Campo de ordenação: `title`, `year`|

```bash
curl "http://localhost:8080/api/v1/movies?title=The&page=1&limit=5&sort=year"
```

```json
{
  "data": [
    { "id": 1, "title": "The Shawshank Redemption", "year": "1994" }
  ],
  "page": 1,
  "limit": 5,
  "total": 1
}
```

### GET /api/v1/movies/{id}

Busca um filme pelo ID.

```bash
curl http://localhost:8080/api/v1/movies/1
```

```json
{ "id": 1, "title": "The Shawshank Redemption", "year": "1994" }
```

Erros possíveis:
- `400 Bad Request` — ID inválido (não é número)
- `404 Not Found` — filme não encontrado

### POST /api/v1/movies

Cria um novo filme de forma **assíncrona**. Publica na fila RabbitMQ e retorna `202 Accepted` imediatamente.

```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "Inception", "year": "2010"}'
```

```json
HTTP/1.1 202 Accepted
{ "message": "movie creation accepted" }
```

> O filme será persistido em background pelo movies-service assim que consumir a mensagem da fila.

Erros possíveis:
- `400 Bad Request` — body inválido ou campos obrigatórios ausentes

### DELETE /api/v1/movies/{id}

Remove um filme pelo ID.

```bash
curl -X DELETE http://localhost:8080/api/v1/movies/1
```

`204 No Content` (sem body)

Erros possíveis:
- `400 Bad Request` — ID inválido
- `404 Not Found` — filme não encontrado

---

## 🧪 Testes

```bash
# Todos os testes do workspace
go test ./...

# Testes de um módulo específico
cd movies-service && go test ./...

# Com cobertura
go test -cover ./...

# Verbose
go test -v ./...
```

### O que está coberto

- **Domain (3 testes, sem mock)** — validação de `NewMovie`, erros de título vazio, ano inválido
- **Service (13 testes, com mock)** — todos os casos de sucesso e erro dos 4 métodos do CRUD
- **gRPC adapter (4 testes, table-driven)** — mapeamento de domain errors para gRPC status codes

Os testes de service usam um `MockMovieRepository` que implementa a interface `output.MovieRepository`. Isso significa que os 13 testes passam sem MongoDB, DynamoDB ou qualquer banco rodando.

---

## 📚 Swagger

A documentação interativa fica disponível após subir os serviços:

```
http://localhost:8080/swagger/index.html
```

Para regenerar após alterar anotações nos handlers:

```bash
cd api-gateway
swag init -g cmd/main.go
```

---

## ☸️ Kubernetes

Os manifestos estão em `k8s/` e foram testados com Docker Desktop (Kubernetes habilitado).

```bash
# Buildar as imagens localmente antes de aplicar
docker build -t api-gateway:latest -f api-gateway/Dockerfile .
docker build -t movies-service:latest -f movies-service/Dockerfile .

# Aplicar todos os manifestos
kubectl apply -f k8s/

# Verificar status
kubectl get pods
kubectl get services

# Acessar a API via port-forward
kubectl port-forward svc/api-gateway 8080:8080
```

Os manifestos criam Deployments e Services para api-gateway, movies-service e MongoDB, além de um PersistentVolumeClaim para persistência dos dados do MongoDB.

---

## 🔄 LocalStack + DynamoDB

O movies-service tem dois adapters que implementam a mesma interface `MovieRepository`: um para MongoDB e outro para DynamoDB via LocalStack. Trocar entre eles exige mudar uma única linha no `main.go` — o service, o gRPC server e os testes não precisam de nenhuma alteração.

O LocalStack emula o DynamoDB localmente sem custo e sem necessidade de conta AWS. As credenciais configuradas no `docker-compose.yml` são fake — o LocalStack não as valida, mas o AWS SDK exige que existam.

O `docker-compose.yml` já inclui o container do LocalStack. Para usar o DynamoDB no lugar do MongoDB, basta trocar o adapter no `main.go` do movies-service:

```go
// MongoDB (padrão)
repo := mongodb.NewMongoRepository(...)

// DynamoDB via LocalStack
repo := dynamodb.NewDynamoRepository(...)
```

Depois rode normalmente:

```bash
docker compose up
```

---

## 🐛 Troubleshooting

**movies-service demora para ficar disponível no primeiro start**

O seed popula o banco sincronamente antes do servidor gRPC iniciar. Com 250 filmes no DynamoDB (um `PutItem` por vez), isso pode levar alguns segundos. Os logs mostram `PutItem => 200` enquanto o seed roda — é comportamento esperado.

**EOF durante `docker compose up --build`**

Buildar os serviços separadamente resolve:
```bash
docker compose build movies-service
docker compose build api-gateway
docker compose up
```

**Pods com `ErrImageNeverPull` no Kubernetes**

As imagens precisam ser buildadas localmente antes de aplicar os manifestos:
```bash
docker build -t api-gateway:latest -f api-gateway/Dockerfile .
docker build -t movies-service:latest -f movies-service/Dockerfile .
kubectl rollout restart deployment/api-gateway deployment/movies-service
```

**LocalStack falhando com `License activation failed`**

O `docker-compose.yml` já usa `localstack/localstack:3.8`. Versões mais recentes exigem licença paga — não altere a tag.

---

## 📞 Contato

- **Email**: jeffhonorato230@gmail.com
- **GitHub**: [github.com/FranciscoHonorat](https://github.com/FranciscoHonorat)

---

**Versão Go**: 1.25.0 | **Versão da API**: 1.0.0