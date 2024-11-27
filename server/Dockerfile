# server/Dockerfile
# ベースイメージ
FROM golang:1.16

# 作業ディレクトリの設定
WORKDIR /app

# Goモジュールのコピー
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# アプリのソースコードをコピー
COPY . .

# アプリのビルド
RUN go build -o main .

# ポートの公開
EXPOSE 8080

# アプリの起動
CMD ["./main"]
