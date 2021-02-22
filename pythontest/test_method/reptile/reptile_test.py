import urllib.request
import urllib.parse
import urllib.error
import socket

# urlopen()
# url = 'https://www.python.org'
# response = urllib.request.urlopen(url)
# """urlopen() 的 方法 及 属性练习"""
# print(response.read().decode('utf-8'))

# 查看 响应的数据类型
# print(type(response))
# 返回的数据是HTTPResponse 类型的对象
# 主要包含read（）.readinto（）、getheader（name），getheaders（）、fileno（）等方法
# 以及msg，version，status，reason，debuglevel.closed等属性。
# print(response.status)
# print(response.getheaders())
# print(response.getheader('Content-Type'))
# urlopen()  -data 参数
# 如果传递了 data 参数 那么他的请求方式不再是 get 而是 post
# decode是解码，encode是编码
# 解码代表bytes类型转成str类型
# 编码代表str类型转成bytes类型
# 而bytes类型的数据一般在写入文件时需要用到
# url_02 = 'http://httpbin.org/post'
# data = bytes(urllib.parse.urlencode({'word': 'hello'}), encoding='utf-8')
# response_02 = urllib.request.urlopen(url_02, data=data)
# bytes 转化为 str
# print(str(response_02.read(), 'utf-8'))
# print(response_02.read().decode())

# timeout 参数 超时时间
#
# try:
#     response_03 = urllib.request.urlopen('http://httpbin.org/get', timeout=0.1)
#     # print(response_03.read().decode())
# except socket.timeout as e:
#     print('Time out ')


# data = bytes(urllib.parse.urlencode({'add': 'shanghai'}), encoding='utf-8')
# data = bytes('zhangsan', 'utf-8')
# requset = urllib.request.Request('https://python.org', data=data)
# response = urllib.request.urlopen(requset)
# print(response.read().decode('utf-8'))

# class urllib. request. Request ( url, data=None, headers={}, origin_req_host=None,
# unverifiable=False, method =None)
# url = 'http://80.251.211.152:8801/api/login'
#
# data = bytes(
#     urllib.parse.urlencode({
#         "type": 3,
#         "appID": 1,
#         "equipmentType": 1,
#         "machine": ""
#     }),
#     encoding='utf-8')
# request = urllib.request.Request(url, data=data,
#                                  headers={'UseCrypt': '0', 'contentType': "application/json;charset-UTF-8"})
# response = urllib.request.urlopen(request)
# print(response.read().decode())

# from urllib.error import URLError
# from urllib.request import ProxyHandler, build_opener
#
# proxy_handler = ProxyHandler({
#     'http': '111.206.37.248:80'
# })
# opener = build_opener(proxy_handler)
#
# try:
#     response = opener.open('http://www.baidu.com')
#     print(response.read().decode('utf-8'))
# except URLError as e:
#     print(e.reason)

# import http.cookiejar, urllib.request
#
# cookies = http.cookiejar.CookieJar()
# handler = urllib.request.HTTPCookieProcessor(cookies)
# opener = urllib.request.build_opener(handler)
# request = opener.open()

# from multiprocessing import cpu_count
#
# print(cpu_count())

# import urllib
# import urllib.parse
# import urllib.request
# import urllib.error
#
# url = 'http://www.baidu.com'
# try:
#
#     request = urllib.request.urlopen(url)
# except urllib.error.URLError as e:
#     print(e.reason)


import requests
import json

# files = {'files': open('favicon.ico', 'rb')}
# r = requests.post("http://httpbin.org/post", files=files)
# print(r.text)

r = requests.get("https://www.baidu.com")
print(r.cookies)
for key, value in r.cookies.items():
    print(key + "=" + value)

# data = json.dumps({
#     "type": 3,
#     "appID": 1,
#     "equipmentType": 1,
#     "machine": "2323232"
# })
# headers = {
#     "UseCrypt": "0",
#     "User-Agent": "j",
#     "Content-Type": "application/json"
# }
# r = requests.post('http://80.251.211.152:8801/api/login', data=data, headers=headers)
# print(r.status_code)
# print(r.json())
# print(r.text)
# print(r.cookies)
# print(r.content)

# r = requests.get('https://github.com/favicon.ico')
# with open('favicon.ico', 'wb') as f:
#     f.write(r.content)
