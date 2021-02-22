import hashlib
import urllib.request
import urllib.parse
import urllib.error
import json
from paysettings import Paysettings
from get_data import Gdata
import tkinter

result = ""


def backparams(userid, appid, vipid, pay_set):
    bp_key = ["appId", "key", "mchId", "userId", "vipId"]
    bp_value = [appid, pay_set.key, pay_set.mch_id, int(userid), int(vipid)]
    bp_dict = dict(zip(bp_key, bp_value))
    # 转为json 串
    back_params = json.dumps(bp_dict, separators=(',', ':'))
    print(back_params)
    mdbparams = hashlib.md5(back_params.encode(encoding='UTF-8')).hexdigest()
    print('back_params：' + mdbparams)
    return mdbparams


def creat_md5sign(userid, appid, vipid, pay_set):
    """生成订单签名"""
    # String check = Md5Utils.hash(Md5Utils.hash(channelId + vipId + appId + userId));
    params = pay_set.channel_id + vipid + appid + userid
    sign_create = hashlib.md5((hashlib.md5(params.encode(encoding='UTF-8')).hexdigest()).encode(encoding='UTF-8')).hexdigest()
    print('签名①：' + sign_create)
    return sign_create


def create_md5backsign(money, pay_set, out_trade_no):
    """回调订单签名"""
    params = out_trade_no + money + pay_set.mch_id + pay_set.key
    sign_create = hashlib.md5(params.encode(encoding='utf-8')).hexdigest()
    print('签名②：' + sign_create)
    return sign_create


def create_order(userid, appid, vipid, pay_set, sign):
    """生成订单"""
    url_pay = 'http://80.251.211.152:8803/pay/service/order/create/{}/{}/{}/{}/{}' \
        .format(pay_set.channel_id, vipid, appid, userid, sign)
    headers = {
        'Authorization': 'Bearer {}'.format(pay_set.token),
        'User-Agent': 'zhiyefashiliuhaizhu'
    }

    try:
        request = urllib.request.Request(url_pay, headers=headers, method='POST')
        response = urllib.request.urlopen(request)
        order_no = response.read().decode()
        print('订单：' + order_no)
        return order_no
    except urllib.error.HTTPError as e:
        print(e.reason)


def back_ldpay(money, pay_set, sign_back, out_trade_no, back_params):
    """零度 支付 回调"""
    global result
    result = ""
    url_payback = 'http://80.251.211.152:8803/pay/service/ld/notify' \
                  '?status={}&total_fee={}&out_trade_no={}' \
                  '&back_params={}&type={}&sign={}'. \
        format(pay_set.status, money, out_trade_no, back_params, pay_set.type_pay, sign_back)
    try:
        response_payback = urllib.request.urlopen(url_payback)
        result = response_payback.read().decode()
        print('结果：' + result)
        return result
    except urllib.error.HTTPError as e:
        print(e.reason)


def pay_run(GUI):
    pay_set = Paysettings()
    gdata = Gdata()
    appname = GUI.app_name_com.get()
    cardname = GUI.Point_card_name_com.get()
    money = GUI.money_input_text.get()
    userid = GUI.userid_input_text.get()
    appid, vipid = gdata.sql_get_appcardid(appname, cardname)
    # 生成 back_params   回调参数
    back_params = backparams(userid, appid, vipid, pay_set)
    # 生成订单签名
    sign = creat_md5sign(userid, appid, vipid, pay_set)
    # 生成订单
    out_trade_no = create_order(userid, appid, vipid, pay_set, sign)
    # 生成回调签名
    sign_back = create_md5backsign(money, pay_set, out_trade_no)
    # 回调
    back_ldpay(money, pay_set, sign_back, out_trade_no, back_params)


def status_result():
    print(result)
    if result == "success":
        tkinter.messagebox.showinfo('提示', "支付成功!!!")
    else:
        tkinter.messagebox.showerror('错误', "支付异常!!!")
