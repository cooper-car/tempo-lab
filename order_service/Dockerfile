# 使用 Go 的官方映像
FROM golang:1.22 AS builder

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源碼
COPY . .

# 編譯應用程式
RUN GOOS=linux GOARCH=amd64 go build -o order_service . && chmod +x order_service

# 使用輕量級的映像
FROM alpine:latest

# 複製編譯好的二進制檔案
COPY --from=builder /app/order_service .

# 設定睡5分鐘
#CMD ["sleep", "300"]
CMD ["./order_service"]