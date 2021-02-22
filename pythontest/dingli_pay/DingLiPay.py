import time
import tkinter
import tkinter.messagebox
import tkinter.ttk
from tkinter import END

from get_datas import *
from parameters import *
from main import pay

LOG_LINE_NUM = 0


class GUI:
    def __init__(self, window, info):
        window.title('tools')  # 设置标题
        window.geometry('400x380')  # 设置大小
        """第一行 app名称"""
        self.app_name_text = tkinter.Label(window, text='APP名称：')
        self.app_name_text.place(x=0, y=10, anchor='nw')
        self.app_name_com = tkinter.ttk.Combobox(window, width=13, state="readonly")
        self.app_name_com['value'] = info.app_names_list
        self.app_name_com.place(x=60, y=10, anchor='nw')
        self.app_name_com.bind("<<ComboboxSelected>>", self.get_appname)

        """第二行 用户ID"""
        self.user_id_text = tkinter.Label(window, text='用户ID：')
        self.user_id_text.place(x=0, y=40, anchor='nw')
        self.user_id_input_text = tkinter.Entry(window, width=15)
        self.user_id_input_text.place(x=60, y=40, anchor='nw')
        """第三行 套餐类型"""
        self.Point_card_type = tkinter.Label(window, text='套餐类型：')
        self.Point_card_type.place(x=0, y=70, anchor='nw')
        self.Point_card_type_com = tkinter.ttk.Combobox(window, width=13, state="readonly")
        self.Point_card_type_com['value'] = ['vip', '金币']
        self.Point_card_type_com.place(x=60, y=70, anchor='nw')
        self.Point_card_type_com.bind("<<ComboboxSelected>>", self.get_commodity_type)
        """第四行 套餐名称"""
        self.Point_card_name = tkinter.Label(window, text='套餐名称：')
        self.Point_card_name.place(x=0, y=100, anchor='nw')
        self.Point_card_name_com = tkinter.ttk.Combobox(window, width=13, state="readonly")
        self.Point_card_name_com.place(x=60, y=100, anchor='nw')
        self.Point_card_name_com.bind("<<ComboboxSelected>>", self.get_commodity_name)

        """第五行 金额"""
        self.money_text = tkinter.Label(window, text='订单金额：')
        self.money_text.place(x=0, y=130, anchor='nw')
        self.money_input_text_value = tkinter.StringVar()
        self.money_input_text = tkinter.Entry(window, width=15, textvariable=self.money_input_text_value, fg="#969696")
        self.money_input_text.place(x=60, y=130, anchor='nw')
        """第六行 支付方式"""
        self.pay_type_name = tkinter.Label(window, text='支付方式：')
        self.pay_type_name.place(x=0, y=160, anchor='nw')
        self.pay_type_name_com = tkinter.ttk.Combobox(window, width=13, state="readonly")
        self.pay_type_name_com.place(x=60, y=160, anchor='nw')
        self.pay_type_name_com['value'] = ['支付宝', '微信']
        self.pay_type_name_com.bind("<<ComboboxSelected>>", self.get_pay_type)
        """第七行 支付状态"""
        self.pay_state_name = tkinter.Label(window, text='支付状态：')
        self.pay_state_name.place(x=0, y=190, anchor='nw')
        self.pay_state_name_com = tkinter.ttk.Combobox(window, width=13, state="readonly")
        self.pay_state_name_com.place(x=60, y=190, anchor='nw')
        self.pay_state_name_com['value'] = ['成功', '失败']
        self.pay_state_name_com.bind("<<ComboboxSelected>>", self.get_pay_state)
        """第六行 生成订单按钮"""
        self.create_order_button = tkinter.Button(window, text='立即支付', command=self.pay_button)
        self.create_order_button.place(x=60, y=220, anchor='nw')

        """第七行 日志打印"""
        self.print_logs = tkinter.Text(window, width=100, height=10)
        self.print_logs.place(x=1, y=250, anchor='nw')

        """换账号"""
        # self.tittle_change = tkinter.Label(window, text='更换设备用户')
        # self.tittle_change.place(x=240, y=10, anchor='nw')
        #
        # self.user_id_change = tkinter.Label(window, text='用户ID：')
        # self.user_id_change.place(x=180, y=40, anchor='nw')
        # self.user_id_change_input = tkinter.Entry(window, width=15)
        # self.user_id_change_input.place(x=230, y=40, anchor='nw')
        self.change_id = tkinter.Button(window, text='更换用户', command=self.change_button)
        self.change_id.place(x=130, y=220, anchor='nw')

    def input_id_check(self, id):
        if id != "":
            try:
                int(id)
            except Exception:
                tkinter.messagebox.showwarning("警告", "请输入整数id！")
                return
            if len(id) > 9:
                tkinter.messagebox.showwarning("警告", "用户id不符合规则！")
                return
        else:
            tkinter.messagebox.showwarning("警告", "请输入用户id！")
            return

    def change_button(self):
        self.input_id_check(self.user_id_input_text.get())
        info.need_changed_user_id = self.user_id_input_text.get()
        result = change_user(info)
        if result == 1:
            tkinter.messagebox.showinfo('提示', "更换成功！\n请重新安装APK！")
        elif result == "is_changed":
            tkinter.messagebox.showinfo('提示', "此账户已被更换！\n请重新安装APK！")
        else:
            print(result)
            tkinter.messagebox.showwarning('提示', "更换异常！")

    def get_appname(self, event):
        info.app_name = self.app_name_com.get()
        if info.app_name == "long_video":
            info.app_type = 1
            info.database = "long_video_db"
        elif info.app_name == "novel":
            info.app_type = 4
            info.database = "novel_db"
        # 获取视频套餐信息
        select_commodity(info)

    def get_commodity_type(self, event):
        # 充值类型 1vip  2金币
        # 获取套餐类型
        self.Point_card_name_com['value'] = ['']
        self.Point_card_name_com.current(0)
        self.money_input_text.delete(0, 'end')
        if self.Point_card_type_com.get() == "vip":
            info.commodity_type = 1
        elif self.Point_card_type_com.get() == "金币":
            info.commodity_type = 2
        # print("commodity_type:", info.commodity_type)
        # 套餐名称 赋值
        self.set_commodity_name()

    def set_commodity_name(self):
        commodity_info = []
        if info.commodity_type == 1:
            commodity_info = info.commodity_vip_info
        elif info.commodity_type == 2:
            commodity_info = info.commodity_gold_info
        info.card_name_list = []
        for name in commodity_info:
            info.card_name_list.append(name["commodity_name"])
        self.Point_card_name_com["value"] = info.card_name_list

    def get_commodity_name(self, event):
        commodity_name = self.Point_card_name_com.get()
        index = info.card_name_list.index(commodity_name)
        if info.commodity_type == 1:
            info.commodity_id = info.commodity_vip_info[index]["commodity_id"]
            info.commodity_amount = info.commodity_vip_info[index]["commodity_in_fact_amount"]
        elif info.commodity_type == 2:
            info.commodity_id = info.commodity_gold_info[index]["commodity_id"]
            info.commodity_amount = info.commodity_gold_info[index]["commodity_in_fact_amount"]
        # print(info.commodity_amount)
        # 设置 套餐金额
        self.money_input_text_value.set(info.commodity_amount)

    def get_pay_type(self, event):
        if self.pay_type_name_com.get() == "支付宝":
            info.channel_id = 53  # 53 支付宝
        elif self.pay_type_name_com.get() == "微信":
            info.channel_id = 54  # 54 微信

    def get_pay_state(self, event):
        if self.pay_state_name_com.get() == "成功":
            info.trade_status = "TRADE_SUCCESS"
        elif self.pay_state_name_com.get() == "失败":
            info.trade_status = "Failed"

    def pay_button(self):
        self.input_id_check(self.user_id_input_text.get())
        info.user_id = self.user_id_input_text.get()
        if info.app_type == 0:
            tkinter.messagebox.showwarning("警告", "请选择app名称！")
            return
        if info.commodity_type == 0:
            tkinter.messagebox.showwarning("警告", "请选择套餐类型！")
            return
        if self.Point_card_name_com.get() == "":
            tkinter.messagebox.showwarning("警告", "请选择套餐名称！")
            return
        if info.channel_id == 0:
            tkinter.messagebox.showwarning("警告", "请选择支付方式！")
            return
        if info.trade_status == "":
            tkinter.messagebox.showwarning("警告", "请选择支付状态！")
            return
        # 查询用户是否存在
        select_user_id_whether_exists(info)
        if info.user_whether_exists == 0:
            tkinter.messagebox.showwarning("警告", "用户不存在！")
            return
        print("**********************************")
        print("app_type:", info.app_type)
        print("user_id:", info.user_id)
        print("commodity_type:", info.commodity_type)
        print("commodity_id:", info.commodity_id)
        print("commodity_amount:", info.commodity_amount)
        print("channel_id:", info.channel_id)
        print("trade_status:", info.trade_status)
        print("**********************************")
        # 支付
        get_machine_code(info)
        result = pay(info)
        print(result)
        # 打印日志
        self.print_logs_write()
        if result == "success":
            tkinter.messagebox.showinfo('提示', "订单创建成功\n订单状态:已支付")
        elif result == "fail":
            tkinter.messagebox.showinfo('提示', "订单创建成功\n订单状态:未支付")
        elif result == "订单创建异常":
            tkinter.messagebox.showinfo('提示', "订单创建异常!!!")
        else:
            tkinter.messagebox.showerror('错误', "支付异常!!!")

    def print_logs_write(self):
        global LOG_LINE_NUM
        current_time = time.strftime('%H:%M:%S', time.localtime(time.time()))
        app_name = self.app_name_com.get()
        card_name = self.Point_card_name_com.get()
        money = self.money_input_text.get()
        user_id = self.user_id_input_text.get()
        pay_type = self.pay_type_name_com.get()
        pay_state = self.pay_state_name_com.get()
        logmsg = app_name + " " + user_id + " " + card_name + " 金额:" + money + " " + pay_type + " " + pay_state
        logmsg_in = str(logmsg) + " 时间 " + str(current_time) + "\n"  # 换行
        # with open("./record.txt", 'a') as fw:
        #     fw.write(logmsg_in)
        if LOG_LINE_NUM <= 8:
            self.print_logs.insert(END, logmsg_in)
            LOG_LINE_NUM = LOG_LINE_NUM + 1
        else:
            self.print_logs.delete(1.0, 2.0)
            self.print_logs.insert(END, logmsg_in)


#  初始化窗口
window = tkinter.Tk()
info = All_infos()
GUI(window, info)

# 主窗口循环显示
window.mainloop()
