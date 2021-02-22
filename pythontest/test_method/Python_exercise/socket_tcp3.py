import socket

# 创建socket 对象
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# 创建连接
s.connect(('127.0.0.1', 9999))
# 接收欢迎信息
print(s.recv(1024).decode('utf-8'))
for data in [b'Michael', b'Tracy', b'sssss']:
    #  发送数据
    s.send(data)
    print(s.recv(1024).decode('utf-8'))
s.send(b'exit')
s.close()
