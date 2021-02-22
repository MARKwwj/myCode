import requests

url = "http://80.251.211.152:8801/api/login"

payload = "{\r\n    \"type\":3,\r\n    \"appID\":1,\r\n    \"equipmentType\":1,\r\n    \"machine\":\"xxxxxxx1\",\r\n    \"resUseCrypt\": \"False\",\r\n    \"appCurrentVersion\": \"1.12.0\"," \
          "\r\n    \"timeStamp\": 1597974774000\r\n} \r\n "

headers = {
    'UseCrypt': '0',
    'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload).json()

print(type(response))
print(response)
assert response == {'msg': 'Service exception', 'code': 500}
