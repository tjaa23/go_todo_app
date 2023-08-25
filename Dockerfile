## deploy-builder:リリース用のビルドを行うコンテナイメージを作成するステージ
# デプロイ用にコンテナに含めるバイナリを作成するコンテナ
FROM golang:1.21.0-bullseye AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY sample .
RUN go build -trimpath -ldflags "-w -s" -o app

#----------------------------------------------

# deploy:マネージドサービス上で動かすことを想定したリリース用のコンテナイメージを作成するステージ
# デプロイ用のコンテナ
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

#----------------------------------------------

# ローカルで開発するときに利用するコンテナイメージを作成するステージ
# ローカル開発環境で利用するホットリロード環境
FROM golang:1.21.0 as dev
WORKDIR /app
# ファイルが更新されたことを検知して実行中のgoプログラムを再起動してくれる
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", ".air.toml"]