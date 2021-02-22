import urllib
import urllib.error
import urllib.parse
import urllib.request
import time
import json

# 获取当前时间戳

# http://80.251.211.152:8801/api/login 登录ap
# http://80.251.211.152:8801/api/chargeConfig 金币套餐
# http://192.168.100.105:8801/api/vipChargeConfig vip套餐


url_gold = 'http://80.251.211.152:8801/api/chargeConfig'


def login():
    """登录 获取token"""
    url_login = 'http://80.251.211.152:8801/api/login'

    data = urllib.parse.urlencode({
        "type": '3',
        "appID": '1',
        "equipmentType": '1',
        "machine": 'de519d9e5aeb3017'
    }).encode('utf-8')
    # json.dumps() 将json 字典转化为字符串
    data = bytes(json.dumps({
        "type": '3',
        "appID": '1',
        "equipmentType": '1',
        "machine": 'de519d9e5aeb3017'
    }), 'utf-8')
    headers = {
        'UseCrypt': '0',
        'Content-Type': 'application/json',
        'User-Agent': 'ppp',
        'Conetent-Length': '{}'.format(len(data))
    }
    try:
        request = urllib.request.Request(url_login, data=data, headers=headers, method='POST')
        response = urllib.request.urlopen(request)
        # json.loads() 函数是将json格式数据转换为字典
        result = json.loads(response.read().decode('utf-8'))
        token = result['data']['token']
        print(token)
    except urllib.error.URLError as e:
        print(e.reason)
    return token


def charge(token):
    timestamp = str(time.time()).split('.')
    cc = int(timestamp[0])
    print(type(cc))
    url_vip = 'http://80.251.211.152:8801/api/chargeConfig'
    headers = {
        'Authorization': 'Bearer {}'.format(token),
        'Content-Type': 'application/json',
        'UseCrypt': '0',
        'User-agent': 'xxxx'
    }
    data = bytes(json.dumps({
        "appId": '1',
        "timeStamp": '{}'.format(cc)
    }), 'utf-8')
    request = urllib.request.Request(url_vip, data=data, headers=headers)
    response = urllib.request.urlopen(request)
    result = json.loads(response.read().decode('utf-8'))
    print(result)


if __name__ == '__main__':
    token = login()
    charge(token)
