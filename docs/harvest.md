# Harvest Report: Astro-Karuta v1.0.0

**Date**: 2026-04-04
**Tag**: v1.0.0
**Repository**: akaitigo/astro-karuta

---

## Project Summary

Astro-Karuta is a real-time multiplayer educational game that teaches astronomy through the traditional Japanese card game "karuta". Players compete to grab the correct picture card when a reading clue is announced via WebSocket, collecting celestial objects (constellations, planets, phenomena) along the way.

### Tech Stack

| Layer     | Technology                             |
|-----------|----------------------------------------|
| Frontend  | TypeScript, Next.js 15, React 19, PWA |
| Real-time | WebSocket (gorilla/websocket)          |
| Backend   | Go 1.23, net/http stdlib router        |
| Testing   | Vitest + Testing Library, go test      |
| CI        | GitHub Actions                         |

---

## Statistics

### Issues

| Metric     | Count |
|------------|-------|
| Total      | 5     |
| Closed     | 5     |
| Open       | 0     |

### Pull Requests

| Metric | Count |
|--------|-------|
| Total  | 10    |
| Merged | 10    |
| Open   | 0     |

| PR # | Title | Type |
|------|-------|------|
| #6   | Card/Deck CRUD API with seed data | feat |
| #7   | WebSocket battle engine with matchmaking | feat |
| #8   | Collection feature (backend + frontend) | feat |
| #9   | Seasonal deck and observation missions | feat |
| #10  | Game frontend UI | feat |
| #11  | README and CHANGELOG for v1.0.0 release | docs |
| #12  | Review round 1 - security and robustness | fix |
| #13  | Review round 2 - CRITICAL/HIGH security and robustness fixes | fix |
| #14  | Review round 4 - errcheck, UUID user ID, promise leak | fix |
| #15  | Review round 5 - game timeout, reconnect state, graceful shutdown | fix |

**Merge rate**: 10/10 (100%)

### Tests

| Scope    | Tests |
|----------|-------|
| Backend  | 76+ (race detector enabled) |
| Frontend | 65    |

---

## Review Loop History

### Round 1 (セキュリティ基本)
**PR #12** | 検出: CRITICAL 2, HIGH 5 → 全解消

| 重要度 | 指摘 | 修正 |
|--------|------|------|
| CRITICAL | WebSocket CheckOrigin が全オリジン許可 | CORS_ORIGIN 環境変数で検証 |
| CRITICAL | config.go にハードコード DB credential | 空デフォルトに変更 |
| HIGH | createGameState nil 戻り値未チェック | nil ガード追加 |
| HIGH | MarshalMessage エラー無視 10箇所 | 全てログ付きハンドリング |
| HIGH | 不正解 grab でカードが進行するバグ | 正解時のみ進行に修正 |
| HIGH | WebSocket 接続数無制限 | 200接続上限追加 |
| HIGH | WS URL パス不一致 | /api/v1/ws に統一 |

### Round 2 (セキュリティ深掘り+入力検証+型統一)
**PR #13** | 検出: CRITICAL 5, HIGH 10 → 全解消

3つの専門エージェント（コードレビュー、セキュリティ監査、破壊的検証）を並列で実施。

| 重要度 | 指摘 | 修正 |
|--------|------|------|
| CRITICAL | 認証メカニズム完全欠如 | ADR-003で猶予判断、user_id UUID検証追加 |
| CRITICAL | PlayerName/RoomCode バリデーション不在 | 空文字・長さ・制御文字チェック追加 |
| CRITICAL | HandleGrab にゲーム状態チェックなし | playing 以外は拒否 |
| CRITICAL | config.Load() 完全未使用 | main.go で一本化 |
| CRITICAL | WebSocket Origin 空許可 | 本番では Origin 必須に |
| HIGH | reconnect で他人の player_id 乗っ取り | clientID 不一致をログ警告 |
| HIGH | 不正解 grab 回数無制限（総当たり） | 1ラウンド1回制限 (HasGrabbed) |
| HIGH | REST API に CORS ミドルウェア未実装 | middleware/cors.go 追加 |
| HIGH | FE/BE 型定義不一致 | card_grabbed 削除、bonus_card 統一、game_state ハンドラ追加 |
| HIGH | Makefile migrate 未実装 | 削除 |
| HIGH | CompleteMission lat=0/lng=0 未検証 | 明示的拒否追加 |
| HIGH | time.Local 使用→本番 TZ ズレ | time.UTC に統一 |
| HIGH | DB スキーマ user_missions 欠落 | マイグレーション追加 |
| HIGH | ミッション→コレクション未連携 | CompleteMission で AddToCollection 呼び出し |
| HIGH | WebSocket setTimeout 接続待ち | waitForConnection Promise に置換 |

### Round 3-4 (CI+機能バグ+Promiseリーク)
**PR #14** | 検出: HIGH 3, MEDIUM 2 → 全解消

| 重要度 | 指摘 | 修正 |
|--------|------|------|
| HIGH | golangci-lint errcheck 3箇所で CI 失敗 | fmt.Fprintf, Conn.Close() の戻り値チェック |
| HIGH | DEFAULT_USER_ID="user-1" が非UUID→400エラー | 有効な UUID に変更 |
| HIGH | lobby の as 型アサーション不統一 | — (in演算子チェック付きで許容) |
| MEDIUM | waitForConnection Promise リーク | 再接続不可時に全 reject |
| MEDIUM | lobby/page.tsx as アサーション | 実質的な型安全性は担保済み |

### Round 5 (ゲームロジック+運用品質)
**PR #15** | 検出: HIGH 6, MEDIUM 8 → HIGH 全解消

