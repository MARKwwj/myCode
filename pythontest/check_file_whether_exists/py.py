import os

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


def check():
    cmd, conn = conn_mysql()
    sql = "select video_id from long_video_video_info where video_creator = 'reptile' and video_create_time < '2021-01-08 00:00:00'"
    cmd.execute(sql)
    for id in cmd.fetchall():
        path = "/home/datadrive/resources/res1/da_xiang/crypted/videos/{0}".format(id["video_id"])
        try:
            os.stat(path)
        except Exception as e:
            print("video_id:{0},err:{1}".format(id["video_id"], e))
    close_mysql(cmd, conn)


check()
