from datetime import datetime
import logging
import random
import time

import pymysql

fmt = '[%(asctime)s] %(levelname)s [%(funcName)s: %(filename)s, %(lineno)d] %(message)s'
logging.basicConfig(level=logging.DEBUG,  # log level
                    # filename='insert_log.txt',  # log名字
                    format=fmt,  # log格式
                    datefmt='%Y-%m-%d %H:%M:%S',  # 日期格式
                    filemode='a')  # 追加模式


class all_infos:
    def __init__(self):
        self.host = '192.168.100.51'
        self.port = 3306
        self.user = 'root'
        self.password = '123456'
        self.database = 'long_video_db'

        """sql 字段值"""
        self.user_id = 1042980
        self.nick_name = 'Hello World'
        self.commodity_name = '月卡'
        self.commodity_value = 10
        self.commodity_gift = 10
        self.commodity_in_fact_amount = 10
        self.commodity_type = random.randint(1, 2)
        self.order_status = random.randint(0, 3)
        self.recharge_method = random.randint(0, 2)
        str_order = '2020925' + str(round(time.time() * 1000))
        self.order_number = str_order
        self.third_party_order_number = str_order
        self.del_flag = 2
        month = str(random.randint(1, 12))
        day = str(random.randint(1, 28))
        str_p = '2020-{0}-{1} 16:06:23'.format(month, day)
        self.create_time = str_p
        self.success_time = str_p
        self.coin_balance = 0
        self.vip_end_time = str(datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
        self.recharge_amount = 1.00
        channel_list = [1, 2, 3, 4, 5, 6, 7, 8, 11, 12, 14, 15, 16, 17, 18, 19, 20, 21, 23, 24, 25,
                        27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 43, 44, 45, 46, 47, 48]
        select_channel_id = random.sample(channel_list, 1)
        self.channel_id = select_channel_id[0]

        # logging.info('**************************************************************************')
        # logging.info('user_id:{0}'.format(self.user_id))
        # logging.info('nick_name:{0}'.format(self.nick_name))
        # logging.info('commodity_name:{0}'.format(self.commodity_name))
        # logging.info('commodity_value:{0}'.format(self.commodity_value))
        # logging.info('commodity_gift:{0}'.format(self.commodity_gift))
        # logging.info('commodity_in_fact_amount:{0}'.format(self.commodity_in_fact_amount))
        # logging.info('commodity_type:{0}'.format(self.commodity_type))
        # logging.info('order_status:{0}'.format(self.order_status))
        # logging.info('recharge_method:{0}'.format(self.recharge_method))
        # logging.info('order_number:{0}'.format(self.order_number))
        # logging.info('third_party_order_number:{0}'.format(self.third_party_order_number))
        # logging.info('del_flag:{0}'.format(self.del_flag))
        # logging.info('create_time:{0}'.format(self.create_time))
        # logging.info('success_time:{0}'.format(self.success_time))
        # logging.info('coin_balance:{0}'.format(self.coin_balance))
        # logging.info('vip_end_time:{0}'.format(self.vip_end_time))
        # logging.info('recharge_amount:{0}'.format(self.recharge_amount))
        # logging.info('channel_id:{0}'.format(self.channel_id))
        # logging.info('**************************************************************************')


def insert_sql():
    conn = pymysql.connect(
        host='192.168.100.51',
        port=3306,
        user='root',
        password='123456',
        database='long_video_db',
    )
    try:
        cmd = conn.cursor()
        for i in range(500000):
            info = all_infos()
            sql = 'insert into long_video_payment_order_record(' \
                  'user_id,nick_name,commodity_name,' \
                  'commodity_value,commodity_gift,commodity_in_fact_amount,' \
                  'commodity_type,order_status,recharge_method,' \
                  'order_number,third_party_order_number,del_flag,' \
                  'create_time,success_time,coin_balance,' \
                  'vip_end_time,recharge_amount,channel_id) values(' \
                  '{0},\'{1}\',\'{2}\',' \
                  '{3},{4},{5},' \
                  '{6},{7},{8},' \
                  '\'{9}\',\'{10}\',{11},' \
                  '\'{12}\',\'{13}\',{14},' \
                  '\'{15}\',{16},{17})'.format(info.user_id, info.nick_name, info.commodity_name,
                                               info.commodity_value, info.commodity_gift, info.commodity_in_fact_amount,
                                               info.commodity_type, info.order_status, info.recharge_method,
                                               info.order_number, info.third_party_order_number, info.del_flag,
                                               info.create_time, info.success_time, info.coin_balance,
                                               info.vip_end_time, info.recharge_amount, info.channel_id)
            cmd.execute(sql)
        conn.commit()
    except Exception as e:
        logging.error(e)
    finally:
        conn.close()


if __name__ == '__main__':
    insert_sql()
