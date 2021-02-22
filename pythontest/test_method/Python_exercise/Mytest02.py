import requests
import json
r = requests.get('https://www.baidu.com')
code = r.status_code
html = str(r.content, 'utf-8')
raw_this = r.raw

print('code:{},\ntext:{}\nraw:{}'.format(code, html, raw_this))
