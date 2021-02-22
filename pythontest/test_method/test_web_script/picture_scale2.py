import multiprocessing
import os
import re
import time

import numpy as np
from PIL import Image
import pymongo
import json


def cal_time(func):
    def wrapper(*args):
        start_time = time.time()
        res = func(*args)
        all_time = round(time.time() - start_time, 1)
        print("func {0} run time {1} s".format(func.__name__, all_time))
        return res

    return wrapper


def GetScale(file_path):
    """
    计算图片比例方法
    :param file_path: 文件路径
    :return: 图片比例
    """
    scale = 0
    try:
        img = Image.open(file_path)
        img_width, img_height = img.width, img.height
        scale = round(img_width / img_height, 3)
    except Exception as e:
        print(e)
        scale = 0

    return scale  # 返回图片比例


@cal_time
def encrypt(picRead, pic_name):
    """
    资源 加密 /解密 方法
    :param picRead:二进制文件 （open("rb")）
    :return: res
    """
    # 读取二进制文件
    picData = picRead.read()
    picLen = len(picData)
    # 需要写出的文件路径
    cur_path = os.getcwd() + "/" + pic_name
    f_write = open(cur_path, "wb", 1024 * 1024 * 5)
    # 获取Key
    data = [89, 61, 148, 105, 40, 237, 255, 227, 195, 111, 172, 197, 236, 82, 43, 163,
            121, 78, 245, 50, 67, 117, 136, 3, 26, 132, 52, 204, 182, 83, 13, 146,
            21, 215, 47, 48, 189, 96, 181, 23, 1, 158, 219, 184, 86, 112, 51, 84,
            46, 59, 191, 106, 140, 4, 65, 173, 109, 248, 88, 53, 152, 153, 36, 115,
            37, 242, 177, 90, 179, 196, 143, 217, 239, 251, 69, 193, 55, 42, 147, 76,
            134, 218, 9, 174, 141, 216, 138, 129, 124, 68, 107, 234, 249, 102, 64, 169,
            185, 122, 56, 229, 41, 252, 127, 18, 79, 205, 186, 199, 108, 214, 160, 16,
            226, 94, 240, 209, 253, 187, 166, 99, 5, 213, 34, 155, 154, 0, 211, 97,
            72, 175, 238, 167, 70, 119, 31, 113, 222, 2, 66, 156, 162, 194, 207, 206,
            157, 100, 246, 202, 171, 20, 54, 11, 208, 220, 62, 126, 14, 114, 95, 32,
            73, 165, 254, 35, 244, 81, 149, 137, 135, 176, 93, 15, 44, 57, 164, 92,
            241, 19, 203, 87, 6, 243, 235, 151, 12, 24, 183, 33, 188, 144, 224, 178,
            150, 27, 39, 233, 116, 25, 103, 110, 201, 85, 192, 159, 232, 190, 212, 123,
            131, 22, 210, 45, 74, 10, 125, 180, 130, 77, 221, 133, 60, 230, 80, 75,
            198, 128, 247, 29, 231, 118, 63, 250, 225, 120, 168, 104, 161, 28, 145, 223,
            170, 58, 38, 8, 142, 98, 71, 91, 30, 101, 200, 7, 17, 139, 49, 228]
    i = 0
    # 进行循环加密
    for j in range(picLen):
        newByte = data[i] ^ picData[j]
        i = (i + 1) % len(data)
        f_write.write(bytes([newByte]))
    f_write.close()
    picRead.close()
    return cur_path


def ConnMongo():
    """
    连接数据库
    :return: mycol
    """
    # 测试
    # myclient = pymongo.MongoClient("mongodb://80.251.211.152:17027/")
    # 正式
    # myclient = pymongo.MongoClient("mongodb://80.251.223.182:38028/")
    myclient = pymongo.MongoClient("mongodb://64.64.241.63:7777/")
    mydb = myclient["test"]
    mycol = mydb["cartoon"]
    return mycol


def WriteTxt(id, mycol, txt_name):
    """
    写入txt
    :param id:
    :param mycol:
    :return: is_alter
    """
    is_alter = False
    info_path = os.getcwd() + '/' + txt_name
    info = open(info_path, 'w', )
    for x in mycol.find({"_id": id}, {"_id": 0, "chapters": [{"chapterUrl": 1}]}):
        info.write((json.dumps(x, indent=2, ensure_ascii=False)).encode("gbk", 'ignore').decode("gbk", "ignore"))
        if re.search('jpg_', str(x)) is not None:
            is_alter = True
            print("此漫画已经被修改过！！！ id: {0} ".format(id))
            break
    info.close()
    return is_alter


def UpdTxt(upd_path_res_old_list, upd_path_res_list, txt_name):
    """
    替换txt 内容
    :param upd_path_res_old_list: 旧的值
    :param upd_path_res_list: 新的值
    :return:
    """
    info_path = os.getcwd() + '/' + txt_name
    file_data = ''
    with open(info_path, 'r') as file:
        for line in file:
            for old_path in upd_path_res_old_list:
                cur_index = upd_path_res_old_list.index(old_path)
                if old_path in line:
                    line = line.replace(old_path, upd_path_res_list[cur_index])
            file_data += line
    with open(info_path, 'w') as file:
        file.write(file_data)


