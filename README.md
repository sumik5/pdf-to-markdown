# pdf-to-markdown

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/shivase/pdf-to-markdown)](https://goreportcard.com/report/github.com/shivase/pdf-to-markdown)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

PDFファイルをMarkdownに変換するGo製CLIツール。テキスト抽出・見出し検出・リスト検出を行い、構造化されたMarkdownを出力する。

## 機能

- テキスト抽出（座標・フォントサイズ・Bold属性付き）
- 見出し検出（フォントサイズとBold属性によるh1〜h3判定）
- リスト検出（箇条書き・番号付き）
- インデントレベル検出
- メタデータ抽出（Title・Author・Subject・Keywords等）
- アウトライン（目次）抽出
- ページ区切り付きMarkdown出力
- 変換統計をstderrにJSON出力（ページ数・文字数）

## インストール

### 前提条件

- [Go 1.26以上](https://go.dev/dl/)
- [mise](https://mise.jdx.dev/)（タスク管理用、オプション）

### ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/shivase/pdf-to-markdown.git
cd pdf-to-markdown

# miseを使う場合（推奨）
mise run build
sudo cp build/pdf-to-markdown /usr/local/bin/

# またはgoコマンドで直接ビルド
go build -o pdf-to-markdown .
sudo mv pdf-to-markdown /usr/local/bin/
```

## クイックスタート

```bash
pdf-to-markdown input.pdf output.md
```

変換が完了すると、`output.md` が生成される。変換統計はstderrにJSON形式で出力される。

```json
{"pages":10,"characters":5432,"outputFile":"output.md"}
```

## CLIリファレンス

| 書式 | 説明 |
|------|------|
| `pdf-to-markdown <input.pdf> <output.md>` | PDFをMarkdownに変換する |

| 引数 | 説明 |
|------|------|
| `<input.pdf>` | 変換元のPDFファイルパス |
| `<output.md>` | 出力先のMarkdownファイルパス |

**stderr出力（JSON）**

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `pages` | number | 処理したページ数 |
| `characters` | number | 出力Markdownの文字数 |
| `outputFile` | string | 出力ファイルパス |

## 開発者向け

### セットアップ

```bash
# miseをインストール（まだの場合）
curl https://mise.jdx.dev/install.sh | sh

# miseの設定を信頼
mise trust

# 利用可能なタスクを確認
mise run help
```

### ビルド

```bash
# 本番ビルド
mise run build

# 開発ビルド（race detector付き）
mise run dev
```

### テスト

```bash
# テスト実行
mise run test

# カバレッジレポート生成（coverage.html）
mise run test-coverage
```

### Lintとフォーマット

```bash
# Lintチェック（golangci-lint）
mise run lint

# コードフォーマット
mise run fmt

# すべてのチェックを一括実行
mise run check-all
```

### miseタスク一覧

| タスク | 説明 |
|--------|------|
| `mise run build` | バイナリをビルド（`build/pdf-to-markdown`） |
| `mise run dev` | race detector付きでビルド |
| `mise run test` | テスト実行 |
| `mise run test-coverage` | カバレッジレポート生成 |
| `mise run lint` | golangci-lintでLintチェック |
| `mise run fmt` | `go fmt`でコードフォーマット |
| `mise run check-all` | test・lint・fmtを一括実行 |
| `mise run install` | `/usr/local/bin`にインストール |
| `mise run clean` | ビルド成果物を削除 |
| `mise run mod-update` | Goモジュールを更新 |
| `mise run run` | `PDF_INPUT`・`PDF_OUTPUT`環境変数で実行 |

### プロジェクト構造

```
pdf-to-markdown/
├── main.go                          # CLIエントリポイント
├── main_test.go                     # メインテスト
├── go.mod                           # Goモジュール定義
├── go.sum                           # 依存関係チェックサム
├── .mise.toml                       # miseタスク設定
├── .golangci.yaml                   # golangci-lint設定
└── internal/
    ├── model/
    │   └── types.go                 # 共通型定義
    ├── extractor/
    │   ├── extractor.go             # PDF抽出層
    │   ├── metadata.go              # メタデータ抽出
    │   └── outline.go               # アウトライン抽出
    ├── converter/
    │   ├── converter.go             # ページ変換処理
    │   ├── grouper.go               # 行グループ化
    │   └── detector.go              # 見出し・リスト検出
    └── markdown/
        └── builder.go               # Markdown構築
```

### 依存ライブラリ

| ライブラリ | バージョン | ライセンス |
|-----------|-----------|-----------|
| [ledongthuc/pdf](https://github.com/ledongthuc/pdf) | v0.0.0-20250511090121 | BSD-3-Clause |

Pure GoによるPDF解析ライブラリを使用しているため、CGOは不要。

## ライセンス

このプロジェクトは[MIT License](LICENSE)の下で公開されています。

## 作者

[shivase](https://github.com/shivase)
