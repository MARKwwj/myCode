# coding:utf-8
# version 20181027 by san
import re, time, os
from urllib import request
from lxml import etree  # python xpath 单独使用导入是这样的
import platform
import ssl

ssl._create_default_https_context = ssl._create_unverified_context  # 取消全局证书


# 爬虫电影之类
class getMovies:
    def __init__(self, url, Thuder):
        '''        实例初始化       '''
        self.url = url
        self.Thuder = Thuder

    def getResponse(self, url):
        url_request = request.Request(self.url)
        url_response = request.urlopen(url_request)
        return url_response  # 返回这个对象

    def newMovie(self):
        ''' 获取最新电影 下载地址与url '''
        http_response = self.getResponse(webUrl)  # 拿到http请求后的上下文对象(HTTPResponse object)
        data = http_response.read().decode('gbk')
        # print(data)  #获取网页内容
        html = etree.HTML(data)
        newMovies = dict()
        lists = html.xpath('/html/body/div[1]/div/div[3]/div[2]/div[2]/div[1]/div/div[2]/div[2]/ul/table//a')
        for k in lists:
            if "app.html" in k.items()[0][1] or "最新电影下载" in k.text:
                continue
            else:
                movieUrl = webUrl + k.items()[0][1]
                movieName = k.text.split('《')[1].split("》")[0]
                newMovies[k.text.split('《')[1].split("》")[0]] = movieUrl = webUrl + k.items()[0][1]
        return newMovies


def Movieurl(self, url):
    ''' 获取评分和更新时间 '''
    url_request = request.Request(url)
    movie_http_response = request.urlopen(url_request)
    data = movie_http_response.read().decode('gbk')
    if len(re.findall(r'豆瓣评分.+?.+users', data)):  # 获取评分;没有评分的返回null
        pingf = re.findall(r'豆瓣评分.+?.+users', data)[0].split('/')[0].replace("\u3000", ":")

    else:
        pingf = "豆瓣评分:null"
    desc = re.findall(r'简\s+介.*', data)[0].replace("\u3000", "").replace('<br />', "").split("src")[0].replace('&ldquo',
                                                                                                               "").replace(
        '&rdquo', "").replace('<img border="0"', "")

    times = re.findall(r'发布时间.*', data)[0].split('\r')[0].strip()  # 获取影片发布时间
    html = etree.HTML(data)
    murl = html.xpath('//*[@id="Zoom"]//a')
    for k in murl:
        for l in k.items():
            if "ftp://" in l[1]:
                print
                return l[1], times, pingf, desc


def check_end(self, fiename, path):
    ''' 检测文件是否下载完成 '''
    return os.path.exists(os.path.join(save_path, fiename))


def check_start(self, filename):
    ''' 检测文件是否开始下载 '''
    cache_file = filename + ".xltd"
    return os.path.exists(os.path.join(save_path, cache_file))


def DownMovies(self, name, url):
    ''' windows下载 '''
    PlatForm = platform.system()
    print("即将下载的电影是：%s" % name)
    if PlatForm == "Windows":
        try:
            print(r'"{0}" "{1}"'.format(self.Thuder, url))
            os.system(r'"{0}" {1}'.format(self.Thuder, url))
        except Exception as e:
            print(e)
    else:
        print("当前系统平台不支持")


def Main(self):
    ''' 最终新电影存储在字典中 '''
    NewMoveis = dict()
    Movies = self.newMovie()  # 获取电影的字典信息
    for k, v in Movies.items():
        NewMoveis[k] = self.Movieurl(v), v
    return NewMoveis


def NewMoives(self):
    ''' 查看已经获取到的电影信息 '''
    Today = time.strftime("%Y-%m-%d")
    print("今天是：%s" % Today)
    Movies = self.Main()
    print("最近的 %s 部新电影：" % len(Movies.keys()))
    for k, v in Movies.items():
        #           print(v[0][1].split(":")[0])
        if Today in v[0][1].split(":")[0]:
            print("++++++++++++++++今天刚更新++++++++++++++：", "\n", k, "-->", v, "\n")
        else:
            print("========================================")
            print(k, "-->", v, "\n")


if __name__ == '__main__':
    # 以下依据您个人电影迅雷的相关信息填写即可
    save_path = "D:\mp4"  # 电影下载保存位置 （需要填写）
    Thuder = "D:\download\软件\新建文件夹 (3)\Thunder\Program\ThunderStart.exe"  # Thuder: 迅雷Thuder.exe路径 （需要填写）
    webUrl = 'http://www.dytt8.net'  # 电影网站
    test = getMovies(webUrl, Thuder)  # 实例化
    test.NewMoives()
    Movies = test.Main()

    for k, v in Movies.items():
        movies_name = v[0][0].split('/')[3]
        socre = v[0][2].split(":")[1]
        check_down_status = test.check_end(movies_name, save_path)
        #       print(check_down_status)
        if check_down_status:
            print("电影: %s 已经下载" % movies_name)
            continue
        elif socre == 'null':
            continue
        elif float(socre) > 7.0:
            print(movies_name, socre)
            test.DownMovies(k, v[0][0])
            time.sleep(10)
