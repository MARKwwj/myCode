import DBUtils.PooledDB
import pymysql
import traceback
from pymysql.cursors import DictCursor
import uuid


class Gdata:
    def __init__(self):
        self.pool = DBUtils.PooledDB.PooledDB(
            creator=pymysql,
            mincached=1,
            maxcached=2,
            maxconnections=5,
            blocking=True,
            host="80.251.211.152",
            port=6033,
            user="root",
            passwd="ZsNice2020.",
            database="manage_db"
        )
        self.conn = self.pool.connection()

    def sql_get_card_name(self, name, type):
        self.is_connected()
        cmd = self.conn.cursor(DictCursor)
        cmd2 = self.conn.cursor(DictCursor)
        if type == 0:
            sql_card_name = 'SELECT	b.card_name FROM app_type_config a,	app_pay_vip b WHERE	a.app_id = b.app_id AND ( type = 1 OR type = 0 ) ' \
                            'AND b.del_flag = 0 	AND b.app_category = (select app_category from app_type_config where app_name =\'{0}\' )' \
                            'AND b.type = \'{1}\' ORDER BY	b.sort DESC '.format(name, type)
            sql_card_money = 'select b.vip_actual_money from app_type_config a,app_pay_vip b ' \
                             ' where a.app_id =b.app_id and (type = 1 or type = 0) and b.del_flag = 0 and' \
                             ' b.app_category = (select app_category from app_type_config where app_name =\'{0}\' ) and b.type =\'{1}\' ORDER BY b.sort DESC'.format(name, type)
        else:
            sql_card_name = 'select b.card_name,b.gold_money from app_type_config a,app_pay_vip b  where' \
                            ' a.app_id =b.app_id and (type = 1 or type = 0) and b.del_flag = 0 and ' \
                            'b.app_category = (select app_category from app_type_config where app_name =\'{0}\' ) and b.type =\'{1}\' ORDER BY b.sort DESC'.format(name, type)
            sql_card_money = 'select b.gold_money from app_type_config a,app_pay_vip b  where' \
                             ' a.app_id =b.app_id and (type = 1 or type = 0) and b.del_flag = 0 and ' \
                             'b.app_category = (select app_category from app_type_config where app_name =\'{0}\' ) and b.type =\'{1}\' ORDER BY b.sort DESC'.format(name, type)
        cmd.execute(sql_card_name)
        cmd2.execute(sql_card_money)
        card_money = cmd2.fetchall()
        card_name = cmd.fetchall()
        n = 0
        card_name_list = []
        card_money_list = []
        if type == 0:
            while n < len(card_name):
                card_name_list.append(card_name[n]['card_name'])
                card_money_list.append(card_money[n]['vip_actual_money'])
                n = n + 1
        else:
            while n < len(card_name):
                card_name_list.append(card_name[n]['card_name'])
                card_money_list.append(card_money[n]['gold_money'])
                n = n + 1
        cmd.close()
        cmd2.close()
        return card_name_list, card_money_list

    def get_app_nameid_list(self):
        self.is_connected
        cmd = self.conn.cursor(DictCursor)
        sql_appname = 'select app_name from app_type_config where app_type !=0'
        cmd.execute(sql_appname)
        app_nameid_list = cmd.fetchall()
        card_name_list = []
        n = 0
        while n < len(app_nameid_list):
            card_name_list.append(app_nameid_list[n]['app_name'])
            n = n + 1
        cmd.close()
        return card_name_list

    def sql_get_appcardid(self, appname, cardname):
        self.is_connected
        cmd = self.conn.cursor()
        sql = 'SELECT	b.id,b.app_id FROM	app_type_config a,	' \
              'app_pay_vip b WHERE	a.app_id = b.app_id AND b.del_flag = 0 	AND ' \
              'b.app_category = (select app_category from app_type_config where app_name =\'{0}\')	AND b.card_name =\'{1}\' '.format(appname, cardname)
        cmd.execute(sql)
        res = cmd.fetchall()
        vipid = str(res[0][0])
        appid = str(res[0][1])
        cmd.close()
        print(vipid)
        print(appid)
        return appid, vipid

    def sql_getidexist(self, userid, app_name):
        self.is_connected
        cmd = self.conn.cursor()
        sql = "select b.app_type,b.user_id from " \
              "app_type_config a,sys_user b where a.app_type =b.app_type " \
              "and user_Id =\'{0}\' and a.app_name =\'{1}\' ".format(userid, app_name)
        cmd.execute(sql)
        res = cmd.fetchall()
        cmd.close()
        return res

    def change_userid_exist(self, userid):
        self.is_connected
        cmd = self.conn.cursor()
        sql = "select b.app_type,b.user_id from " \
              "app_type_config a,sys_user b where a.app_type =b.app_type " \
              "and user_Id =\'{0}\'".format(userid)
        cmd.execute(sql)
        res = cmd.fetchall()
        cmd.close()
        return res

    def change_userid(self, userid):
        self.is_connected
        cmd = self.conn.cursor()
        phone_uuid = uuid.uuid1()
        sql = "UPDATE sys_user set phone_udid =\'{0}\' where user_id = \'{1}\'".format(phone_uuid, userid)
        cmd.execute(sql)
        self.conn.commit()
        res = cmd.rowcount
        cmd.close()
        return res

    def is_connected(self):
        """Check if the server is alive"""
        try:
            self.conn.ping(reconnect=True)
            print("db is connecting")
        except:
            traceback.print_exc()
            self.conn = pymysql.connect(
                host="80.251.211.152",
                port=6033,
                user="root",
                passwd="ZsNice2020.",
                database="manage_db"
            )
            print("db reconnect")
