# Étape de construction (Build)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Télécharger les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le code source
COPY . .

# Compiler l'application (Binaire statique léger)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Étape finale (Image minimale pour l'exécution)
FROM alpine:latest  

WORKDIR /root/

# Copier le binaire depuis l'étape de construction
COPY --from=builder /app/main .

# Exposer le port
EXPOSE 8080

# Lancer l'app
CMD ["./main"]