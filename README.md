# 🎬 Movies API - Microserviços em Go

Uma arquitetura de microserviços robusta e escalável para gerenciar uma coleção de filmes. Implementada com **Go**, **Gin Framework**, **gRPC**, **MongoDB** e **Docker**.

## 📋 Visão Geral

Este projeto demonstra as melhores práticas de desenvolvimento de microserviços em Go, incluindo:

- ✅ **Arquitetura de Microserviços** - API Gateway + Movies Service
- ✅ **Comunicação gRPC** - Entre serviços
- ✅ **REST API** - Interface pública com Gin
- ✅ **Banco de Dados** - MongoDB
- ✅ **Containerização** - Docker e Docker Compose
- ✅ **Documentação** - Swagger/OpenAPI 3.0
- ✅ **Testes Unitários** - Testes automatizados
- ✅ **Logging Estruturado** - slog

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                     Cliente (HTTP)                          │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP/JSON
                       │ :8080
┌──────────────────────▼──────────────────────────────────────┐
│              API Gateway (Gin)                              │
│  - REST Endpoints                                           │
│  - Request/Response handling                               │
│  - gRPC Client                                             │
└──────────────────────┬──────────────────────────────────────┘
                       │ gRPC
                       │ :50051
┌──────────────────────▼──────────────────────────────────────┐
│         Movies Service (gRPC Server)                        │
│  - Business Logic                                           │
│  - Data Persistence                                         │
│  - Validations                                              │
└──────────────────────┬──────────────────────────────────────┘
                       │ Driver
                       │ :27017
┌──────────────────────▼──────────────────────────────────────┐
│              MongoDB Database                               │
│  - movies_db database                                       │
│  - Persistent storage                                       │
└─────────────────────────────────────────────────────────────┘
```

### Componentes

#### 1. **API Gateway** 🚪
- Framework: Gin (HTTP)
- Porta: `8080`
- Função: Expõe endpoints REST para clientes
- Comunica com Movies Service via gRPC
- Localização: `api-gateway/`

#### 2. **Movies Service** 🎥
- Protocolo: gRPC
- Porta: `50051`
- Função: Lógica de negócio e persistência
- Conecta ao MongoDB
- Localização: `movies-service/`

#### 3. **Proto Definitions** 📝
- Define as mensagens gRPC
- Define o serviço gRPC
- Localização: `proto/`

#### 4. **MongoDB** 🗄️
- Versão: 8.0
- Banco: `movies_db`
- Porta: `27017`
- Username: `root`
- Password: `password`

---

## 🚀 Início Rápido

### Pré-requisitos

- **Go** 1.21.0 ou superior
- **Docker** e **Docker Compose**
- **Git**
- Opcionais:
  - `swag` para regenerar documentação Swagger
  - `protoc` para recompilar proto files

### Instalação

#### 1. Clonar o repositório

```bash
git clone https://github.com/FranciscoHonorat/movies.git
cd movies
```

#### 2. Instalar dependências Go

```bash
# Baixar todas as dependências
go mod download

# (Opcional) Atualizar dependências
go mod tidy
```

#### 3. Dependências estão prontas

As variáveis de ambiente estão configuradas no `docker-compose.yml`.

### Executar com Docker Compose

#### Build e Start dos serviços

```bash
# Build de todas as imagens
docker compose build

# Start dos serviços
docker compose up -d

# Ver logs em tempo real
docker compose logs -f

# Parar serviços
docker compose down
```

#### Verificar se está rodando

```bash
# Health check
curl http://localhost:8080/health

# Listar filmes
curl http://localhost:8080/api/v1/movies
```

### Executar Localmente (sem Docker)

#### 1. Iniciar MongoDB

```bash
# Docker
docker run -d \
  --name movies-mongo \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=root \
  -e MONGO_INITDB_ROOT_PASSWORD=password \
  mongo:8.0
