# caty

用户认证

## 监听端口

8120

## go-swagger本地简单运用

下载：go get -u github.com/go-swagger/go-swagger/cmd/swagger  
生成文档：swagger generate spec -o ./docs/swagger.json  
文档解引用：swagger flatten --with-expand ./docs/swagger.json -o ./docs/swagger.json  
本地简单运用：swagger serve -F=swagger ./docs/swagger.json  
文档编写规则请参考：  
官方文档：https://goswagger.io/  
其他文档：https://zhuanlan.zhihu.com/p/136521497

## golangci-lint使用

下载：go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1  
使用：golangci-lint run -c ./build/golangci-lint-demo.yml --tests=false --out-format=json  > golangci-lint.json 2>&1  
make命令：make lint