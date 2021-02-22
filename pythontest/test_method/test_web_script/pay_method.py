import hashlib
import json
import time
import urllib.error
import urllib.parse
import urllib.request

from pymysql.cursors import DictCursor


class Settings:
    def __init__(self, user_id):
        """初始化 支付参数"""
        self.channel_id = '27'
        self.vip_id = '12'
        self.app_id = '1'
        self.user_id = user_id
        self.token = 'eyJhbGciOiJIUzUxMiJ9.eyJsb2dpbl91c2VyX2tleSI6IjUyZjVhYTQ4LTgwZDAtNDA1OC1iM2RlLTMwMTBjMTEyZWEyOSJ' \
                     '9.w0qu_6SygJ8D5JaXICu8z_JwxUhLlyJQ8KPFhxRW8u5OkZhxCE4G-PAhcNKqmVDJNhC049UWSsY2vOEyeKbuag'
        # 订单状态
        self.status = '1'
        # 金额
        self.total_fee = '100'
        # 支付类型
        self.type_pay = '1'
        # 商户id
        self.mch_id = '15520208888'
        # 商户密钥
        self.key = 'e07a7c0bedc507a53172e76d4b60af3c'


def creat_md5sign(pay_set):
    """生成订单签名"""
    # String check = Md5Utils.hash(Md5Utils.hash(channelId + vipId + appId + userId));
    params = pay_set.channel_id + pay_set.vip_id + pay_set.app_id + pay_set.user_id
    sign_create = hashlib.md5((hashlib.md5(params.encode(encoding='UTF-8')).hexdigest()).encode(encoding='UTF-8')).hexdigest()
    print('签名①：' + sign_create)
    return sign_create


def backparams(pay_set):
    bp_key = ["appId", "key", "mchId", "userId", "vipId"]
    bp_value = [pay_set.app_id, pay_set.key, pay_set.mch_id, int(pay_set.user_id), int(pay_set.vip_id)]
    bp_dict = dict(zip(bp_key, bp_value))
    # 转为json 串
    back_params = json.dumps(bp_dict, separators=(',', ':'))
    mdbparams = hashlib.md5(back_params.encode(encoding='UTF-8')).hexdigest()
    print('back_params：' + mdbparams)
    return mdbparams


def create_md5backsign(pay_set, out_trade_no):
    """回调订单签名"""
    params = out_trade_no + pay_set.total_fee + pay_set.mch_id + pay_set.key
    sign_create = hashlib.md5(params.encode(encoding='utf-8')).hexdigest()
    print('签名②：' + sign_create)
    return sign_create


def create_order(pay_set, sign):
    """生成订单"""
    url_pay = 'http://80.251.211.152:8803/pay/service/order/create/{}/{}/{}/{}/{}'.format(pay_set.channel_id, pay_set.vip_id, pay_set.app_id, pay_set.user_id, sign)
    headers = {
        'Authorization': 'Bearer {}'.format(pay_set.token),
        'User-Agent': 'zhiyefashiliuhaizhu'
    }

    try:
        request = urllib.request.Request(url_pay, headers=headers, method='POST')
        response = urllib.request.urlopen(request)
    except urllib.error.URLError as e:
        print(e.reason)
    finally:
        order_no = response.read().decode()
        print('订单：' + order_no)
        return order_no


def back_ldpay(pay_set, sign_back, out_trade_no, back_params):
    """零度 支付 回调"""
    result = ""
    url_payback = 'http://80.251.211.152:8803/pay/service/ld/notify' \
                  '?status={}&total_fee={}&out_trade_no={}' \
                  '&back_params={}&type={}&sign={}'. \
        format(pay_set.status, pay_set.total_fee, out_trade_no, back_params, pay_set.type_pay, sign_back)
    try:
        response_payback = urllib.request.urlopen(url_payback)
        result = response_payback.read().decode()
        print('结果：' + result)
    except Exception as e:
        print("请求异常 10s 后重试！！！")
        time.sleep(10)
        print("尝试再次请求！！！")
        response_payback = urllib.request.urlopen(url_payback)
        result = response_payback.read().decode()
        print('结果：' + result)
    finally:
        return result


def pay_run(user_id):
    pay_set = Settings(user_id)
    # 生成订单签名
    sign = creat_md5sign(pay_set)
    # 生成订单
    out_trade_no = create_order(pay_set, sign)
    # 生成回调签名
    sign_back = create_md5backsign(pay_set, out_trade_no)
    # 回调参数
    back_params = backparams(pay_set)
    # 回调
    back_ldpay(pay_set, sign_back, out_trade_no, back_params)


def select_userid(bus_id):
    import pymysql
    import traceback
    try:
        conn = pymysql.connect(
            host="80.251.211.152",
            port=6033,
            user="root",
            passwd="ZsNice2020.",
            database="manage_db"
        )
        print("数据库连接成功!!! 当前数据库版本信息 %s " % conn.get_server_info())
        cmd = conn.cursor(DictCursor)
        sql = "select user_id from sys_user where ancestors like \"%{0}%\" and user_type = \"02\"".format(bus_id)
        cmd.execute(sql)
        res = cmd.fetchall()
        n = 0
        userid = []
        while n < len(res):
            userid.append(res[n]["user_id"])
            n += 1
        return userid
    except Exception:
        print("处理异常 " + traceback.format_exc())  # 打印异常信息
    finally:
        conn.close()


if __name__ == '__main__':
    # 传入商务id 自动查找所有下级代理的下级用户
    id_all = select_userid(185360)
    id_all.reverse()
    # 打印所有下级用户
    print(id_all)
    # 循环给每个下级用户进行充值
    # id_all = [179728, 179729, 179730, 179731, 179732, 179733, 179734, 179735, 179736, 179737, 179738, 179739]
    for id in id_all:
        id_str = str(id)
        print("充值账号%s" % id_str)
        print(id_all.index(id) + 1)
        pay_run(id_str)
