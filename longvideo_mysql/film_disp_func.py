import logging
import os
import re
import shutil
import subprocess
import time

import pymysql

# from film_dispose.parameters import all_infos
# from film_dispose.conf import *
from pymysql.cursors import DictCursor

fmt = '[%(asctime)s] %(levelname)s [%(funcName)s: %(filename)s, %(lineno)d] %(message)s'
logging.basicConfig(level=logging.DEBUG,  # log level
                    filename='film_dis.txt',  # log名字
                    format=fmt,  # log格式
                    datefmt='%Y-%m-%d %H:%M:%S',  # 日期格式
                    filemode='a')  # 追加模式


def cal_time(func):
    def wrapper(*args):
        start_time = time.time()
        res = func(*args)
        all_time = round(time.time() - start_time, 1)
        logging.info("func {0} run time {1} s".format(func.__name__, all_time))
        return res

    return wrapper


@cal_time
def show_files(path, path2):
    if not os.path.exists(path2):
        os.makedirs(path2)
    # 首先遍历当前目录所有文件及文件夹
    file_list = os.listdir(path)
    # 准备循环判断每个元素是否是文件夹还是文件，是文件的话，把名称传入list，是文件夹的话，递归
    for file in file_list:
        # 利用os.path.join()方法取得路径全名，并存入cur_path变量，否则每次只能遍历一层目录
        cur_path = path + '/' + file
        # 判断是否是文件夹
        if os.path.isdir(cur_path):
            cur_path2 = path2 + '/' + file
            isExists = os.path.exists(cur_path2)
            if not isExists:
                os.makedirs(cur_path2)
                logging.info(cur_path2 + ' 创建成功')
            else:
                # 如果目录存在则不创建，并提示目录已存在
                logging.info(cur_path2 + ' 目录已存在')
            show_files(cur_path, cur_path2)
        else:
            encrypt(path, path2, file)
            # all_files.append(file)
    logging.info('{0} 加密完成！'.format(path))
    # return all_files


def encrypt(filePath, path2, name):
    # 需要加密的文件路径
    picRead = open(filePath + "/" + name, "rb", 1024 * 1024 * 5)
    # 读取二进制文件
    picData = picRead.read()
    picLen = len(picData)
    # 需要写出的文件路径
    f_write = open(path2 + "/" + name, "wb", 1024 * 1024 * 5)
    # 获取Key
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
    f_write.close()
    picRead.close()


def dispose_m3u8(films_dir):
    """
     1.获取视频 总时长
     2.截去 m3u8文件 ts/开头
    """
    m3u8_data = ''
    films_time = 0.0
    m3u8_path = films_dir + '/xindex.m3u8'
    # 判断么m3u8文件是否存在
    result = os.path.exists(m3u8_path)
    if result:
        with open(m3u8_path, 'r', 1024 * 10245 * 5) as fr:
            for line in fr:
                # 获取视频时长
                row_regex = re.compile(r'#EXTINF:\d+\.\d+,')
                row_res = re.findall(row_regex, line)
                if row_res:
                    time_regex = re.compile(r'\d+\.\d+')
                    film_time = re.findall(time_regex, row_res[0])
                    films_time += float(film_time[0])
                if 'ts/' in line:
                    # 截去 m3u8 内文件的 ts/前缀
                    line = line.replace('ts/', '')
                m3u8_data += line
        with open(m3u8_path, 'w', 1024 * 10245 * 5) as fw:
            fw.write(m3u8_data)
    films_time = int(films_time)
    return films_time


def cal_films_gold(films_time):
    """计算 视频 金币数"""
    goldCoin = 0
    if films_time < 60 * 30:
        goldCoin = 10
    elif films_time < 60 * 60:
        goldCoin = 20
    elif films_time < 60 * 120:
        goldCoin = 30
    else:
        goldCoin = 40
    return goldCoin


def get_film_info(film_dir):
    """获取视频标题，简介"""
    title_path = film_dir + '/标题.txt'
    title = ''
    introduce = ''
    # 判断么title_path文件是否存在
    result = os.path.exists(title_path)
    if result:
        with open(title_path, 'r', 1024 * 1024 * 5, encoding='utf-8') as fr:
            res = fr.read()
            # 匹配 标题
            title_regex = re.compile(r'标题开始\[.+\]标题结束')
            title_res = re.findall(title_regex, res)
            if title_res:
                title = ((str(title_res[0]).split('标题开始[', 1))[1].split(']标题结束'))[0]
            # 匹配 简介
            introduce_regex = re.compile(r'介绍开始\[.+\]介绍结束')
            introduce_res = re.findall(introduce_regex, res)
            if introduce_res:
                introduce = ((str(introduce_res[0]).split('介绍开始[', 1))[1].split(']介绍结束'))[0]
    return title, introduce


def cut_film_cover(film_dir):
    """视频截取封面图"""
    # 创建封面图目录
    isExists = os.path.exists(film_dir + '/covers')
    if not isExists:
        os.mkdir(film_dir + '/covers')
    # 遍历 电影目录下的文件
    i = 0
    width = ''
    height = ''
    for root, dirs, files in os.walk(film_dir):
        for file in files:
            if file[-2:] == 'ts':
                i = i + 1
                # ts文件路径
                film_path = film_dir + '/' + file
                # 获取 宽高
                if i == 1:
                    width, height = film_width_and_height(film_path)
                    if width is False:
                        logging.info('ts文件宽高获取失败! 当前目录:{0}'.format(film_path))
                        continue
                cmd = 'ffmpeg -i {0} -t 0.0001 -s {1}x{2} -y -vcodec libwebp {3}/covers/cover-{4}.webp'.format(film_path, width, height, film_dir, i)
                # logging.info('执行命令:{0}'.format(cmd))
                result = os.system(cmd)
                if result != 0:
                    logging.error('{0} {1}x{2} 封面图截取失败！！！'.format(film_path, width, height))
    # 获取封面图文件夹下的文件数量
    covers_quantity = cal_covers_quantity(film_dir)
    return covers_quantity


