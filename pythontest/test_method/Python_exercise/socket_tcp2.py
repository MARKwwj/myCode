import socket

import threading
import time

# 创建一个socket 服务器
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# 绑定 ip 和端口
s.bind(('127.0.0.1', 9999))
# 设置监听，参数是监听数量
s.listen(5)


def tcplink(sock, addr):
    print('Accept new connection from %s:%s...' % addr)
    sock.send(b'welcome!')
    while True:
        data = sock.recv(1024)
        time.sleep(1)
        if not data or data.decode('utf-8') == 'exit':
            break
        sock.send((b'Hello %s ' % data.decode('utf-8').encode('utf-8')))
    sock.close()
    print('Connection from %s:%s closed.' % addr)


print('waiting for connection...')
while True:
    # 接收一个新连接
    sock, addr = s.accept()
    # 创建线程来出来TCP的连接
    t = threading.Thread(target=tcplink, args=(sock, addr))
    # 启动线程
    t.start()
