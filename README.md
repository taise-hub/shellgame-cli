# シェルゲ-
シェルゲーは対戦型コマンドラインクイズゲームです。
プレイヤーごとに独立したシェルゲーコンテナを用意しており、WebSocketを用いてシェルゲーコンテナに接続します。

## 🚧DEMO🚧
![デモ](./static/demo.mov)
 
## 🚧Usage🚧
シェルゲーサーバを立てる
```
$ git clone https://github.com/taise-hub/shellgame-cli
$ cd server
$ go run cmd/shellgame/main.go
```
 
シェルゲークライアントを実行する
```bash
$ git clone https://github.com/taise-hub/shellgame-cli
$ cd client
$ go run cmd/shellgame/main.go
```