```

#### 2. Iniciar Movies Service

```bash
cd movies-service
go run cmd/main.go
```

#### 3. Iniciar API Gateway

```bash
cd api-gateway
GRPC_SERVER_URL=localhost:50051 go run cmd/main.go
```

Agora acesse: `http://localhost:8080`

---

## 📁 Estrutura do Projeto

```
movies/
├── api-gateway/                    # API Gateway (HTTP/Gin)
│   ├── cmd/
│   │   └── main.go                # Entry point
│   ├── internal/
│   │   └── handlers/
│   │       ├── movie-handler.go    # Movie CRUD handlers
│   │       ├── health.go           # Health check
│   │       └── errors.go           # Error handling
│   ├── docs/
│   │   ├── swagger.yaml            # OpenAPI spec
│   │   ├── swagger.json            # OpenAPI JSON
│   │   └── docs.go                 # Swagger metadata
│   ├── go.mod
│   └── Dockerfile
│
├── movies-service/                 # Movies Service (gRPC)
│   ├── cmd/
│   │   └── main.go                # Entry point
│   ├── internal/
│   │   ├── adapters/
│   │   │   ├── grpc_server/       # gRPC server implementation
│   │   │   │   ├── server.go
│   │   │   │   └── errors.go
│   │   │   ├── mongodb/            # MongoDB adapter
│   │   │   │   └── movie_repo.go
│   │   │   └── seed/               # Data seeding
│   │   │       └── seed.go
│   │   └── core/
│   │       ├── domain/
│   │       │   ├── movie.go        # Movie entity
│   │       │   └── errors.go
│   │       ├── ports/
│   │       │   ├── input/          # Input ports (service interface)
│   │       │   └── output/         # Output ports (repository interface)
│   │       └── service/
│   │           └── movie_service.go # Business logic
│   ├── movies.json                 # Sample data
│   ├── go.mod
│   └── Dockerfile
│
├── proto/                           # Protocol Buffer definitions
│   ├── movies.proto                # gRPC service definition
│   ├── movies_grpc.pb.go           # Generated gRPC code
│   ├── movies.pb.go                # Generated message code
│   ├── go.mod
│   └── go.sum
│
├── docker-compose.yml              # Orquestração dos serviços
├── go.work                         # Workspace do Go
├── Dockerfile                      # Dockerfile raiz (se existir)
├── .dockerignore                   # Docker build ignore
├── README.md                       # Este arquivo
└── go.work.sum
```

---

## 🔌 API Endpoints

### Base URL
```
http://localhost:8080
```

### Health Check

```http
GET /health
```

**Resposta 200:**
```json
{
  "status": "ok",
  "service": "api-gateway"
}
```

### Movies API (v1)

#### Obter filme por ID

```http
GET /api/v1/movies/{id}
```

**Exemplo:**
```bash
curl http://localhost:8080/api/v1/movies/1
```

**Resposta 200:**
```json
{
  "id": 1,
  "title": "The Shawshank Redemption",
  "year": "1994"
}
```

#### Listar filmes

```http
GET /api/v1/movies?title=&year=&page=&limit=&sort=
```

**Query Parameters:**
- `title` (string, optional) - Buscar por título
- `year` (string, optional) - Filtrar por ano
- `page` (int, optional) - Número da página (padrão: 1)
- `limit` (int, optional) - Filmes por página (padrão: 10, máximo: 100)
- `sort` (string, optional) - Campo de ordenação: `title` ou `year` (padrão: title)

**Exemplo:**
```bash
curl "http://localhost:8080/api/v1/movies?title=The&year=1994&page=1&limit=10&sort=title"
```

**Resposta 200:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "The Shawshank Redemption",
      "year": "1994"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1
}
```

#### Criar filme

```http
POST /api/v1/movies
Content-Type: application/json

