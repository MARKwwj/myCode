import os
import time
import zipfile
import zlib
import codecs
import numpy as np


def cal_time(func):
    def wrapper(*args):
        start_time = time.time()
        res = func(*args)
        all_time = round(time.time() - start_time, 1)
        print("func {0} run time {1} s".format(func.__name__, all_time))
        return res

    return wrapper


def show_files(path, path2, all_files):
    # 修改txt文件编码为 utf-8
    # alter_utf(path)
    # 将文本压缩 为zip 再将后缀修改问 txt
    # zip_txt(path)
    # 创建加密后存储txt 的文件夹
    os.makedirs(path2)
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
            all_files.append(file)
    return all_files


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
    for j in range(picLen):
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


def alter_utf(file_path):
    os.chdir(file_path)
    for f in os.listdir():
        try:
            file1 = codecs.open(f, 'r', 'ansi', buffering=1024 * 1024 * 5)
            res = file1.read()
            file1.close()
            os.remove(f)
            file2 = codecs.open(f, 'w', encoding='utf-8', buffering=1024 * 1024 * 5)
            file2.write(res)
            file2.close()
        except Exception as e:
            print('alter_utf_Exception:' + e)


def zip_txt(file_path):
    # file_path = 'D:\\desktop\\秦时世界里的主神余孽'
    # 切换到需要压缩的文件 的路径 （压缩后的文件会生成到此目录下）
    os.chdir(file_path)
    for f in os.listdir():
        if f[0:5] == "cover":
            continue
        # zipfile.ZipFile(filename,mode)  filename 是压缩后的名字
        z = zipfile.ZipFile(f[:-4] + '.zip', mode="w", compression=zipfile.ZIP_DEFLATED)
        # f 需要压缩的文件
        z.write(f)
        # 删除原来的 txt 文件
        os.remove(f)
        z.close()
    for f in os.listdir():
        if f[0:5] == "cover":
            continue
            # os.rename(f, f[:-4] + '.jpg')
        os.rename(f, f[:-4] + '.txt')


if __name__ == '__main__':
    path = 'D:\\desktop\\zxz'
    for dir in os.listdir(path):
        novel_dir_path = os.path.join(path, dir)
        # alter_utf(novel_dir_path)
        zip_txt(novel_dir_path)
    # zip_txt('D:\\desktop\\333')
