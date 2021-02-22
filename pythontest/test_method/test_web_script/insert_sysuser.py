from datetime import datetime
import logging
import random
import time
import uuid
import pymysql

fmt = '[%(asctime)s] %(levelname)s [%(funcName)s: %(filename)s, %(lineno)d] %(message)s'
logging.basicConfig(level=logging.DEBUG,  # log level
                    # filename='insert_log.txt',  # log名字
                    format=fmt,  # log格式
                    datefmt='%Y-%m-%d %H:%M:%S',  # 日期格式
                    filemode='a')  # 追加模式


class user_infos:
    def __init__(self):
        """sys_user 表字段"""
        self.app_type = '2'
        self.phone_udid = str(uuid.uuid1())
        self.user_name = 'test-user'
        self.nick_name = self.get_nick_name()
        self.user_type = '02'
        self.sex = random.randint(0, 1)
        self.avatar = self.get_avatar()
        self.create_time = str(datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
        self.update_time = str(datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
        city_list = ['广东·广州', '上海', '北京', '重庆', '四川·成都', '湖南·长沙', '广东·深圳', '浙江·宁波', '江苏·苏州', '浙江·杭州']
        self.city = (random.sample(city_list, 1))[0]
        str_bir = str(random.randint(1980, 2005)) + '-' + '0' + str(random.randint(1, 9)) + '-' + str(random.randint(10, 28))
        self.birthday = str_bir
        self.mobile_type = 'small-video-test-user'
        self.equipment_type = str(random.randint(1, 2))

    def get_nick_name(self):
        name_txt_path = "D:\\desktop\\昵称.txt"
        nick_name_Lsit = []
        with open(name_txt_path, 'r', encoding='utf-8') as fr:
            for line in fr.readlines():
                res = line.strip('\n')
                nick_name_Lsit.append(res)
        nick_name = (random.sample(nick_name_Lsit, 1))[0]
        return nick_name

    def get_avatar(self):
        avatar_url_path = "D:\\desktop\\头像.txt"
        avatar_url_path_list = []
        with open(avatar_url_path, 'r', encoding='utf-8') as fr:
            for line in fr.readlines():
                res = line.strip('\n')
                avatar_url_path_list.append(res)
        avatar = (random.sample(avatar_url_path_list, 1))[0]
        return avatar


def insert_sql():
    conn = pymysql.connect(
        host='192.168.100.51',
        port=3306,
        user='root',
        password='123456',
        database='jtest',
    )
    try:
        cmd = conn.cursor()
        for i in range(100):
            info = user_infos()
            sql = "insert into sys_user(" \
                  "app_type,phone_udid,user_name," \
                  "nick_name,user_type,sex," \
                  "avatar,create_time,update_time," \
                  "city,birthday,mobile_type,equipment_type) " \
                  "values(\'{0}\',\'{1}\',\'{2}\',\'{3}\',\'{4}\',\'{5}\',\'{6}\'," \
                  "\'{7}\',\'{8}\',\'{9}\',\'{10}\',\'{11}\',\'{12}\')".format(
                info.app_type, info.phone_udid, info.user_name,
                info.nick_name, info.user_type, info.sex,
                info.avatar, info.create_time, info.update_time,
                info.city, info.birthday, info.mobile_type,
                info.equipment_type)
            cmd.execute(sql)
        conn.commit()
    except Exception as e:
        logging.error(e)
    finally:
        conn.close()


insert_sql()
