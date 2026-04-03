# ADR-001: WebSocket over SSE for Real-time Communication

## Status
Accepted

## Context
かるた対戦ではリアルタイムの双方向通信が必要。選択肢は WebSocket と Server-Sent Events (SSE)。

## Decision
WebSocket を採用する。

## Rationale
- **双方向性**: SSEはサーバー→クライアントの一方向のみ。かるたの取り札アクションはクライアント→サーバーの低遅延送信が必須。SSEだとPOSTリクエストを別途送る必要があり、遅延が増える。
- **レイテンシ**: WebSocketは1つの接続でフレーム単位の送受信。SSE+POSTは2つのHTTP接続が必要で、取り札の先着判定（ミリ秒単位）に不利。
- **状態管理**: WebSocket接続はサーバー側で状態を保持しやすく、切断/再接続の管理が容易。
- **採用実績**: gorilla/websocket はGoの標準的なWebSocketライブラリで、十分な実績がある。

## Consequences
- WebSocket非対応のプロキシ/ファイアウォール環境では接続できない可能性がある（PWAターゲットのため許容範囲）
- サーバー側の接続管理コストが発生する
