# Build stage
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Instalar dependências
RUN apk add --no-cache git

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /app/server ./movies-service/cmd

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar executável do builder
COPY --from=builder /app/server .
COPY --from=builder /app/movies-service/movies.json ./movies-service/movies.json

# Expor porta do gRPC
EXPOSE 50051

# Variáveis de ambiente
ENV MONGO_URI=mongodb://mongo:27017

CMD ["./server"]
