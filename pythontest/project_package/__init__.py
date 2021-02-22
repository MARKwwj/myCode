import json

import requests


def update_res():
    url = "http://199.180.114.169:8880/prod-api/user/login"

    payload = json.dumps({
        "password": "admin123",
        "username": "istest",
        "code": "608700",
    })
    headers = {
        'Content-Type': 'application/json'
    }
    response = requests.request("POST", url, headers=headers, data=payload)
    result = response.json()
    token = result["data"]["token"]
    print(token)


update_res()
