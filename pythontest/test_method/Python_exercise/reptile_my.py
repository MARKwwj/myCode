import requests
from bs4 import BeautifulSoup

url = "https://www.qu.la/book/12763/10664294.html"
req = requests.get(url)
req.encoding = 'utf-8'
soup = BeautifulSoup(req.text, 'html.parser')
content = soup.find(id='content')
with open('D:\\desktop\\a.txt', 'w') as f:
    # string = content.text.replace('\u3000', '').replace('\t', '').replace('\n', '').replace('\r', '').replace('『', '“') \
    #     .replace('』', '”')  # 去除不相关字符
    string = string.split('\xa0')  # 编码问题解决
    string = list(filter(lambda x: x, string))
    for i in range(len(string)):
        string[i] = '    ' + string[i]
        if "本站重要通知" in string[i]:  # 去除文末尾注
            t = string[i].index('本站重要通知')
            string[i] = string[i][:t]
    string = '\n'.join(string)
    print(string)
    f.write(string)
