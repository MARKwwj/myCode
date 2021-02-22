import hashlib
import urllib.request
import urllib.parse
import urllib.error
import json


class Settings:
    def __init__(self):
        """初始化 支付参数"""
        self.channel_id = '27'
        self.vip_id = '12'
        self.app_id = '1'
        self.user_id = '942256'
        self.token = "eyJhbGciOiJIUzUxMiJ9.eyJsb2dpbl91c2VyX2tleSI6ImU4NTgwOTMyLWQ2YzctNGRmZC1iNGM5LTA3ZDFkZjViMTNiM" \
                     "CJ9.QsNOOLy-RQsHEolm8ArBl4V1LupS7EvtfyoVIGjk5DZ12oA11ZaL7WM3zcQwxLYbhsC8tqLPSqEd17ZeunzhaQ"
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
    url_pay = 'http://192.168.100.105:8803/pay/service/order/create/{}/{}/{}/{}/{}'.format(pay_set.channel_id, pay_set.vip_id, pay_set.app_id, pay_set.user_id, sign)
    headers = {
        'Authorization': 'Bearer {}'.format(pay_set.token),
        'User-Agent': 'zhiyefashiliuhaizhu'
    }

    try:
        request = urllib.request.Request(url_pay, headers=headers, method='POST')
        response = urllib.request.urlopen(request)
    except urllib.error.URLError as e:
        print(e.reason)
    order_no = response.read().decode()
    print('订单：' + order_no)
    return order_no


def back_ldpay(pay_set, sign_back, out_trade_no, back_params):
    """零度 支付 回调"""
    url_payback = 'http://192.168.100.105:8803/pay/service/ld/notify' \
                  '?status={}&total_fee={}&out_trade_no={}' \
                  '&back_params={}&type={}&sign={}'. \
        format(pay_set.status, pay_set.total_fee, out_trade_no, back_params, pay_set.type_pay, sign_back)
    response_payback = urllib.request.urlopen(url_payback)
    result = response_payback.read().decode()
    print('结果：' + result)
    return result


def pay_run():
    pay_set = Settings()
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


if __name__ == '__main__':
    for i in range(100):
        pay_run()
        print(i)
