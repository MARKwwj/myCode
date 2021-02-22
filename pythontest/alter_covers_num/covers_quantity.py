import os
import time

import pymysql
from pymysql.cursors import DictCursor


def conn_mysql():
    """连接数据库"""
    conn = pymysql.connect(
        host="58.82.232.37",
        port=6033,
        user='root',
        passwd='ZsNice2020.',
        database='res_video_db',
    )
    cmd = conn.cursor(DictCursor)
    return cmd, conn


def close_mysql(cmd, conn):
    """关闭连接"""
    conn.commit()
    cmd.close()
    conn.close()


def update_covers_num(txt_name):
    cmd, conn = conn_mysql()
    with open(txt_name) as fr:
        for line in fr.readlines():
            video_id = line.strip("\r\n")
            # cur_video_covers_path = os.path.join(videos_path, video_id, "covers")
            # num = len(os.listdir(cur_video_covers_path))
            num = 1
            sql = "update long_video_video_info set video_cover_quantity = {0} where video_id = {1}".format(num, video_id)
            print(sql)
    #         cmd.execute(sql)
    # close_mysql(cmd, conn)


videos_path = "/home/datadrive/resources/res1/da_xiang/crypted/videos"
update_covers_num("1.txt")
