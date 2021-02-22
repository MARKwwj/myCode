import tkinter
import tkinter.ttk
import tkinter.messagebox
import paymethod as mm
import time
from tkinter import *

current_app_name = ''
current_card_type = 3
app_name_list = ()
LOG_LINE_NUM = 0


class GUI:
    def __init__(self, init_window, gdata):
        self.gdata = gdata
        init_window.title('tool')  # 设置标题
        # init_window.iconbitmap('D:\\python_test\\test_method\\pay_gui\\xxxx.ico')  # 设置图标
        init_window.geometry('324x285')  # 设置大小
        """第一行 app名称"""
        self.app_name_text = tkinter.Label(init_window, text='APP名称：')
        self.app_name_text.place(x=0, y=10, anchor='nw')
        self.cv = tkinter.StringVar()
        self.app_name_com = tkinter.ttk.Combobox(init_window, textvariable=self.cv, width=13, state="readonly")
        self.app_name_com['value'] = self.gdata.get_app_nameid_list()
        self.app_name_com.place(x=60, y=10, anchor='nw')
        self.app_name_com.bind("<<ComboboxSelected>>", self.get_appname)

        """第二行 用户ID"""
        self.user_id_text = tkinter.Label(init_window, text='用户ID：')
        self.user_id_text.place(x=0, y=40, anchor='nw')
        self.cv_id = tkinter.StringVar()
        self.userid_input_text = tkinter.Entry(init_window, textvariable=self.cv_id, width=15)
        self.userid_input_text.place(x=60, y=40, anchor='nw')

        """第三行 点卡类型"""
        self.Point_card_type = tkinter.Label(init_window, text='点卡类型：')
        self.Point_card_type.place(x=0, y=70, anchor='nw')
        cv_2 = tkinter.StringVar()
        self.Point_card_type_com = tkinter.ttk.Combobox(init_window, textvariable=cv_2, width=13, state="readonly")
        self.Point_card_type_com['value'] = ['vip', '金币', '']
        self.Point_card_type_com.current(0)
        self.Point_card_type_com.place(x=60, y=70, anchor='nw')
        self.Point_card_type_com.bind("<<ComboboxSelected>>", self.get_cardtype)

        """第四行 点卡名称"""
        self.Point_card_name = tkinter.Label(init_window, text='点卡名称：')
        self.Point_card_name.place(x=0, y=100, anchor='nw')
        cv_3 = tkinter.StringVar()
        self.Point_card_name_com = tkinter.ttk.Combobox(init_window, textvariable=cv_3, width=13, state="readonly")
        self.Point_card_name_com.place(x=60, y=100, anchor='nw')
        self.Point_card_name_com.bind("<<ComboboxSelected>>", self.set_card_name_money)
        """第五行 金额"""
        self.money_text = tkinter.Label(init_window, text='订单金额：')
        self.money_text.place(x=0, y=130, anchor='nw')
        self.cv_3 = tkinter.StringVar()
        self.money_input_text = tkinter.Entry(init_window, textvariable=self.cv_3, width=15, fg="#969696")
        self.money_input_text.place(x=60, y=130, anchor='nw')
        self.cv_3.set("不建议修改")

        """第六行 生成订单按钮"""
        self.create_order_button = tkinter.Button(init_window, text='立即支付', command=self.button_result)
        self.create_order_button.place(x=60, y=160, anchor='nw')

        """第七行 日志打印"""
        self.print_logs = tkinter.Text(init_window, width=45, height=6)
        self.print_logs.place(x=1, y=195, anchor='nw')

        """换账号"""
        self.tittle_change = tkinter.Label(init_window, text='更换设备用户')
        self.tittle_change.place(x=200, y=10, anchor='nw')

        self.user_id_change = tkinter.Label(init_window, text='用户ID：')
        self.user_id_change.place(x=190, y=40, anchor='nw')
        self.cv_id_change = tkinter.StringVar()
        self.user_id_change_input = tkinter.Entry(init_window, textvariable=self.cv_id_change, width=15)
        self.user_id_change_input.place(x=190, y=40, anchor='nw')

        self.change_id = tkinter.Button(init_window, text='change', command=self.change_button)
        self.change_id.place(x=215, y=70, anchor='nw')

    def get_appname(self, event):
        global current_app_name
        self.Point_card_type_com.current(2)
        self.Point_card_name_com['value'] = ['']
        self.Point_card_name_com.current(0)
        self.cv_3.set("不建议修改")
        # print(self.app_name_com.get())
        current_app_name = self.app_name_com.get()

    def get_cardtype(self, event):
        global current_card_type
        self.Point_card_name_com['value'] = ['']
        self.Point_card_name_com.current(0)
        self.cv_3.set("不建议修改")
        # print(self.Point_card_type_com.get())
        cardtype = self.Point_card_type_com.get()
        if cardtype == '金币':
            current_card_type = 1
            card, actual_money = self.gdata.sql_get_card_name(current_app_name, current_card_type)
            self.Point_card_name_com['value'] = card
        elif cardtype == 'vip':
            current_card_type = 0
            card, actual_money = self.gdata.sql_get_card_name(current_app_name, current_card_type)
            self.Point_card_name_com['value'] = card

    def set_card_name_money(self, event):
        # print(card, money)set_card_name_money()
        print(self.Point_card_name_com.get())
        c, a = self.gdata.sql_get_card_name(current_app_name, current_card_type)
        # print(type(self.Point_card_name_com.get()))
        card_index = c.index(self.Point_card_name_com.get())
        self.cv_3.set(a[card_index])

    def button_result(self):
        userid = self.userid_input_text.get()
        print(len(userid))
        money = self.money_input_text.get()
        appname = self.app_name_com.get()
        judge_switch = True
        try:
            int(money)
        except ValueError:
            tkinter.messagebox.showwarning("警告", "请输入整数金额！")
            judge_switch = False

        try:
            int(userid)
        except ValueError:
            tkinter.messagebox.showwarning("警告", "请输入整数ID！")
            judge_switch = False

        if judge_switch is True:
            if len(userid) != 6:
                tkinter.messagebox.showwarning("警告", "用户ID不正确！")
            elif int(money) <= 1:
                tkinter.messagebox.showwarning("警告", "金额太少拉！")
            elif int(money) >= 1000:
                tkinter.messagebox.showwarning("警告", "金额太大拉！")
            elif len(userid) == 6:
                res = self.gdata.sql_getidexist(userid, appname)
                if len(res) == 0:
                    tkinter.messagebox.showwarning("警告", "用户不存在！")
                else:
                    mm.pay_run(self)
                    mm.status_result()
                    self.print_logs_write()

    def change_button(self):
        userid = self.user_id_change_input.get()
        print(userid)
        judge_switch = True
        try:
            int(userid)
        except ValueError:
            tkinter.messagebox.showwarning("警告", "请输入整数ID！")
            judge_switch = False
        if judge_switch is True:
            if len(userid) != 6:
                tkinter.messagebox.showwarning("警告", "用户ID不正确！")
            elif len(userid) == 6:
                res = self.gdata.change_userid_exist(userid)
                if len(res) == 0:
                    tkinter.messagebox.showwarning("警告", "用户不存在！")
                else:
                    res = self.gdata.change_userid(userid)
                    if res != 0:
                        tkinter.messagebox.showinfo('提示', "更换成功,请重装apk！")
                    else:
                        tkinter.messagebox.showwarning('警告', "更换失败!该账号设备码已被修改过,请重装apk！")

    def print_logs_write(self):
        global LOG_LINE_NUM
        current_time = time.strftime('%H:%M:%S', time.localtime(time.time()))
        appname = self.app_name_com.get()
        cardname = self.Point_card_name_com.get()
        money = self.money_input_text.get()
        userid = self.userid_input_text.get()
        logmsg = appname + " " + userid + " " + cardname + " 金额:" + money
        logmsg_in = str(logmsg) + " 时间 " + str(current_time) + "\n"  # 换行
        if LOG_LINE_NUM <= 5:
            self.print_logs.insert(END, logmsg_in)
            LOG_LINE_NUM = LOG_LINE_NUM + 1
        else:
            self.print_logs.delete(1.0, 2.0)
            self.print_logs.insert(END, logmsg_in)
