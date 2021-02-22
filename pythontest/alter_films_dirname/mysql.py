import DBUtils.PooledDB
import pymysql
import traceback


def con_sql():
    try:
        pool = DBUtils.PooledDB.PooledDB(
            creator=pymysql,
            mincached=1,  # 连接池中空闲连接的初始数量
            maxcached=2,  # 连接池中空闲连接的最大数量
            maxconnections=5,  # 创建连接池的最大数量
            blocking=True,  # 超过最大连接数量时候的表现，为True等待连接数量下降，为false直接报错处理
            host="199.180.114.169",
            port=6033,
            user="root",
            passwd="ZsNice2020.",
            database="res_video_db"
        )
        conn = pool.connection()
        return conn
    except Exception:
        print('处理异常: ' + traceback.format_exc())
