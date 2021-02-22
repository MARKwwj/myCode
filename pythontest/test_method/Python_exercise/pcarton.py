import os
import time
import requests
from selenium import webdriver
from lxml import etree


def getChapterUrl(url):
    headers = {
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
    }
    part_url = "http://ac.qq.com"
    res = requests.get(url, headers=headers)
    html = res.content.decode()
    el = etree.HTML(html)
    li_list = el.xpath('//*[@id="chapter"]/div[2]/ol[1]/li')
    for li in li_list:
        for p in li.xpath("./p"):
            for span in p.xpath("./span[@class='works-chapter-item']"):
                item = {}
                list_title = span.xpath("./a/@title")[0].replace(' ', '').split('：')
                if list_title[1].startswith(('第', '序')):
                    getChapterFile(part_url + span.xpath("./a/@href")[0], list_title[0], list_title[1])


def getChapterFile(url, path1, path2):
    # path = os.path.join(path)
    # 漫画名称目录
    path = os.path.join(path1)
    if not os.path.exists(path):
        os.mkdir(path)
    # 章节目录
    path = path + '\\' + path2
    if not os.path.exists(path):
        os.mkdir(path)
    chrome = webdriver.Chrome()
    # "http://ac.qq.com/ComicView/index/id/505435/cid/2"
    chrome.get(url)
    time.sleep(4)
    imgs = chrome.find_elements_by_xpath("//div[@id='mainView']/ul[@id='comicContain']//img")
    for i in range(0, len(imgs)):
        js = "document.getElementById('mainView').scrollTop=" + str((i) * 1280)
        chrome.execute_script(js)
        time.sleep(3)
        print(imgs[i].get_attribute("src"))
        with open(path + '\\' + str(i) + '.png', 'wb') as f:
            f.write(requests.get(imgs[i].get_attribute("src")).content)
    chrome.close()
    print('下载完成')


if __name__ == '__main__':
    getChapterUrl('https://ac.qq.com/Comic/comicInfo/id/531490')
