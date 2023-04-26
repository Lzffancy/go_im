Instant messaging server 即时通讯服务器
基于go channel设计,使用tcp连接

部署
修改main.go 中的ip和port
1.编译
go build -o im_server user.go main.go server.go
2.后台运行
nohup ./im_server > im_server_out.log 2>&1 &
3.检查日志 
cat im_server.log
[im_server]2023/04/26 15:47:15.105781 main.go:27: -----imserver start run in pid:9843-----
[im_server]2023/04/26 15:47:15.105787 server.go:94: -----imserver run in pid:9843 ok-----

ps -ef |grep "im_server"
root     26186  6451  0 16:33 pts/0    00:00:00 ./im_server

4.客户端使用nc命令
nc 127.0.0.1 8888
进入界面
=============Welcome to chat room============
off: offline your client #关闭
who: list all client #在线用户
all: send a broadcast message. exp all:hello ==> send 'hello' #给所有人发消息
rename: change your name.exp rename:xxx ===> your name is xxx #修改自己昵称
sendto| send message to someone. exp sendto|username|msg #给某人发消息
===========================================

5.体验地址

nc 106.52.49.30 8887