新しい視点で深掘り: ゲームロジック正確性、FE堅牢性、データ整合性、運用品質、テスト品質

| 重要度 | 指摘 | 修正 |
|--------|------|------|
| HIGH | ゲームタイムアウト未発動（PRD の5分制限が機能しない） | time.AfterFunc + endGame 時 Stop |
| HIGH | reconnect 後に CurrentCard 未送信 | statePayload に reading_text, candidates 追加 |
| HIGH | グレースフルシャットダウンで goroutine 残留 | Hub.Shutdown() 実装 |
| HIGH | 待機キューメモリリーク（切断クライアント永久残留） | HandleDisconnect で除去 |
| HIGH | 季節デッキ並列生成競合 | sync.Mutex + ダブルチェック |
| HIGH | math/rand シード未設定の懸念 | Go 1.23 デフォルト動作をコメント明記 |

### Round 6 (最終確認)
**検証のみ** | CRITICAL=0, HIGH=0 → **PASS**

残存 MEDIUM 4件（全て許容判断済み）:
- joinExistingGame TOCTOU — endGame ガードで整合性保持
- lobby as 型アサーション — in 演算子チェック付き
- 0枚カードゲーム開始 — LoadCards 失敗時は Fatalf
- distractor 不足 — シードデータ4枚以上保証

### 累計
| 重要度 | 検出総数 | 修正数 | 残存 |
|--------|----------|--------|------|
| CRITICAL | 7 | 7 | 0 |
| HIGH | 24 | 24 | 0 |
| MEDIUM | 12+ | 8+ | 4 (許容) |

---

## Lessons Learned

### 1. レビューは1ラウンドでは終わらない
Round 1 で「CRITICAL=0, HIGH=0」と判定して終了したが、Round 2 で別視点（セキュリティ監査、破壊的検証）を入れたら CRITICAL 5 + HIGH 10 が出た。**1つの観点で問題ゼロでも、別の観点から見ると重大な問題が残っている**。最低3ラウンド・3観点（コード品質、セキュリティ、破壊的検証）は必須。

### 2. CIが通っている≠品質が担保されている
errcheck が Round 1 からずっと失敗していたが、テストは全パスだったためスルーされていた。**lint も含めた `make check` 全通過を修正完了の条件にすべき**。

### 3. ゲームロジックは「仕様 vs 実装」のギャップが埋もれやすい
PRD に「5分制限」と書いたが、タイムアウトの実装が完全に抜けていた。reconnect 後のカード状態送信も同様。**PRD の受け入れ条件を1つずつテストで検証する工程**が必要。

### 4. 並列エージェントで多面的レビューは有効
Round 2 で コードレビュー / セキュリティ監査 / 破壊的検証 の3エージェントを並列実行し、合計38件の指摘を一度に洗い出せた。後出しがなくなる。

### 5. FE/BE の型定義ドリフトは必ず起きる
TypeScript の型とGoのモデルが手動同期のため、`card_grabbed` (BE側に存在しない)、`bonus_card` のフィールド不足、`game_state` ハンドラ未実装など複数のズレが発生。**スキーマファースト（OpenAPI or protobuf）で生成すべき**。

### 6. 「開発時の仮実装」は本番まで残る
`DEFAULT_USER_ID = "user-1"`、`CheckOrigin: return true`、`config.Load()` の未使用、`time.Local` の使用 — 全て「あとで直す」つもりの仮実装がそのまま残った。**仮実装には `// FIXME:` ではなく、lint で検出できる仕組み**が必要。

### 7. メモリリークはレビューしないと見つからない
待機キューの切断クライアント残留、Promise の永久 pending — テストでは再現しにくく、運用で初めて顕在化するタイプ。**リソース解放のレビュー観点を明示的にチェックリストに入れるべき**。

---

## Template / Skill Improvement Suggestions

### pipeline-*.md テンプレートへの反映

1. **Review Loop の最低ラウンド数を3に引き上げ**
   - Round 1: コード品質（デッドコード、命名、テスト品質）
   - Round 2: セキュリティ + 破壊的検証（3エージェント並列）
   - Round 3: ゲームロジック / ビジネスロジック正確性 + 運用品質
   - 各ラウンドで観点を変え、同じ視点の再レビューを避ける

2. **CI全通過を修正完了の必須条件に明記**
   - `make check`（lint含む）が通らない状態でレビュー完了と判定しない
   - golangci-lint errcheck は初期スキャフォールドから有効化

3. **PRD 受け入れ条件の逆引きテスト工程を追加**
   - Ship 前に PRD の各受け入れ条件に対応するテストが存在するか確認
   - 「5分制限」のようなテストされていない仕様を発見するため

4. **仮実装検出の仕組み化**
   - `grep -rn "TODO\|FIXME\|user-1\|return true\|time.Local" ` を CI に組み込み
   - `CheckOrigin.*return true` パターンを golangci-lint カスタムルールで検出

5. **FE/BE 型同期のガードレール**
   - TypeScript の型定義と Go のモデルの対応表を CLAUDE.md に記載
   - または OpenAPI / protobuf からの型生成をスキャフォールドに組み込み

6. **レビューエージェントの指示テンプレート改善**
   - 「後出し禁止」だけでなく「前ラウンドで見たファイルも再読必須」を明記
   - リソース解放（goroutine leak, channel close, Promise reject）を明示的チェック項目に

7. **Harvest にレビューラウンド詳細を記録する指示を強化**
   - 「Round 1 で問題ゼロ → 終了」としないためのガードレール
   - 各ラウンドの検出数・修正数・残存数の構造化テーブルを必須化
