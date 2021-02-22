#!/usr/bin/python2.7
from types import MethodType


# 创建一个类
class Student(object):
    pass


# 创建一个函数
def set_age(self, age):
    self.age = age


def print_hello(self):
    print('hello')


s1 = Student()
# 给一个实例绑定的方法，对另一个实例是不起作用的
s1.set_age = MethodType(set_age, s1)
s1.print_hello = MethodType(print_hello, s1)  # 第一个参数是要绑定的方法，第二个参数是要绑定的对象，第三个参数是类名（可省略）
s1.set_age(23)
print(s1.age)
s1.print_hello()

s2 = Student()
# s2.set_age(24)
# print(s2.age)
# s2.print_hello()

# 为了给所有实例都绑定方法，可以给 class 绑定方法
Student.set_age = set_age
# 给 class 绑定方法后，所有实例均可调用
s2.set_age(24)
print(s2.age)
