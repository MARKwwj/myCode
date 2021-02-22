import requests
import json
import uuid
import random
import time


def binduser(agent_id):
    """
    :param agent_id:  代理账号的列表
    :return:
    """
    login_url = 'http://192.168.100.113:8801/api/login'
    random_machine = str(uuid.uuid1())
    timestamp1 = int(round(time.time() * 1000))
    headers = {
        "UseCrypt": "0",
        "User-Agent": 'jun/test',
        "Content-Type": "application/json"
    }
    login_data = {
        "type": "3",
        "appID": "1",  # app 类型
        "equipmentType": "1",  # 设备类型
        "machine": "{0}".format(random_machine),
        "resUseCrypt": "False",
        "appCurrentVersion": "1.12.0",
        "timeStamp": "{0}".format(timestamp1)
    }
    login_data_json = json.dumps(login_data)
    try:
        res = requests.post(url=login_url, data=login_data_json, headers=headers)
        response = res.json()
        token = response["data"]["token"]
        print(res.status_code)
        print(random_machine)
        print(response["data"]["token"])
    except Exception:
        print("请求异常 10s 后重试！！！")
        time.sleep(10)
        print("尝试再次请求！！！")
        res = requests.post(url=login_url, data=login_data_json, headers=headers)
        response = res.json()
        token = response["data"]["token"]
        print(res.status_code)
        print(random_machine)
        print(response["data"]["token"])

    bind_url = 'http://80.251.211.152:8801/api/bindCode'
    code_list = agent_id  # 传入一个 代理账号列表
    select_bind_code = random.sample(code_list, 1)
    bind_code = select_bind_code[0]
    # 转为16进制
    code = hex(int(bind_code)).split("0x")[1]
    timestamp = int(round(time.time() * 1000))

    bind_headers = {
        "UseCrypt": "0",
        "User-Agent": 'jun/test',
        "Content-Type": "application/json",
        "Authorization": "Bearer {0}".format(token)
    }
    bind_data = {
        "code": "{0}".format(code),
        "type": "0",
        "timeStamp": "{0}".format(timestamp)
    }
    bind_data_json = json.dumps(bind_data)
    try:
        bind_res = requests.post(url=bind_url, data=bind_data_json, headers=bind_headers)
        print(bind_res.json())
        print(bind_res.text)
    except Exception:
        print("请求异常 10s 后重试！！！")
        time.sleep(10)
        print("尝试再次请求！！！")
        bind_res = requests.post(url=bind_url, data=bind_data_json, headers=bind_headers)
        print(bind_res.json())
        print(bind_res.text)


if __name__ == '__main__':
    # agent_id = ["179416", "179417", "179419", "184272", "184312", "184405"]
    agent_id = ["179416"]
    # for i in range(4):
    binduser(agent_id)
