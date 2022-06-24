# json mask

## 用途

在不直接修改类型源码的前提下，“外挂式”地获取其掩码化的 json 。

## 在什么情况推荐使用json mask

假设有一个类型，我不能直接修改它的源码，只能在同一个包的另一个文件赋予它方法（比如，类型是代码生成的），同时希望在它 json 串行化时进行掩码脱敏。

json mask 可以满足以上需求，因为 json mask 是基于`NeedMask`接口工作，而不是字段 tag 。
只需要为类型实现`NeedMask`接口（可在同一个包的另一个源码文件中），然后通过`JsonMasked`函数就能获取它的掩码化 json 。

如果类型的源码你能够修改，那么[go-mask](https://github.com/ggwhite/go-masker)或许更适合你，它是基于字段 tag 的。

## 意图、设计思路

主要是为了实现[对 protobuf 生成类型的掩码化json](test/personpb)。因为直接修改生成代码不优雅、不安全，所以采用这种“外挂式”的实现。

另外，很大程度上参考了这篇文章：[Golang 打平无结构化的 Json 文本进行脱敏打码](https://0ne.store/2021/11/30/20211130-goland-unstruct-json-mask/)。

## Demo

目前有且只有一种推荐用法，见[TestExample](jsonmask_test.go)。