def cal_proportion(width, height):
    """ 宽高等比例缩小 """
    if height > width:
        proportion = 720 / float(height)
        height = 720
        width = int(width * proportion)
    else:
        proportion = 720 / float(width)
        width = 720
        height = int(height * proportion)
    width = str(width)
    height = str(height)
    return width, height


def film_width_and_height(film_path):
    """获取分片的宽高"""
    width = ''
    height = ''
    # film_path  分片的路径
    cmd = 'ffmpeg -i {0}'.format(film_path)
    try:
        # 判断么film_path文件是否存在
        film_info = subprocess.getstatusoutput(cmd)
        row_regex = re.compile(r'Stream #0:0.*')
        row_res = re.findall(row_regex, film_info[1])
        w_h_regex = re.compile(r'[1-9][0-9][0-9]*x[1-9][0-9][0-9]*')
        w_h_regex_res = re.findall(w_h_regex, row_res[0])
        width_and_height = w_h_regex_res[0]
        if width_and_height:
            width = int((width_and_height.split('x'))[0])
            height = int((width_and_height.split('x'))[1])
        # 修改图片比例
        cal_width, cal_height = cal_proportion(width, height)
    except Exception:
        return False, False
    return cal_width, cal_height


def cal_video_shard_size(film_dir):
    """计算分片文件总大小"""
    size = 0
    for root, dirs, files in os.walk(film_dir):
        for name in files:
            if name[-2:] == 'ts':
                size += os.path.getsize(os.path.join(root, name))
    return size


def cal_covers_quantity(film_dir):
    """计算封面图数量"""
    covers_path = film_dir + '/covers'
    # 获取封面图文件夹下的文件数量
    covers_list = os.listdir(covers_path)
    covers_quantity = len(covers_list)
    return covers_quantity


def set_cover_dir(film_dir, set_cover):
    """设置封面图专属文件夹"""
    covers_path = film_dir + '/covers'
    unique_cover_path = film_dir + '/setcover'
    isExists = os.path.exists(unique_cover_path)
    if not isExists:
        os.mkdir(unique_cover_path)
    cmd = 'cp {0}/{1} {2}/cover.webp'.format(covers_path, set_cover, unique_cover_path)
    logging.info('执行命令:{0}'.format(cmd))
    result = os.system(cmd)
    if result == 0:
        logging.info('设置封面图专属文件夹成功!')
    else:
        logging.error('设置封面图专属文件夹失败!')


def delete_files(film_dir):
    """删除未加密文件"""
    isExists = os.path.exists(film_dir)
    if isExists:
        shutil.rmtree(film_dir)


def jar_encrypt(old_path, new_path):
    result = os.system("java -jar XORDir.jar {0} {1}".format(old_path, new_path))
    if result == 0:
        os.rename(old_path, old_path + "-altered")
        return "加密成功"
    else:
        return "加密失败"


def check_whether_encrypt(film_dir):
    expression = re.compile(".*-altered$")
    result = re.findall(expression, film_dir)
    if result:
        logging.info('视频已完成加密 跳过！目录{0}'.format(film_dir))
        return "is_encrypted"
    else:
        return "un_encrypt"


def cal_covers_num(covers_path, film):
    """更新封面图数量到数据库"""
    conn = pymysql.connect(
        host="58.82.232.37",
        port=6033,
        user='root',
        passwd='ZsNice2020.',
        database='res_video_db',
    )
    cmd = conn.cursor(DictCursor)
    num = len(os.listdir(covers_path))
    sql = "update long_video_video_info set video_cover_quantity = {0} where video_id = {1}".format(num, film)
    cmd.execute(sql)
    """关闭连接"""
    conn.commit()
    cmd.close()
    conn.close()


def run(cur_video_path, res_video_path, video_id):
    logging.info('视频处理开始！！！')
    # 当前电影的路径
    video_dir = cur_video_path + '/' + video_id
    logging.info('当前电影的路径:{0}'.format(video_dir))
    # 验证是否已加密
    encrypt_check = check_whether_encrypt(video_dir)
    if encrypt_check == "is_encrypted":
        return
    covers_path = video_dir + '/covers'
    isExists_covers = os.path.exists(covers_path)
    if isExists_covers:
        logging.info("covers文件夹已存在! 跳过截取封面图！")
    else:
        logging.info('正在截取封面图！')
        cut_film_cover(video_dir)
        logging.info('封面图截取结束！')
    # logging.info('计算封面图数量 更新到数据库！')
    # cal_covers_num(covers_path, video_id)
    # logging.info('更新完成！')
    logging.info('开始加密！')
    new_path = os.path.join(res_video_path, video_id)
    result = jar_encrypt(video_dir, new_path)
    logging.info(result)
    logging.info('加密结束！')

    logging.info('当前视频处理完成！！！')

# cur_video_path = "/home/datadrive/resources/res9/smallvideo/crypted/videos"
# res_video_path = "/home/datadrive/resources/res1/da_xiang/crypted/videos"
# cur_video_path="D:\\desktop\\r1"
# res_video_path="D:\\desktop\\r2"
# run(cur_video_path, res_video_path, 1)
