import io
import os
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
    # file_path = 'D:\\desktop\\1.jpg'
    img = Image.open(file_path)
    img_width, img_height = img.width, img.height
    scale = round(img_width / img_height, 3)
    return scale  # 返回图片比例


def show_files(path, path2, all_files):
    # 首先遍历当前目录所有文件及文件夹
    file_list = os.listdir(path)
    # 准备循环判断每个元素是否是文件夹还是文件，是文件的话，把名称传入list，是文件夹的话，递归
    for file in file_list:
        # 利用os.path.join()方法取得路径全名，并存入cur_path变量，否则每次只能遍历一层目录
        cur_path = os.path.join(path, file)
        # 判断是否是文件夹
        if os.path.isdir(cur_path):
            cur_path2 = os.path.join(path2, file)
            isExists = os.path.exists(cur_path2)
            if not isExists:
                os.makedirs(cur_path2)
                print(cur_path2 + ' 创建成功')
            else:
                # 如果目录存在则不创建，并提示目录已存在
                print(cur_path2 + ' 目录已存在')
            show_files(cur_path, cur_path2, all_files)
        else:
            encrypt(path, path2, file)
            print(file)
            all_files.append(file)

    return all_files


@cal_time
def encrypt(filePath, path2, name):
    # 需要加密的文件路径
    picRead = open(filePath + "/" + name, "rb", 1024 * 1024 * 5)
    # 读取二进制文件
    picData = picRead.read()
    picLen = len(picData)
    # 需要写出的文件路径
    f_write = open(path2 + "/" + name, "wb", 1024 * 1024 * 5)
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
    # key = numpy.fromnumeric(data, dtype=numpy.uint8)
    key = np.array(data, dtype='uint8')
    i = 0
    # 进行循环加密
    for j in range(len(picData)):
        if i == len(key) - 1:
            newByte = key[i] ^ picData[j]
            i = 0
            # 写出二进制文件
            f_write.write(bytes([newByte]))
        else:
            newByte = key[i] ^ picData[j]
            i = i + 1
            # 写出二进制文件
            f_write.write(bytes([newByte]))
    f_write.close()
    picRead.close()


def ConnMongo():
    """
    连接数据库
    :return: mycol
    """
    myclient = pymongo.MongoClient("mongodb://64.64.241.63:7777/")
    mydb = myclient["magic"]
    mycol = mydb["cartoon"]
    return mycol


def WriteTxt(id, mycol):
    """
    写入txt
    :param id:
    :param mycol:
    :return:
    """
    info_path = os.getcwd()
    print(info_path + '/info.txt')
    info = open(info_path + '/info.txt', 'w')
    for x in mycol.find({"_id": id}, {"_id": 0, "chapters": [{"chapterUrl": 1}]}):
        info.write(json.dumps(x, indent=2, ensure_ascii=False))
    info.close()


def UpdTxt(upd_path_res_old_list, upd_path_res_list):
    """
    替换txt 内容
    :param upd_path_res_old_list: 旧的值
    :param upd_path_res_list: 新的值
    :return:
    """
    info_path = os.getcwd() + '/info.txt'
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
    print(remote_path)
    print(chapter)
    # /home/resources/wwwroot/limao/crypted/cartoonInfo/964cb5df130e327f708cb25dafc523a6/    115
    return remote_path, chapter


def UpdMongo(id, mycol):
    """
    修改数据库 漫画图片信息
    :param mycol:
    :return:
    """
    info_path = os.getcwd() + '/info.txt'
    with open(info_path, 'r') as file:
        file_data = file.readlines()
    condition = {'_id': id}
    print(type(condition))
    total_num_change = eval(((str(file_data).replace("'", "")).replace("\\n,", "")).strip("[]"))
    print(total_num_change)
    result = mycol.update_one(condition, {'$set': total_num_change})
    print("修改了：{}".format(result.modified_count))


@cal_time
def PicScaUpd(remote_path, chapter):
    """
    拼接漫画图片比例
    :param remote_path:
    :param chapter:
    :return:
    """
    try:
        upd_path_res_old_list = []
        upd_path_res_list = []
        show_files(remote_path, "/home/lighting/decode_pic", [])  # 创建文件 并解密到/home/lighting/decode_pic
        for i in range(1, chapter + 1):
            cur_remote_path = str("/home/lighting/decode_pic/" + str(i) + '/')
            print(cur_remote_path)
            # noinspection PyBroadException
            try:
                remote_files = os.listdir(cur_remote_path)
            except Exception:
                break
            for file in remote_files:  # 遍历读取远程目录里的所有文件
                scale = GetScale(cur_remote_path + file)  # 读取文件比例
                upd_path_res_old = (remote_path.split('crypted/', 1))[1] + str(i) + '/' + file
                upd_path_res = (remote_path.split('crypted/', 1))[1] + str(i) + '/' + file + "_" + str(scale)
                print(upd_path_res)
                upd_path_res_old_list.append(upd_path_res_old)
                upd_path_res_list.append(upd_path_res)

        UpdTxt(upd_path_res_old_list, upd_path_res_list)

    except Exception as e:
        raise e


def run(id):
    # 传入空的list接收文件名
    mycol = ConnMongo()
    WriteTxt(id, mycol)
    remote_path, chapter = Mongogetpath(id, mycol)
    PicScaUpd(remote_path, chapter)
    UpdMongo(id, mycol)


if __name__ == '__main__':
    run(22)
