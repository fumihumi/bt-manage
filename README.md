# bt-manage

macOS で Bluetooth デバイスの **一覧/接続/切断** を行うための CLI（必要に応じて TUI で選択）です。

- 対応OS: macOS
- バックエンド: `blueutil`（外部コマンド）

## インストール

### 前提: blueutil のインストール

`bt-manage` は内部で `blueutil` を呼び出します。

Homebrew を使う場合:

```bash
brew install blueutil
```

### bt-manage のビルド

```bash
go build ./cmd/bt-manage
```

生成された `bt-manage` を PATH の通った場所に置いてください。

## 使い方

### デバイス一覧

```bash
bt-manage list
```

出力形式:

```bash
bt-manage list --format tsv
bt-manage list --format json
bt-manage list --format tsv --no-header
```

### 接続

```bash
bt-manage connect <name-or-prefix>
```

- `<name-or-prefix>` を省略すると、TTY の場合は TUI picker が起動します。
- 複数候補にマッチした場合も、TTY なら picker で選択します。

出力形式:

```bash
bt-manage connect <name-or-prefix> --format tsv
bt-manage connect <name-or-prefix> --format json
bt-manage connect <name-or-prefix> --no-header
```

### 切断

```bash
bt-manage disconnect <name-or-prefix>
```

出力形式:

```bash
bt-manage disconnect <name-or-prefix> --format tsv
bt-manage disconnect <name-or-prefix> --format json
bt-manage disconnect <name-or-prefix> --no-header
```

### インタラクティブ動作の強制/抑制

```bash
bt-manage connect --interactive
bt-manage disconnect --interactive
```

- 非TTY環境で `--interactive` を指定した場合はエラーになります。

### ドライラン

```bash
bt-manage connect <name-or-prefix> --dry-run
bt-manage disconnect <name-or-prefix> --dry-run
```

- 実行はせず、**機械可読な出力のみ**を行います。

### verbose

```bash
bt-manage --verbose list
bt-manage --verbose connect <name-or-prefix>
```

- 内部で実行する `blueutil` コマンドを stderr に出します。

### バージョン

```bash
bt-manage version
```

## 開発

```bash
go test ./...
```

## 既知の制約

- macOS 専用です。
- `blueutil` の出力仕様に依存します（環境によっては出力が異なる可能性があります）。
