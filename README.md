
# tinyWM on Go

tinyWMのGo向け移植です。

### Build

```
go build gotiny.go
```

あとはこのバイナリファイルを、tinyWMと同じような形で実行します。
詳細説明面倒草津。

### 個人的に思ったこと

 - XlibをGoで触るのは難がありまくるかもしれない。特にXEventに関しては、(cgoがCのunion共同体をサポートしていない..ということで)ラッパー関数を無駄につらつらつらつら書かなければ使えないハメになっている。
   - XcbやGtkでは対処できるのだろうか？ 
   - **unsafe.Pointerを使う事で解決できるようである。** *(*XButtonEvent)(unsafe.Pointer(event))といった感じ。
