# client/Dockerfile
# ベースイメージ
FROM node:14

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係のインストール
COPY package*.json ./
RUN npm install

# アプリのソースコードをコピー
COPY . .

# アプリをビルド
RUN npm run build

# ポートの公開
EXPOSE 3000

# アプリの起動
CMD ["npm", "start"]
