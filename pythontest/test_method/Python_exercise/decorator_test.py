# 创建装饰器，要求如下：
#  1. 创建add_log装饰器，被装饰的函数打印日志信息；
#  2. 日志格式为: [字符串时间] 函数名: xxx， 运行时间：xxx, 运行返回值结果:xxx
import time


def add_log(func):
    def wrapper(*args, **kwargs):
        print("[{0}] 函数名：{1}".format(time.asctime(time.localtime(time.time())), func.__name__))
        return func(*args, **kwargs)

    return wrapper


@add_log
def print_test():
    print("2020")


print_test()
