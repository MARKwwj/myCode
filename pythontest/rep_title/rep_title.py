import lxml
import requests
from fake_useragent import UserAgent
from bs4 import BeautifulSoup


def run(page_num):
    us = UserAgent()
    url = 'http://www.w3397.com/movie/list-1-%E5%85%A8%E9%83%A8-%E5%85%A8%E9%83%A8-%E5%85%A8%E9%83%A8-0-{0}.html'.format(page_num)
    header = {'user-agent': us.chrome}
    response = requests.post(url, headers=header)
    html = response.text
    soup = BeautifulSoup(html, "lxml")
    res = soup.find_all('a', {'class': 'poster-pic-name width_quan'})
    with open('title.txt', 'a', 1024 * 1024 * 1, encoding='utf-8') as f:
        for i in res:
            f.write(i.string + '\n')


# for i in range(1, 2257):
#     run(i)

def rep_pho():
    us = UserAgent()
    url = 'https://www.runoob.com/docker/docker-architecture.html'
    header = {'user-agent': us.chrome}
    response = requests.post(url, headers=header).text
    print(response)


rep_pho()
