# Goのイメージを使用
FROM golang:1.16 AS builder

# 作業ディレクトリを設定
WORKDIR /app

# go.modとgo.sumをコピー
COPY go.mod ./
COPY go.sum ./

# 依存関係をダウンロード
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 実行用のイメージを作成
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates

# ビルドしたバイナリをコピー
WORKDIR /root/
COPY --from=builder /app/main .

# アプリケーションを実行
CMD ["./main"]
