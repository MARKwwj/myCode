import threading
import time

exit_flag = 0


# class myThread(threading.Thread):  # 继承父类的方法
#     def __init__(self, threadID, name, counter):
#         threading.Thread.__init__(self)  # 继承父类的属性
#         self.threadID = threadID
#         self.name = name
#         self.counter = counter
#
#     def run(self):
#         print("开始线程：" + self.name)
#         print_time(self.name, self.counter, 5)
#         print("退出线程：" + self.name)

def run(a):
    s = a * a
    print(s)


if __name__ == '__main__':
    # 创建新线程
    thread1 = threading.Thread(target=run, args=(100,))
    thread2 = threading.Thread(target=run, args=(50,))
    # 开启新线程
    thread1.start()
    thread2.start()
    thread1.join()
    thread2.join()
