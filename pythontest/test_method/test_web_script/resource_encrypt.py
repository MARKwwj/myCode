import os
import time


def cal_time(func):
    def wrapper(*args):
        start_time = time.time()
        res = func(*args)
        all_time = round(time.time() - start_time, 1)
        print("func {0} run time {1} s".format(func.__name__, all_time))
        return res

    return wrapper


@cal_time
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
    # data = [89, 61, 148, 105, 40, 237, 255, 227, 195, 111, 172, 197, 236, 82, 43, 163,
    #         121, 78, 245, 50, 67, 117, 136, 3, 26, 132, 52, 204, 182, 83, 13, 146,
    #         21, 215, 47, 48, 189, 96, 181, 23, 1, 158, 219, 184, 86, 112, 51, 84,
    #         46, 59, 191, 106, 140, 4, 65, 173, 109, 248, 88, 53, 152, 153, 36, 115,
    #         37, 242, 177, 90, 179, 196, 143, 217, 239, 251, 69, 193, 55, 42, 147, 76,
    #         134, 218, 9, 174, 141, 216, 138, 129, 124, 68, 107, 234, 249, 102, 64, 169,
    #         185, 122, 56, 229, 41, 252, 127, 18, 79, 205, 186, 199, 108, 214, 160, 16,
    #         226, 94, 240, 209, 253, 187, 166, 99, 5, 213, 34, 155, 154, 0, 211, 97,
    #         72, 175, 238, 167, 70, 119, 31, 113, 222, 2, 66, 156, 162, 194, 207, 206,
    #         157, 100, 246, 202, 171, 20, 54, 11, 208, 220, 62, 126, 14, 114, 95, 32,
    #         73, 165, 254, 35, 244, 81, 149, 137, 135, 176, 93, 15, 44, 57, 164, 92,
    #         241, 19, 203, 87, 6, 243, 235, 151, 12, 24, 183, 33, 188, 144, 224, 178,
    #         150, 27, 39, 233, 116, 25, 103, 110, 201, 85, 192, 159, 232, 190, 212, 123,
    #         131, 22, 210, 45, 74, 10, 125, 180, 130, 77, 221, 133, 60, 230, 80, 75,
    #         198, 128, 247, 29, 231, 118, 63, 250, 225, 120, 168, 104, 161, 28, 145, 223,
    #         170, 58, 38, 8, 142, 98, 71, 91, 30, 101, 200, 7, 17, 139, 49, 228]
    data = [0x59, 0x3d, 0x94, 0x69, 0x28, 0xed, 0xff, 0xe3, 0xc3, 0x6f, 0xac, 0xc5, 0xec, 0x52, 0x2b, 0xa3,
            0x79, 0x4e, 0xf5, 0x32, 0x43, 0x75, 0x88, 0x03, 0x1a, 0x84, 0x34, 0xcc, 0xb6, 0x53, 0x0d, 0x92,
            0x15, 0xd7, 0x2f, 0x30, 0xbd, 0x60, 0xb5, 0x17, 0x01, 0x9e, 0xdb, 0xb8, 0x56, 0x70, 0x33, 0x54,
            0x2e, 0x3b, 0xbf, 0x6a, 0x8c, 0x04, 0x41, 0xad, 0x6d, 0xf8, 0x58, 0x35, 0x98, 0x99, 0x24, 0x73,
            0x25, 0xf2, 0xb1, 0x5a, 0xb3, 0xc4, 0x8f, 0xd9, 0xef, 0xfb, 0x45, 0xc1, 0x37, 0x2a, 0x93, 0x4c,
            0x86, 0xda, 0x09, 0xae, 0x8d, 0xd8, 0x8a, 0x81, 0x7c, 0x44, 0x6b, 0xea, 0xf9, 0x66, 0x40, 0xa9,
            0xb9, 0x7a, 0x38, 0xe5, 0x29, 0xfc, 0x7f, 0x12, 0x4f, 0xcd, 0xba, 0xc7, 0x6c, 0xd6, 0xa0, 0x10,
            0xe2, 0x5e, 0xf0, 0xd1, 0xfd, 0xbb, 0xa6, 0x63, 0x05, 0xd5, 0x22, 0x9b, 0x9a, 0x00, 0xd3, 0x61,
            0x48, 0xaf, 0xee, 0xa7, 0x46, 0x77, 0x1f, 0x71, 0xde, 0x02, 0x42, 0x9c, 0xa2, 0xc2, 0xcf, 0xce,
            0x9d, 0x64, 0xf6, 0xca, 0xab, 0x14, 0x36, 0x0b, 0xd0, 0xdc, 0x3e, 0x7e, 0x0e, 0x72, 0x5f, 0x20,
            0x49, 0xa5, 0xfe, 0x23, 0xf4, 0x51, 0x95, 0x89, 0x87, 0xb0, 0x5d, 0x0f, 0x2c, 0x39, 0xa4, 0x5c,
            0xf1, 0x13, 0xcb, 0x57, 0x06, 0xf3, 0xeb, 0x97, 0x0c, 0x18, 0xb7, 0x21, 0xbc, 0x90, 0xe0, 0xb2,
            0x96, 0x1b, 0x27, 0xe9, 0x74, 0x19, 0x67, 0x6e, 0xc9, 0x55, 0xc0, 0x9f, 0xe8, 0xbe, 0xd4, 0x7b,
            0x83, 0x16, 0xd2, 0x2d, 0x4a, 0x0a, 0x7d, 0xb4, 0x82, 0x4d, 0xdd, 0x85, 0x3c, 0xe6, 0x50, 0x4b,
            0xc6, 0x80, 0xf7, 0x1d, 0xe7, 0x76, 0x3f, 0xfa, 0xe1, 0x78, 0xa8, 0x68, 0xa1, 0x1c, 0x91, 0xdf,
            0xaa, 0x3a, 0x26, 0x08, 0x8e, 0x62, 0x47, 0x5b, 0x1e, 0x65, 0xc8, 0x07, 0x11, 0x8b, 0x31, 0xe4]
    i = 0
    # 进行循环加密
    for j in range(picLen):
        newByte = data[i] ^ picData[j]
        i = (i + 1) % len(data)
        f_write.write(bytes([newByte]))
        # if i == len(data) - 1:
        #     newByte = data[i] ^ picData[j]
        #     i = 0
        #     # 写出二进制文件
        #     f_write.write(bytes([newByte]))
        # else:
        #     newByte = data[i] ^ picData[j]
        #     i = i + 1
        #     # 写出二进制文件
        #     f_write.write(bytes([newByte]))

    f_write.close()
    picRead.close()


if __name__ == '__main__':
    # 解密本地文件
    # 传入空的list接收文件名
    contents = show_files("D:\\desktop\\rr", "D:\\desktop\\rr1", [])
    # 循环打印show_files函数返回的文件名列表
    # for content in contents:
    #     print(content)
