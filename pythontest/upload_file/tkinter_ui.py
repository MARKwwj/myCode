import time
import tkinter as tk
from tkinter import filedialog
import os
import hashlib

from upload_file.mongoDB import *
# from .sftp import *
from upload_file.resource_encrypt import *
from upload_file.sftp import *


# class Application(tk.Frame):
#     def __init__(self, master=None):
#         super().__init__(master)
#         self.master = master
#         self.filePath = tk.StringVar()
#         self.pack()
#         self.create_widgets()
#
#     def create_widgets(self):
#         # 获取文件
#         self.getFile_bt = tk.Button(self)
#         self.getFile_bt['width'] = 15
#         self.getFile_bt['height'] = 1
#         self.getFile_bt["text"] = "打开"
#         self.getFile_bt["command"] = self._getFile
#         self.getFile_bt.pack(side="top")
#
#         # 显示文件路径
#         self.filePath_en = tk.Entry(self, width=30)
#         self.filePath_en.pack(side="top")
#
#         self.filePath_en.delete(0, "end")
#         self.filePath_en.insert(0, "请选择文件")
#
#         # 上传按钮
#         self.upload_bt = tk.Button(self)
#         self.upload_bt['width'] = 15
#         self.upload_bt['height'] = 1
#         self.upload_bt['text'] = '立即上传'
#         self.upload_bt['command'] = self._disposeFile
#         self.upload_bt.pack(side="top")
#
#     # 打开文件并显示路径
#     def _getFile(self):
#         default_dir = r"文件路径"
#         self.filePath = tk.filedialog.askdirectory(title=u'选择文件', initialdir=(os.path.expanduser(default_dir)))
#         print(self.filePath)
#         self.filePath_en.delete(0, "end")
#         self.filePath_en.insert(0, self.filePath)
#
#     def _disposeFile(self):
#         """ 解密文件 """
#         file_path = self.filePath
#         file_path_list = str(file_path).split('/')
#         file_name = file_path_list[len(file_path_list) - 1]
#         new_file_name = (hashlib.md5((file_name + str(round((time.time() * 1000), 1))).encode('utf-8'))).hexdigest()
#         new_file_path = os.path.join(os.path.dirname(file_path) + '/' + new_file_name)
#         print(new_file_path)
#         print(new_file_name)
#         # 解密本地文件 传入空的list接收文件名
#         # file_path 加密前的文件
#         # new_file_path 加密后的文件
#         show_files(file_path, new_file_path, [])
#
#         """ 上传小说到服务器 """
#         # 上传加密的到服务器
#         upload_file(new_file_path, '/wwwroot/novel/uncrypt/chapterDetails/{0}'.format(new_file_name))
#         upload_file(new_file_path, '/wwwroot/novel/crypted/chapterDetails/{0}'.format(new_file_name))
#         """ 存入数据库 """
#         novel_id = cal_id()
#         insert_novel(file_name, novel_id)
#         insert_chapter_detail_mdb(file_path, new_file_name, novel_id)
#

# def run():
#     root = tk.Tk()
#     root.title("Demo")
#     root.geometry("220x150")
#     app = Application(master=root)
#     app.mainloop()



