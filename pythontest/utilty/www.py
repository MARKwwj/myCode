import json
import random
import time

import requests

url = "http://192.168.100.113:8002/NovelInfos"

for i in range(1, 100):
    chapters_list = []
    for j in range(1, 21):
        chapters_dict = {
            "subTitle": "章节{0}".format(j),
            "chapterID": j,
            "chapterSize": random.randint(1000, 8000),
            "chapterTime": random.randint(1000, 8000)
        }
        chapters_list.append(chapters_dict)

    data = {
        "onlyQuery": False,
        "title": "标题-{0}".format(i),
        "introduce": "简介简介简介简介简介简介简介",
        "author": "作者",
        "categoryName": "有声",
        "type": "都市,言情,穿越,古代",
        "totalNum": 2,
        "novelStatus": 1,
        "totalWordNum": 0,
        "novelType": 2,
        "totalSize": random.randint(1000, 8000),
        "chapters": chapters_list
    }
    headers = {
        'Content-Type': 'application/json'
    }
    json_data = json.dumps(data)
    response = requests.post(url, data=json_data, headers=headers)
    print(response.json())
