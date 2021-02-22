import json
import os
import time

import requests

import paramiko
import telegram
from requests_toolbelt import MultipartEncoder

telegram_token = '1352787097:AAF23njz9wpHVxdWHrjHXisWFDTF0oxdMYU'


def get_token():
    # 获取 X-Token
    url_token = "http://199.180.114.169:8880/prod-api/user/login"

    payload_token = json.dumps({
        "password": "admin123",
        "username": "istest",
        "code": "602541",
    })
    headers_token = {
        'Content-Type': 'application/json'
    }
    response_token = requests.request("POST", url_token, headers=headers_token, data=payload_token)
    result_token = response_token.json()
    assert result_token["msg"] == "成功", "X-Token传 失败！！！"
    token = result_token["data"]["token"]
    return token


def upload_zip(token):
    # 上传压缩包
    url_upload = "http://199.180.114.169:8880/prod-api/novel/app/config/package_upload"
    m = MultipartEncoder(
        fields={
            'file': ('uncrypt.zip', open('/Users/mini/Documents/gitauto/magic_app/magic_res/wwwroot/da_xiang/uncrypt/uncrypt.zip', 'rb'),
                     'application/octet-stream')})
    headers_upload = {
        "X-Token": "{0}".format(token),
        'User-Agent': 'test',
        "Content-Type": m.content_type
    }
    re = requests.post(url_upload, data=m, headers=headers_upload).json()
    print(re)
    assert re["msg"] == "成功", "压缩包上传 失败！！！"
    print("压缩包上传成功！！！")


def get_zip_list(token):
    # 获取 压缩包list
    url_zip = "http://199.180.114.169:8880/prod-api/novel/app/config/get_upload_package_list"
    payload_zip = {}
    headers_zip = {
        "X-Token": "{0}".format(token),
    }
    response_zip = requests.request("POST", url_zip, headers=headers_zip, data=payload_zip).json()
    assert response_zip["msg"] == "成功", "压缩包list获取 失败！！！"
    zip_info_dict = response_zip["data"][0]
    return zip_info_dict


def update_collection_to_app(zip_info_dict, token):
    # 同步合辑至app
    url_update = "http://199.180.114.169:8880/prod-api/novel/app/config/update_app_json"
    payload_update = json.dumps(zip_info_dict)
    headers_update = {
        "X-Token": "{0}".format(token),
        'Content-Type': 'application/json'
    }
    response_update = requests.request("POST", url_update, headers=headers_update, data=payload_update).json()
    print(response_update)
    assert response_update["msg"] == "成功", "同步合辑至app 失败！！！"
    print("同步合辑至app成功！！！")


def update_zip_to_app(zip_info_dict, token):
    # 同步zip至app
    url_update2 = "http://199.180.114.169:8880/prod-api/novel/app/config/update_zip"
    payload_update2 = json.dumps(zip_info_dict)
    headers_update2 = {
        "X-Token": "{0}".format(token),
        'Content-Type': 'application/json'
    }
    response_update2 = requests.request("POST", url_update2, headers=headers_update2, data=payload_update2).json()
    print(response_update2)
    assert response_update2["msg"] == "成功", "同步zip至app 失败！！！"
    print("同步zip至app成功！！！")


def upload(local_zip_dir, server_zip_dir):
    local_zip_path = os.path.join(local_zip_dir, 'uncrypt.zip')
    server_zip_path = os.path.join(server_zip_dir, 'uncrypt.zip')
    while True:
        try:
            transport = paramiko.Transport(('199.180.114.169', 26091))
            transport.connect(username='root', password='VzyC3RmT4wRY')
            transport.banner_timeout = 60
            sftp = paramiko.SFTPClient.from_transport(transport)
            ssh = paramiko.SSHClient()
            ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            ssh.connect('199.180.114.169', 26091, 'root', 'VzyC3RmT4wRY', timeout=60)
            print('资源服务器 连接成功！！！')
            break
        except Exception:
            print('资源服务器 连接失败！！！')
            time.sleep(3)
            print('三秒后重连！！！')
            continue

    ssh.exec_command('rm -rf {0}*'.format(server_zip_dir))
    sftp.put(local_zip_path, server_zip_path)
    if sftp.stat(server_zip_path):
        ssh.exec_command('unzip {0} -d {1}'.format(server_zip_path, server_zip_dir))
        ssh.exec_command('zip {0} {1}app.json {2}conf {3}zip'.format(server_zip_path, server_zip_dir, server_zip_dir, server_zip_dir))
        # 更新 zip
        try:
            # 获取token
            token = get_token()
            # 上传zip
            upload_zip(token)
            # 获取最新 zip list
            zip_info_dict = get_zip_list(token)
            # 更新 合辑 到 app
            update_collection_to_app(zip_info_dict, token)
            # 更新 zip 到 app
            update_zip_to_app(zip_info_dict, token)
        except Exception as e:
            print(e)
            send_mes_to_tele()

        sftp.close()
        ssh.close()
    else:
        print("远端 uncrypt.zip 不存在！！！")


def send_mes_to_tele():
    text = "打包异常！@Gaofeiji747 "
    bot = telegram.Bot(telegram_token)
    bot.send_message(-1001387400552, text=text)


def zip(local_zip_dir, server_zip_dir):
    # local_zip_path 本地压缩包路径
    os.chdir(local_zip_dir)
    os.system('rm -rf uncrypt.zip')
    result = os.system('zip -r uncrypt.zip app.json conf zip')
    if result != 0:
        print("本地 uncrypt.zip 压缩失败！！！")
    else:
        print("本地 uncrypt.zip 压缩成功！！！")
        upload(local_zip_dir, server_zip_dir)
    return result


server_zip_dir = '/home/resources/wwwroot/novel/zip/'

local_zip_dir = '/Users/mini/Documents/gitauto/magic_app/magic_res/wwwroot/novel/uncrypt/'

zip(local_zip_dir, server_zip_dir)
