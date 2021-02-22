import tkinter
from get_data import Gdata
from gui import GUI


def run():
    root = tkinter.Tk()  # 创建窗体
    gdata = Gdata()
    my_gui = GUI(root, gdata)
    root.mainloop()


if __name__ == '__main__':
    run()
