import requests
from fake_useragent import UserAgent
from lxml import etree


class wallhaven(object):
    def __init__(self):
        self.url = "https://wallhaven.cc/search?q=id%3A65348&page={}"
        ua = UserAgent(verify_ssl=False)
        for i in range(1, 50):
            self.headers = {
                'User-Agent': ua.random
            }

    def main(self):
        startPage = int(input("起始页:"))
        endPage = int(input("终止页:"))
        for page in range(startPage, endPage + 1):
            url = self.url.format(page)

    def get_page(self, url):
        res = requests.get(url=url, headers=self.headers)
        html = res.content.decode("utf-8")
        return html

    def parse_page(self, html):
        parse_html = etree.HTML(html)
        image_src_list = parse_html.xpath('//figure//a/@href')


if __name__ == '__main__':
    imageSpider = wallhaven()
    imageSpider.main()
