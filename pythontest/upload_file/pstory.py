import os

import requests
from bs4 import BeautifulSoup


# 定义请求函数
def get_html(url, headers):
    response = requests.get(url, headers)
    return response.text


# 定义解析函数
def get_novel(html):
    soup = BeautifulSoup(html, 'lxml')
    novel_name = soup.find('h3', {'class': 'j_chapterName'}).getText()
    novel_txt = soup.find('div', {'class': 'read-content j_readContent'}).getText().replace("　　", "\n")
    return novel_name + '\n' + novel_txt


# 定义获取下一章节的url
def get_next_url(html):
    soup = BeautifulSoup(html, 'lxml')
    next_url = soup.find('a', {'id': 'j_chapterNext'})
    return next_url['href']


def save(novel, i):
    save_path = 'D:\\desktop\\全职高手\\'
    if not os.path.exists(save_path):
        os.mkdir(save_path)
    save_name = str(i) + '.txt'
    full_path = save_path + save_name
    fp = open(full_path, 'w')
    fp.write(novel)
    fp.close()


# 定义主函数
def main():
    url = "https://read.qidian.com/chapter/TtxVU3dYVW81/eSlFKP1Chzg1"
    headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64; rv:61.0) Gecko/20100101 Firefox/61.0'}
    i = 1
    while i:
        # 调用请求函数
        html = get_html(url=url, headers=headers)
        # 调用获取下一章url
        url = 'https:' + get_next_url(html=html)
        # 调用解析函数d
        novel = get_novel(html=html)
        # 调用存储函数
        save(novel=novel, i=i)
        i += 1


main()