{
  "title": "Inception",
  "year": "2010"
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "Inception", "year": "2010"}'
```

**Resposta 201:**
```json
{
  "id": 251,
  "title": "Inception",
  "year": "2010"
}
```

#### Deletar filme

```http
DELETE /api/v1/movies/{id}
```

**Exemplo:**
```bash
curl -X DELETE http://localhost:8080/api/v1/movies/1
```

**Resposta 204:** (Sem corpo)

---

## 📚 Documentação Swagger

A API está totalmente documentada com **Swagger/OpenAPI 3.0**.

### Visualizar Documentação

#### Opção 1: Markdown

A documentação completa da API está disponível online no repositório.

#### Opção 2: Swagger UI (Local)

Após rodar `swag init -g cmd/main.go` na pasta `api-gateway/`:

```bash
cd api-gateway
go run cmd/main.go
```

Acesse: `http://localhost:8080/swagger/index.html`

#### Opção 3: Editor Online
1. Acesse https://editor.swagger.io/
2. Cole o conteúdo de `api-gateway/docs/swagger.yaml`

#### Opção 4: Postman

Importe a coleção Postman disponível no repositório para testes interativos dos endpoints.

---

## 🛠️ Desenvolvimento

### Setup do Ambiente de Desenvolvimento

```bash
# 1. Clonar repositório
git clone https://github.com/FranciscoHonorat/movies.git
cd movies

# 2. Instalar dependências
go mod download

# 3. Setup Docker (MongoDB)
docker compose up -d mongo

# 4. Rodar os serviços
cd movies-service && go run cmd/main.go &
cd ../api-gateway && GRPC_SERVER_URL=localhost:50051 go run cmd/main.go
```

### Modificar Proto Definitions

Se alterar `proto/movies.proto`:

```bash
cd proto
protoc --go_out=. --go-grpc_out=. movies.proto
```

### Atualizar Documentação Swagger

Após adicionar/modificar anotações nos handlers:

```bash
cd api-gateway

# Instalar swaggo (primeira vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar documentação
swag init -g cmd/main.go
```

### Estrutura de Código

#### Domain-Driven Design
O projeto segue princípios de **Clean Architecture**:

```
├── domain/       → Entities e Rules (sem dependências)
├── ports/        → Interfaces (input/output)
├── service/      → Business Logic (usa ports)
└── adapters/     → Implementações (MongoDB, gRPC)
```

#### Padrão de Handlers (API Gateway)

Cada handler segue este padrão:
1. Recebe request do Gin
2. Valida parâmetros
3. Chama client gRPC
4. Trata erros e converte para HTTP
5. Retorna JSON

---

## 🧪 Testes

### Rodar Testes

```bash
# Todos os testes (recursivo)
go test ./...

# Testes de um módulo específico
go test ./internal/adapters/grpc_server/

# Com cobertura
go test -cover ./...

# Com verbose
go test -v ./...

# Teste específico
go test -run TestName ./...
```

### Estrutura de Testes

Testes estão no mesmo diretório dos arquivos originais com sufixo `_test.go`:

```
internal/
├── adapters/
│   └── grpc_server/
│       ├── server.go
│       └── errors_test.go       ← Testes
└── core/
    └── service/
        └── movie_service.go
```

### Exemplo de Teste

```go
func TestGetMovie(t *testing.T) {
    // Arrange
    expected := &Movie{ID: 1, Title: "Test"}
    
    // Act
    result, err := service.GetMovie(ctx, 1)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

---

## 🐳 Docker e Docker Compose

### Arquivos Docker

```
api-gateway/Dockerfile      # Build do API Gateway
movies-service/Dockerfile   # Build do Movies Service
docker-compose.yml          # Orquestração
```

### Variáveis de Ambiente (docker-compose)

```yaml
GRPC_SERVER_URL: movies-service:50051   # API Gateway → Movies Service
MONGODB_URI: mongodb://root:password@mongo:27017/movies_db
```

### Comandos Úteis

```bash
# Build sem cache
docker compose build --no-cache

