import socket

# 创建UDP 服务端 SOCK_DGRAM 指定了这个socket 类型是 udp
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
# 绑定端口 不需要 listen() 监听 直接接收来自任何客户端的数据
s.bind(('127.0.0.1', 9999))
print('Bind udp on 9999...')
while True:
    # 接收数据
    data, addr = s.recvfrom(1024)
    print("Received from %s:%s." % addr)
    s.sendto()
