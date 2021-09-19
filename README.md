# obs

对象服务

## 监听端口

8120

## go-swagger本地简单运用

下载：go get -u github.com/go-swagger/go-swagger/cmd/swagger  
生成文档：swagger generate spec -o ./docs/swagger.json  
文档解引用：swagger flatten --with-expand ./docs/swagger.json -o ./docs/swagger.json  
本地简单运用：swagger serve -F=swagger ./docs/swaggerResult.json  
文档编写规则请参考：  
官方文档：https://goswagger.io/  
其他文档：https://zhuanlan.zhihu.com/p/136521497