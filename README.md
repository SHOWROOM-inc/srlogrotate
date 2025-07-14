# srlogrotate

Go言語でSRサーバ向けのログローテーションを実現するライブラリです。

## 使い方

`io.Writer`インターフェイスで提供しているので、一般のログライブラリの出力先指定にそのまま指定できます。

```go
import (
"github.com/SHOWROOM-inc/srlogrotate"
"go.uber.org/zap"
"go.uber.org/zap/zapcore"
)

func main() {
    infoLogFile := srlogrotate.NewLogger("/var/log/app/default.info.log")
    core := zapcore.NewCore(encoder, zapcore.AddSync(infoLogFile), zapcore.InfoLevel),
    zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
```

## 仕様

- 与えられたファイル名の語尾に `filename.yyyymmdd`を付けたファイル名で出力します。
- ログファイルの行数やサイズで分割は行いません。
- 現時点では、ローテートしたファイルについて圧縮や削除は行いません。
