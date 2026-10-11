# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Встановлюємо сертифікати безпеки
RUN apk --no-cache add ca-certificates tzdata

# Створюємо непривілейованого користувача (UID 10001)
RUN adduser -D -g '' -u 10001 appuser

# Копіюємо залежності для кешування шарів Docker
COPY go.mod ./
RUN go mod download

# Копіюємо вихідний код
COPY . .

# Компіляція статичного бінарника без CGO та з обрізанням налагоджувальної інформації (-s -w)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/api ./cmd/api

# Stage 2: Minimal Runtime Image
FROM alpine:3.21

# Копіюємо сертифікати та часові пояси
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Копіюємо скомпільований бінарник
COPY --from=builder --chown=10001:10001 /app/bin/api /bin/api

# Працюємо від імені непривілейованого користувача
USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["/bin/api"]
