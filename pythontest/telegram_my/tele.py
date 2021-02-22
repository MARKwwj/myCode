import json
import os
import sys

import requests

import telegram

token = '1445191367:AAHCq-M_2cOCzx4GKFSDgHDGHnvrlx4akds'
url_post_dingding = "https://oapi.dingtalk.com/robot/send?access_token=b21316e6fee5e6b1fbebef3a91a5e4f00b22765dd8938dcc0e506bfac5556baa"


def automessage(type):
    global path
    global url
    global name
    # 1 novel 2 cartoon 3 longv
    if type == "1":
        path = os.listdir("/users/mini/.jenkins/workspace/novel")
        url = "http://199.180.114.169:8848/apk/novel/{0}".format(path[0])
        name = "小说app"
    elif type == "2":
        path = os.listdir("/users/mini/.jenkins/workspace/cartoon")
        url = "http://199.180.114.169:8848/apk/cartoon/{0}".format(path[0])
        name = "漫画app"
    elif type == "3":
        path = os.listdir("/users/mini/.jenkins/workspace/da_xiang")
        url = "http://199.180.114.169:8848/apk/longv/{0}".format(path[0])
        name = "长视频app"
    elif type == "4":
        path = os.listdir("/users/mini/.jenkins/workspace/new_long")
        url = "http://199.180.114.169:8848/apk/longv/{0}".format(path[0])
        name = "新长视频app"

    header = {
        "Content-Type": "application/json",
        "Charset": "UTF-8"
    }
    data_dingding = json.dumps({"msgtype": "text",
                                "text": {
                                    "content": "-{0} {1}".format(name, url)
                                },
                                })

    res_dingding = requests.post(url_post_dingding, headers=header, data=data_dingding)
    print(res_dingding.json())
    # 发送消息到飞机
    text = "-{0} {1}".format(name, url)

    bot = telegram.Bot(token)
    bot.send_message(-1001301978742, text=text)


if __name__ == "__main__":
    a = sys.argv[1]
    automessage(a)


