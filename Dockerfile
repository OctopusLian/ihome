FROM golang:latest
#创建工作目录
RUN mkdir -p /go/src/ihome
#进入工作目录
WORKDIR /go/src/ihome
#将当前目录下的所有文件复制到指定位置
COPY . /go/src/ihome
#下载beego和bee

RUN go env -w GOPROXY=https://goproxy.cn
RUN go get github.com/astaxie/beego && go get github.com/beego/bee && go get github.com/go-sql-driver/mysql && go get -u github.com/beego/beego/v2/server/web/session/redis && go get github.com/prometheus/client_golang/prometheus@v1.7.0
#端口
EXPOSE 8080
#运行
CMD ["bee", "run"]