# Ver logs de um serviço
docker compose logs -f api-gateway
docker compose logs -f movies-service
docker compose logs -f mongo

# Executar comando em container
docker compose exec movies-service go test ./...

# Remover tudo (containers, networks, volumes)
docker compose down -v

# Rebuild e start
docker compose up -d --build
```

---

## ☸️ Kubernetes

### Pré-requisitos

- `kubectl` instalado e configurado
- Cluster Kubernetes rodando (Docker Desktop, Minikube, Kind, etc.)
- Imagens Docker dos serviços disponíveis (pushed para registry ou carregadas localmente)

### Manifests

Os manifests Kubernetes estão localizados em `k8s/`:

```
k8s/
├── api-gateway-deployment.yaml      # Deployment do API Gateway
├── api-gateway-service.yaml         # Service do API Gateway
├── movies-service-deployment.yaml   # Deployment do Movies Service
├── movies-service-service.yaml      # Service do Movies Service
├── mongodb-deployment.yaml          # Deployment do MongoDB
├── mongodb-service.yaml             # Service do MongoDB
└── mongodb-pvc.yaml                 # PersistentVolumeClaim para MongoDB
```

### Deploy no Kubernetes

#### 1. Aplicar todos os manifests

```bash
kubectl apply -f k8s/
```

Este comando cria:
- Deployments para API Gateway, Movies Service e MongoDB
- Services para exposição dos serviços
- PersistentVolumeClaim para persistência de dados do MongoDB

#### 2. Verificar o status

```bash
# Ver pods
kubectl get pods

# Ver services
kubectl get svc

# Ver deployments
kubectl get deployments

# Ver volumes
kubectl get pvc
```

#### 3. Acessar os serviços

```bash
# Port-forward API Gateway (local)
kubectl port-forward svc/api-gateway 8080:8080

# Port-forward MongoDB (local)
kubectl port-forward svc/mongodb 27017:27017

# Acessar Swagger
http://localhost:8080/swagger/index.html

# Acessar API
http://localhost:8080/api/v1/movies
```

#### 4. Ver logs

```bash
# Logs do API Gateway
kubectl logs -l app=api-gateway -f

# Logs do Movies Service
kubectl logs -l app=movies-service -f

# Logs do MongoDB
kubectl logs -l app=mongodb -f
```

#### 5. Deletar recursos

```bash
# Remover tudo
kubectl delete -f k8s/

# Ou remover seletivamente
kubectl delete deployment api-gateway
kubectl delete service api-gateway
```

### Configuração de Imagens

Antes de aplicar os manifests, certifique-se que as imagens estão disponíveis:

```bash
# Build das imagens
docker build -t api-gateway:latest api-gateway/
docker build -t movies-service:latest movies-service/

# Se usar um registry (ex: Docker Hub)
docker tag api-gateway:latest seu-usuario/api-gateway:latest
docker push seu-usuario/api-gateway:latest

docker tag movies-service:latest seu-usuario/movies-service:latest
docker push seu-usuario/movies-service:latest
```

### Variáveis de Ambiente no Kubernetes

Os Deployments usam ConfigMaps e Secrets para variáveis de ambiente:

```yaml
# Exemplo no deployment
env:
  - name: GRPC_SERVER_URL
    value: "movies-service:50051"
  - name: MONGODB_URI
    value: "mongodb://root:password@mongodb:27017/movies_db"
```

---

## 🔍 Troubleshooting

### Problema: Conexão recusada com MongoDB

**Erro:**
```
connection refused mongodb://root:password@mongo:27017
```

**Solução:**
```bash
# Verificar se MongoDB está rodando
docker compose logs mongo

# Reiniciar MongoDB
docker compose restart mongo

# Ou rodar manualmente
docker run -d \
  --name movies-mongo \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=root \
  -e MONGO_INITDB_ROOT_PASSWORD=password \
  mongo:8.0
