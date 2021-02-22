import logging

fmt = '[%(asctime)s] %(levelname)s [%(funcName)s: %(filename)s, %(lineno)d] %(message)s'
logging.basicConfig(level=logging.DEBUG,  # log level
                    filename='film_dis.txt',  # log名字
                    format=fmt,  # log格式
                    datefmt='%Y-%m-%d %H:%M:%S',  # 日期格式
                    filemode='a')  # 追加模式
