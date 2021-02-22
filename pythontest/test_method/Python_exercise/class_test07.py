# 定制类
# class Student:
#     def __init__(self, name):
#         self.name = name
#
#     def __str__(self):
#         return 'Student object (name:%s)' % self.name
#
#
# print(Student('jun'))
# s = Student('zhang')

# 错误处理
# try:
#     print('try...')
#     r = 10 / 2  # int('a')
#     print('result:', r)
# except ZeroDivisionError as e:
#     print('except:', e)
# except ValueError as e:
#     print('except：', e)
# else:
#     print('no error!')
# finally:
#     print('finally...')
# print('end')
#
#
# import logging
# def cal(s):
#     return 10 / int(s)
#
#
# def bar(s):
#     return cal(s) * 2
#
#
# def main():
#     try:
#         bar('0')
#     except Exception as e:
#         print('error:', e)
#         # logging.exception(e)
#     finally:
#         print('finally...')
#
#
# main()
# # Python 内置的 logging 模块可以非常容易地记录错误信息
# print('111')

def foo(s):
    n = int(s)
    assert n > 0, 'n要大于0'
    return print(10 / n)


def main():
    foo('0')


main()