```

### Problema: gRPC connection refused

**Erro:**
```
rpc error: code = Unavailable desc = connection refused
```

**Solução:**
```bash
# Verificar se Movies Service está rodando
docker compose logs movies-service

# Ou rodar localmente
cd movies-service
go run cmd/main.go

# Verificar porta
lsof -i :50051
```

### Problema: Swagger não encontrado

**Erro:**
```
404 Not Found /swagger/index.html
```

**Solução:**
```bash
cd api-gateway
swag init -g cmd/main.go
go run cmd/main.go
```

### Problema: "no Go files" ao executar swag

**Erro:**
```
error: execute go list command, exit status 1, stdout:, stderr:no Go files in .../api-gateway
```

**Solução:**
```bash
# Execute na raiz do projeto (onde está o dockerfile)
cd api-gateway
swag init -g cmd/main.go
```

### Problema: Erro ao conectar ao MongoDB (Docker for Desktop no Windows)

**Erro:**
```
connection refused
```

**Solução:**
```bash
# Usar "host.docker.internal" ao invés de "localhost"
# Ou usar o nome do container: "mongo"

# Em docker-compose, sempre use o nome do serviço:
mongodb://root:password@mongo:27017
```

---

## 📊 Monitoramento e Logs

### Logs Estruturados

O projeto usa `log/slog` para logging estruturado:

```go
slog.Error("CreateMovie error", slog.Any("error", err))
slog.Info("Movie created", slog.Any("movie_id", id))
```

### Ver Logs

```bash
# Todos os logs
docker compose logs

# Apenas API Gateway
docker compose logs -f api-gateway

# Últimas 50 linhas
docker compose logs --tail=50

# Com timestamps
docker compose logs -t
```

---

## 🚀 Deploy em Produção

### Checklist de Deploy

- [ ] Variáveis de ambiente configuradas (.env)
- [ ] MongoDB backup configurado
- [ ] TLS/HTTPS habilitado
- [ ] Autenticação implementada (JWT)
- [ ] Rate limiting configurado
- [ ] CORS configurado
- [ ] Logging centralizado
- [ ] Monitoring e alertas
- [ ] CI/CD pipeline
- [ ] Documentação atualizada

### Sugestões para Produção

```dockerfile
# Usar multi-stage build
FROM golang:1.25 AS builder
# ... build

FROM alpine:latest
# ... runtime
```

---

## 🤝 Contribuindo

### Steps para Contribuir

1. Fork o repositório
2. Criar feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Abrir Pull Request

### Padrões de Código

- Seguir `gofmt` para formatação
- Adicionar testes para novas funcionalidades
- Atualizar documentação
- Adicionar anotações Swagger para novos endpoints

---

## 📞 Contato e Suporte

- **Email**: support@movies-api.local
- **GitHub Issues**: [Abrir issue](https://github.com/FranciscoHonorat/movies/issues)
- **Discussions**: [GitHub Discussions](https://github.com/FranciscoHonorat/movies/discussions)

---

## 🙏 Agradecimentos

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [gRPC](https://grpc.io/)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [Swagger/OpenAPI](https://swagger.io/)

---

## 📚 Recursos Adicionais

### Documentação

A documentação completa da API está disponível no repositório, incluindo detalhes de todos os endpoints, parâmetros e exemplos de uso.

### Ferramentas Recomendadas
- [Postman](https://www.postman.com/) - Testes de API
- [MongoDB Compass](https://www.mongodb.com/products/compass) - Visualizador MongoDB
- [VS Code REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) - Testes inline
- [Insomnia](https://insomnia.rest/) - Cliente REST alternativo

### Cursos e Tutoriais
- [Go by Example](https://gobyexample.com/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [gRPC Concepts](https://grpc.io/docs/what-is-grpc/)
- [MongoDB University](https://university.mongodb.com/)

---

**Última atualização**: 24 de Maio de 2026  
**Versão da API**: 1.0.0  
**Versão Go**: 1.25.0  
**Status**: ✅ Pronto para uso
