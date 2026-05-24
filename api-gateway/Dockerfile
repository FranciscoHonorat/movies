# Build stage
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Instalar dependências
RUN apk add --no-cache git

# Copiar código fonte
COPY . .

# Download dependências
RUN go mod download

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./api-gateway/cmd

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar executável do builder
COPY --from=builder /app/server .

# Expor porta HTTP
EXPOSE 8080

# Variáveis de ambiente
ENV GRPC_SERVER_URL=movies-service:50051

CMD ["./server"]
