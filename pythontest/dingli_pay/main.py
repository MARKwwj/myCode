import hashlib
import json

import pymysql
import requests
from pymysql.cursors import DictCursor

# long_video_login_url = "http://199.180.114.169:8001/long_video/user/login"

create_order_url = "http://199.180.114.169:8004/pay_service/order/create"

pay_callback_url = "http://199.180.114.169:8004/pay_service/notify/23"


def get_token(info):
    payload = json.dumps({
        "machineCode": "{0}".format(info.machine_code)
    })
    headers = {
        'Content-Type': 'application/json'
    }
    if info.app_name == "long_video":
        login_url = "http://199.180.114.169:8001/long_video/user/login"
    elif info.app_name == "novel":
        login_url = "http://74.120.174.18:9000/novel/user/login"
    response = requests.request("POST", login_url, headers=headers, data=payload)
    print(response.text)
    token = (response.json())["data"]["token"]
    print("token: {0}".format(token))
    return token


def create_order(token, commodity_id, channel_id, app_type):
    payload = json.dumps({
        "channelId": channel_id,
        "commodityId": commodity_id,
        "appType": app_type
    })
    headers = {
        'X-Token': '{0}'.format(token),
        'Content-Type': 'application/json'
    }
    response_json = requests.request("POST", create_order_url, headers=headers, data=payload).json()
    if response_json["msg"] == "成功":
        print(response_json)
        order_number = response_json["data"]["orderNumber"]
        print("order_number: {0}".format(order_number))
        return order_number
    else:
        print(response_json)
        print("订单创建失败！")
        return 0


def pay_callback(out_trade_no, trade_status, money):
    # 签名方式
    # 参数名字典序，最后拼接商户key a=1b=2c=3key
    sign = "money={0}&out_trade_no={1}&pid={2}&trade_no={3}&trade_status={4}NcLQhG52c3MS6W5V" \
        .format(money, out_trade_no, 8888, out_trade_no, trade_status)
    sign_md5 = hashlib.md5(sign.encode(encoding='utf-8')).hexdigest()
    payload = json.dumps({
        "pid": 8888,
        "out_trade_no": out_trade_no,
        "trade_no": out_trade_no,
        "trade_status": trade_status,
        "money": money,
        "sign": sign_md5
    })
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.request("POST", pay_callback_url, headers=headers, data=payload)
    print(response.text)
    return response.text


def pay(info):
    token = get_token(info)
    order_number = create_order(token, info.commodity_id, info.channel_id, info.app_type)
    if order_number == 0:
        return "订单创建异常"
    result = pay_callback(order_number, info.trade_status, info.commodity_amount)
    return result
