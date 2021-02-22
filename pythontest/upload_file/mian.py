import hashlib
import os
import time

from upload_file.mongoDB import cal_id, insert_novel, insert_chapter_detail_mdb
from upload_file.resource_encrypt import show_files
from upload_file.sftp import upload_file


def _upload_eFile(file_path):
    file_path_list = str(file_path).split('\\')
    file_name = file_path_list[len(file_path_list) - 1]
    new_file_name = (hashlib.md5((file_name + str(round((time.time() * 1000), 1))).encode('utf-8'))).hexdigest()
    new_file_path = os.path.join(os.path.dirname(file_path) + '\\' + new_file_name)
    print(new_file_path)
    print(new_file_name)
    # 解密本地文件 传入空的list接收文件名
    # file_path 加密前的文件
    # new_file_path 加密后的文件
    show_files(file_path, new_file_path, [])

    """ 上传小说到服务器 """
    # 上传加密的到服务器
    upload_file(new_file_path, '/wwwroot/novel/uncrypt/chapterDetails/{0}'.format(new_file_name))
    upload_file(new_file_path, '/wwwroot/novel/crypted/chapterDetails/{0}'.format(new_file_name))
    """ 存入数据库 """
    novel_id = cal_id()
    insert_novel(file_name, novel_id)
    insert_chapter_detail_mdb(file_path, new_file_name, novel_id)


_upload_eFile('D:\\desktop\\梦回大唐')
