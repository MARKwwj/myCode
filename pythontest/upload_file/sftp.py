import os

import paramiko


def upload_file(local_path, server_path):
    print(local_path)
    print(server_path)
    # 用户名和密码
    transport = paramiko.Transport(('192.168.100.51', 2222))
    transport.connect(username='resources', password='ZsAbc2020')
    sftp = paramiko.SFTPClient.from_transport(transport)
    # 上传文件
    file_all = os.listdir(local_path)
    sftp.mkdir(server_path)
    for file in file_all:
        cur_file_path = os.path.join(local_path, file)
        server_file_path = server_path + '/' + file
        sftp.put(cur_file_path, server_file_path)

    transport.close()


# upload_file('D:\\desktop\\主题的黑店', '/home/resources/wwwroot/novel/uncrypt/chapterDetails/dasdfjkfhjsdbvj')
# def upload_file_res():
#     # 用户名和密码
#     transport = paramiko.Transport(('192.168.100.51', 2222))
#     transport.connect(username='resources', password='ZsAbc2020')
#     sftp = paramiko.SFTPClient.from_transport(transport)
#     # sftp.get('/wwwroot/novel/crypted/chapterDetails/01d3c310118d2ff7fee6728edc3231cd/1.txt', 'D:\\desktop\\asd\\1.txt')
#     sftp.mkdir('/wwwroot/novel/crypted/chapterDetails/asd')
#     # # 上传文件
#     # file_all = os.listdir(local_path)
#     # sftp.mkdir(server_path)
#     # for file in file_all:
#     #     cur_file_path = os.path.join(local_path, file)
#     #     server_file_path = server_path + '/' + file
#     #     sftp.put(cur_file_path, server_file_path)
#     #
#     # transport.close()
#
#
# upload_file_res()
