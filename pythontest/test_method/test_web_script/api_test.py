# coding=UTF-8
# This Python file uses the following encoding: utf-8
import json
import re

import requests
import time

from fake_useragent import UserAgent

url_all = [
    "https://add.x69.app/",
    "https://add.routeapi.xyz/",
    "https://every.dx99.app",
    "https://xsmuntp.routeapi.club/",
    "https://zettag.routeapi.club/",
    "https://lan.dx99.app",
    "https://icK.eggrxx.cn",

]


def request_new(url_real):
    count = 0
    global res, exception
    ua = UserAgent()
    while True:
        try:
            headers = {

                'User-Agent': ua.random

            }
            res, exception = "", ""
            response = requests.head(url=url_real, headers=headers)
            code = response.status_code
            if code != 200:
                res = "request failed-code:{0}—{1}".format(response.status_code, url_real)
            # 关闭请求  释放内存
            response.close()
            del response
            break
        except Exception as e:
            print("请求异常！！！")
            exception = str(e)
            if re.search("getaddrinfo failed", exception) is not None:
                res = "request failed —{0}".format(url_real)
            else:
                res, exception = "", ""
            count += 1
            if count == 5:
                res = "request failed —{0}".format(url_real)
                break
            time.sleep(3)
    return res, exception


def automessage(request_res):
    url_post = "https://api.telegram.org/bot1352787097:AAEEMLsuR3jDjA1Kl9i7bVcRu8MNK52LTVU/sendMessage"
    header = {
        "Content-Type": "application/json",
        'Connection': 'close'
    }
    data = json.dumps(
        {
            "chat_id": -435326960,
            "text": "{0}".format(request_res)
        })

    response = requests.post(url_post, headers=header, data=data)
    # 关闭请求  释放内存
    response.close()
    del response


def sendme(execption):
    if execption == "":
        execption = "request failed"
    url_post = "https://api.telegram.org/bot1352787097:AAEEMLsuR3jDjA1Kl9i7bVcRu8MNK52LTVU/sendMessage"
    header = {
        "Content-Type": "application/json",
        'Connection': 'close'
    }
    data = json.dumps(
        {
            "chat_id": -474899309,
            "text": "{0}".format(execption)
        })
    response = requests.post(url_post, headers=header, data=data)
    # 关闭请求  释放内存
    response.close()
    del response


def run():
    while True:
        cur_time = time.asctime(time.localtime(time.time()))
        print("%s 开始检测域名" % cur_time)
        for url_real in url_all:
            request_res, exception = request_new(url_real)
            # request_res非空代表 请求有异常发生
            if request_res != "":
                print(request_res)
                # 发送异常url 到 telegram
                # automessage(request_res)
                # sendme(exception)
            else:
                print("当前url {0}".format(url_real))
        print("暂停检测 60s !")
        time.sleep(60)


if __name__ == "__main__":
    res, exception = "", ""
    run()
