# 1. Aşama: Build
FROM golang:1.22.3 AS builder

WORKDIR /app

# go.mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./
RUN go mod download

# Projenin tüm dosyalarını kopyala
COPY . .

# cmd/api altındaki main.go'yu derle, çıktıyı 'main' olarak adlandır
#RUN go build -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api


# 2. Aşama: Küçük çalışma imajı
#FROM debian:stable-slim
FROM alpine:latest

WORKDIR /app
# HTTPS istekleri için gerekli ?
RUN apk --no-cache add ca-certificates

# Derlenmiş binary'yi kopyala
COPY --from=builder /app/main .

# Eğer uygulama environment dosyasını kullanıyorsa ve container içine dahil etmek istiyorsan:
COPY .env .

# Uygulamanın dinlediği port
EXPOSE 5454

# Container başlatıldığında çalışacak komut
CMD ["./main"]

#docker build -t pwp-backend .
#docker run --rm -p 5454:5454 pwp-backend
#docker run -d --name pwp-api -p 5454:5454 pwp-backend

