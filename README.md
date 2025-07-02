# self-jvm
JVMを学習するためにGoで実装したJVMです。

javaコンパイル
```
# Amazon Corretto 8を使用してコンパイル
docker pull amazoncorretto:8
docker run -v $(pwd):/MakeJVM -w /MakeJVM amazoncorretto:8 javac MakeJVM.java
```

java実行
```
# Goで実装したJVMを使用して実行
go build -o make_jvm main.go
./make_jvm
```

# 参考
https://techbookfest.org/product/6Jzwfirmy97RcFgCmyjBHC?productVariantID=rrbCNvtUtnk01p5Y25zSCr
https://github.com/murakmii/gojiai/tree/main
https://github.com/zxh0/jvm.go/tree/master
https://github.com/platypusguy/jacobin/wiki/Useful-JVM-links?utm_source=chatgpt.com
https://jacobin.org/?utm_source=chatgpt.com
https://docs.oracle.com/javase/specs/jvms/se8/jvms8.pdf