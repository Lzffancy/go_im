Instant messaging server ��ʱͨѶ������
����go channel���,ʹ��tcp����

����
�޸�main.go �е�ip��port
1.����
go build -o im_server user.go main.go server.go
2.��̨����
nohup ./im_server > im_server_out.log 2>&1 &
3.�����־ 
cat im_server.log
[im_server]2023/04/26 15:47:15.105781 main.go:27: -----imserver start run in pid:9843-----
[im_server]2023/04/26 15:47:15.105787 server.go:94: -----imserver run in pid:9843 ok-----

ps -ef |grep "im_server"
root     26186  6451  0 16:33 pts/0    00:00:00 ./im_server

4.�ͻ���ʹ��nc����
nc 127.0.0.1 8888
�������
=============Welcome to chat room============
off: offline your client #�ر�
who: list all client #�����û�
all: send a broadcast message. exp all:hello ==> send 'hello' #�������˷���Ϣ
rename: change your name.exp rename:xxx ===> your name is xxx #�޸��Լ��ǳ�
sendto| send message to someone. exp sendto|username|msg #��ĳ�˷���Ϣ
===========================================
![img.png](img.png)
5.�����ַ

nc 106.52.49.30 8887