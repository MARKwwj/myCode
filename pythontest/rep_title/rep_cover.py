import hashlib
import os
import re
import wget

import requests
from fake_useragent import UserAgent
from bs4 import BeautifulSoup

url_map = {}


def run():
    us = UserAgent()
    url = 'https://www.mzitu.com/page/2/'
    header = {'user-agent': us.chrome}
    response = requests.post(url, headers=header)
    html = response.text
    soup = BeautifulSoup(html, "lxml")
    res = soup.find_all('img', {'class': 'lazy'})
    for i in res:
        reExpression = re.compile("data-original=\"(https://tbweb.iimzt.com/thumbs/[0-9]{4}/[0-9]{2}/[0-9]{6}_450.jpg)\"")
        result = re.findall(reExpression, str(i))
        map_key = str((hashlib.md5((result[0]).encode(encoding='UTF-8')).hexdigest()))
        if not url_map.__contains__(map_key):
            url_map[map_key] = result[0]



def changeWebp():
    print(os.getcwd())
    cur_path = os.getcwd()
    for i in os.listdir(cur_path):
        cur_file_path = os.path.join(cur_path, i)
        if (i.split("."))[1] != "py":
            cover_name = (i.split(".jpg"))[0]
            cover_name_file_path = os.path.join(cur_path, "webpDir", cover_name)
            command = "ffmpeg -i {0} {1}.webp".format(cur_file_path, cover_name_file_path)
            os.system(command)


if __name__ == '__main__':
    for i in range(2, 5):
        run()
    for k, v in url_map.items():
        wget.download(v)
    changeWebp()
