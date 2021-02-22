import re
import uuid

import pymysql
from pymysql.cursors import DictCursor
import redis


def conn_mysql(info):
    """连接数据库"""
    conn = pymysql.connect(
        host=info.host,
        port=info.port,
        user=info.user,
        passwd=info.password,
        database=info.database,
    )
    cmd = conn.cursor(DictCursor)
    return cmd, conn


def close_mysql(cmd, conn):
    """关闭连接"""
    conn.commit()
    cmd.close()
    conn.close()


def select_user_id_whether_exists(info):
    """判断用户是否存在"""
    cmd, conn = conn_mysql(info)
    select_machine_code_sql = "select machine_code from {0}_user where user_id = {1}".format(info.app_name, info.user_id)
    cmd.execute(select_machine_code_sql)
    result = cmd.fetchall()
    if result:
        info.user_whether_exists = 1
        print("该用户存在！\"machine_code\":", result[0]["machine_code"])
    else:
        info.user_whether_exists = 0
        print("该用户不存在")
    close_mysql(cmd, conn)


def select_commodity(info):
    # 充值类型 1vip  2金币
    cmd, conn = conn_mysql(info)
    select_vip_sql = "select commodity_id,commodity_name,commodity_in_fact_amount from {0}_payment_commodity_info where commodity_type = 1  ".format(info.app_name)
    select_gold_sql = "select commodity_id,commodity_name,commodity_in_fact_amount from {0}_payment_commodity_info where commodity_type = 2  ".format(info.app_name)
    cmd.execute(select_vip_sql)
    info.commodity_vip_info = cmd.fetchall()
    cmd.execute(select_gold_sql)
    info.commodity_gold_info = cmd.fetchall()
    close_mysql(cmd, conn)


def change_user(info):
    cmd, conn = conn_mysql(info)
    select_check_machine_code_sql = "select machine_code from {0}_user where user_id = {1}".format(info.app_name, info.need_changed_user_id)
    cmd.execute(select_check_machine_code_sql)
    re_expression = re.compile("^TestCU-.*")
    re_machine_code = re.findall(re_expression, str((cmd.fetchall())[0]["machine_code"]))
    print(re_machine_code)
    if re_machine_code:
        return "is_changed"
    update_change_user_sql = "update {0}_user set machine_code =\"TestCU-{1}\" where user_id= {2} ".format(info.app_name, uuid.uuid1(), info.need_changed_user_id)
    cmd.execute(update_change_user_sql)
    close_mysql(cmd, conn)
    r = redis.Redis(host='199.180.114.169', port=9736, db=0, password="ZsAbc2020", decode_responses=True)
    res = r.delete("LOGIN_USER_INFO_{0}_{1}".format(info.need_changed_user_id, info.app_type))
    r.close()
    return res


def get_machine_code(info):
    cmd, conn = conn_mysql(info)
    cmd = conn.cursor(DictCursor)
    select_machine_code_sql = "select machine_code from {0}_user where user_id = {1}".format(info.app_name, info.user_id)
    execute_result_code = cmd.execute(select_machine_code_sql)
    if execute_result_code == 1:
        sql_result = cmd.fetchall()
        machine_code = sql_result[0]["machine_code"]
        print("machine_code: {0}".format(machine_code))
        info.machine_code = machine_code
    close_mysql(cmd, conn)
