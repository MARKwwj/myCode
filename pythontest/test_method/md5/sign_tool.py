import tkinter
import os
import sys
import time
from tkinter import *
from __init__ import *

sign_back = ''
out_trade_no = ''
token = ''


class MY_GUI:

    def __init__(self, init_window):
        init_window.title('Pay_tool')  # 设置标题
        init_window.iconbitmap('D:\\python_test\\test_method\\md5\\image\\Pay_tool_ico.ico')  # 设置图标
        init_window.geometry('350x250')  # 设置大小
        # root['background'] = 'black'  # 设置背景颜色
        self.label_text = tkinter.Label(init_window, text='user_id')
        self.label_text.pack()
        self.text = tkinter.Text(init_window, width=15, height=1)
        self.text.pack()
        self.button1 = tkinter.Button(init_window, text="start", command=self.get_userid)  # 创建按钮1
        self.button1.pack()  # 将按钮1添加到窗口中
        self.label_text2 = tkinter.Label(init_window, text='back_params')
        self.label_text2.pack()
        self.text2 = tkinter.Text(init_window, width=40, height=1)
        self.text2.pack()
        self.button2 = tkinter.Button(init_window, text="continue", command=self.get_back_params)  # 创建按钮2
        self.button2.pack()  # 将按钮2添加到窗口中

        self.label_text3 = tkinter.Label(init_window, text='pay_log')  # 日志框
        self.label_text3.pack()

        self.log_data_Text = Text(init_window, width=66, height=9)  # 日志框
        self.log_data_Text.pack()

    def get_userid(self):
        num = self.text.get(0.0, END)
        user_id = num.strip('\n')
        print(user_id)
        pay_set = Settings()
        # 生成订单签名
        sign = creat_md5sign(pay_set, user_id)
        # 生成订单
        global out_trade_no
        out_trade_no = create_order(pay_set, sign, user_id)
        # 生成回调签名
        global sign_back
        sign_back = create_md5backsign(pay_set, out_trade_no)

    def get_back_params(self):
        # 获取 回调参数
        pay_set = Settings()
        back_params = self.text2.get(0.0, END)
        # 回调
        result = back_ldpay(pay_set, sign_back, out_trade_no, back_params.strip('\n'))
        self.write_log(result)

    def get_current_time(self):
        current_time = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(time.time()))
        return current_time

    def write_log(self, msg):
        current_time = self.get_current_time()
        self.log_data_Text.insert(END, msg + ' ' + current_time + '\n')


def run():
    root = tkinter.Tk()  # 创建窗体
    my_gui = MY_GUI(root)
    root.mainloop()  # 界面循环及时显示界面变化


if __name__ == '__main__':
    run()