def Mongogetpath(id, mycol):
    """
    获取mongodb漫画信息
    :param id:
    :param mycol:
    :return:
    """
    path = mycol.find_one({"_id": id}, {"_id": 0, 'picUrl': 1})
    remote_path = "/home/resources/wwwroot/limao/crypted/" + (path['picUrl'].split('cover.jpg', 1))[0]
    chapter = (mycol.find_one({"_id": id}, {"_id": 0, 'totalNum': 1}))['totalNum']
    # print(remote_path)
    print("漫画id: {0} 漫画章节： {1} ".format(id, chapter))
    # /home/resources/wwwroot/limao/crypted/cartoonInfo/964cb5df130e327f708cb25dafc523a6/    115
    return remote_path, chapter


def UpdMongo(id, mycol, txt_name):
    """
    修改数据库 漫画图片信息
    :param mycol:
    :return:
    """
    info_path = os.getcwd() + '/' + txt_name
    with open(info_path, 'r') as file:
        file_data = file.readlines()
    condition = {'_id': id}
    total_num_change = eval(((str(file_data).replace("'", "")).replace("\\n,", "")).strip("[]"))
    # print(total_num_change)
    result = mycol.update_one(condition, {'$set': total_num_change})
    print("修改了：{}".format(result.modified_count))


def SelMongo(mycol):
    """
    查询所有漫画id
    :param mycol:
    :return:
    """
    id_list = []
    result = mycol.find({}, {'_id': 1})
    for i in result:
        id_list.append(i['_id'])
    # print("漫画id:{}".format(id_list))
    return id_list


@cal_time
def PicScaUpd(remote_path, chapter, pic_name, thread_name):
    """
    拼接漫画比例
    :param remote_path:
    :param chapter:
    :param pic_name:
    :return:
    """
    try:
        upd_path_res_old_list = []
        upd_path_res_list = []

        for i in range(1, chapter + 1):
            cur_remote_path = str(remote_path + str(i) + '/')
            print(thread_name + " - " + cur_remote_path)
            # noinspection PyBroadException
            try:
                remote_files = os.listdir(cur_remote_path)
                remote_files.sort(key=lambda x: int(x[:-4]))
                print('{0} 当前目录下所有文件：{1}'.format(thread_name, remote_files))
            except Exception as e:
                print(e)
                break
            for file in remote_files:  # 遍历读取远程目录里的所有文件
                picRead = open(cur_remote_path + file, 'rb', 1024 * 1024 * 10)
                decode_pic_path = encrypt(picRead, pic_name)  # 解密文件
                scale = GetScale(decode_pic_path)  # 读取文件比例
                upd_path_res_old = ((((cur_remote_path.split("/home/resources/", 1))[1] + file).split('crypted', 1))[1].split('/', 1))[1]
                upd_path_res = (((((cur_remote_path.split("/home/resources", 1))[1] + file).split('crypted', 1))[1] + "_" + str(scale)).split('/', 1))[1]
                print(thread_name + " - " + upd_path_res)
                upd_path_res_old_list.append(upd_path_res_old)
                upd_path_res_list.append(upd_path_res)
        return upd_path_res_old_list, upd_path_res_list
    except Exception as e:
        print(e)


def run(txt_name, pic_name, id_list, thread_name):
    print("%s is beginning !" % thread_name)
    mycol = ConnMongo()
    for id in id_list:
        is_alter = WriteTxt(id, mycol, txt_name)
        if is_alter:
            continue
        remote_path, chapter = Mongogetpath(id, mycol)
        upd_old_list, upd_new_list = PicScaUpd(remote_path, chapter, pic_name, thread_name)
        UpdTxt(upd_old_list, upd_new_list, txt_name)
        UpdMongo(id, mycol, txt_name)
    print("%s is ending !" % thread_name)


if __name__ == '__main__':
    # 直接解密读取服务本地文件
    # 获取所有漫画id
    mycol = ConnMongo()
    id_list = SelMongo(mycol)
    # 进程 1
    process1 = multiprocessing.Process(target=run, args=('info1.txt', 'pic1.jpg', id_list[770:820], "thread1"))
    # 进程 2
    process2 = multiprocessing.Process(target=run, args=('info2.txt', 'pic2.jpg', id_list[820:880], "thread2"))
    # # 进程 3
    # process3 = multiprocessing.Process(target=run, args=('info3.txt', 'pic3.jpg', id_list[299:200:-1], "thread3"))
    # # 进程 4
    # process4 = multiprocessing.Process(target=run, args=('info4.txt', 'pic4.jpg', id_list[399:300:-1], "thread4"))
    # # 进程 5
    # process5 = multiprocessing.Process(target=run, args=('info5.txt', 'pic5.jpg', id_list[499:400:-1], "thread5"))
    # # 进程 6
    # process6 = multiprocessing.Process(target=run, args=('info6.txt', 'pic6.jpg', id_list[599:500:-1], "thread6"))
    # # 进程 7
    # process7 = multiprocessing.Process(target=run, args=('info7.txt', 'pic7.jpg', id_list[699:600:-1], "thread7"))
    # # 进程 8
    # process8 = multiprocessing.Process(target=run, args=('info8.txt', 'pic8.jpg', id_list[799:700:-1], "thread8"))

    process1.start()
    process2.start()
    # process3.start()
    # process4.start()
    # process5.start()
    # process6.start()
    # process7.start()
    # process8.start()
