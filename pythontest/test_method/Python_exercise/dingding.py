import requests
import os
import json

a = os.listdir("d:\\desktop\\test")
novel = "http://80.251.211.152:8848/apk/novel/{0}".format(a[0])
cartoon = "http://80.251.211.152:8848/apk/cartoon/{0}".format(a[0])
longv = "http://80.251.211.152:8848/apk/longv/{0}".format(a[0])
url = "https://oapi.dingtalk.com/robot/send?access_token=b21316e6fee5e6b1fbebef3a91a5e4f00b22765dd8938dcc0e506bfac5556baa"

header = {
    "Content-Type": "application/json",
    "Charset": "UTF-8"
}
data = json.dumps({"msgtype": "text",
                   "text": {
                       "content": "jenkins自动打包-{0}".format(novel)
                   },
                   "at": {
                       "isAtAll": True  # @全体成员（在此可设置@特定某人）
                   }})

res = requests.post(url, headers=header, data=data)
print(res.status_code)
print(res.text)
