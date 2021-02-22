# 编码
# # ord 获取字符的整数表示 十进制 ascii
# print(ord('a'))
# print(ord('A'))
# # chr()函数把编码转换为对应的字符
# print(chr(97))
# print(chr(65))
# # 以Unicode表示的str通过encode()方法可以编码为指定的bytes
# print('ABC'.encode('ascii'))
# print('中文'.encode('utf-8'))
# print(b'\xe4\xb8\xad\xe6\x96\x87'.decode('utf-8'))
# # 如果bytes中只有一小部分无效的字节，可以传入errors='ignore'忽略错误的字节：
# print(b'\xe4\xb8\xad\xe6\x96ps0d'.decode('utf-8', errors="ignore"))
# # 要计算str包含多少个字符，可以用len()函数：
# print(len(b'\xe4\xb8\xad\xe6\x96\x87'))
# print(len('中文'))
#
# 格式化
# 占位符	替换内容
# %d	    整数
# %f	    浮点数
# %s	    字符串
# %x	    十六进制整数
# print('Hello, %s' % 'world')
# print('Hello, %d' % 100)
# print('Hello, %f' % 0.1)
# # 格式化整数和浮点数还可以指定是否补0和整数与小数的位数(有大于一个占位符%? 后面要带括号)
# print('%2d-%02d' % (3, 1))
# print('%.5f' % 3.1415926)
# # 有些时候，字符串里面的%是一个普通字符怎么办？这个时候就需要转义，用%%来表示一个%：
# print('这是一个百分号 %d %%' % 100)
#
# # list 有序集合
# lis_name = ['赵六', '李四', '张三']
# # 集合的元素个数
# # print(len(lis_name))
# # print(lis_name[0])
# # print(lis_name[len(lis_name)-1])
# # print('________________________')
# # # 为集合增加元素
# lis_name.append('小二')
# # -1 做索引 直接获取集合最后一个元素
# # 以此类推，可以获取倒数第2个、倒数第3个：
# # print(lis_name[-1])
# # print(lis_name[-2])
# # print(lis_name[-3])
# # 也可以把元素插入到指定的位置，比如索引号为1的位置
# lis_name.insert(1, '小三')
# # 打印输出 集合元素
# print('lis_name 集合所有元素')
# for n in range(0, 4):
#     print(lis_name[n])
# print('__________________________________')
# # 要删除list末尾的元素，用pop()方法
# lis_name.pop()
# # 要删除指定位置的元素，用pop(i)方法，其中i是索引位置：
# lis_name.pop(1)
# # 要把某个元素替换成别的元素，可以直接赋值给对应的索引位置：
# lis_name[2] = 'ss'
# # 打印输出集合 所有元素
# print('修改后 lis_name 集合的元素')
# for n in range(0, len(lis_name)):
#     print(lis_name[n])
# print('__________________________________')

# # list里面的元素的数据类型也可以不同，比如：
# lis_dif = [100, '肖肖', 0.2, True]
# # list元素也可以是另一个list，比如：
# s = [1.2, 100, [1, '帅帅', True], '张三']
# print('s 集合的长度 %d' % len(s))
# for n in range(0, len(s)):
#     print(s[n])
#     print(s[2][n])


# tuple
# 另一种有序列表叫元组：tuple。tuple和list非常类似
# 但是tuple一旦初始化就不能修改
# 可以存入不同类型的数据
# tup_01 = ('hello', 'hi', '哦哦', True, 100, 0.2)
# for i in range(0, len(tup_01)):
#     print(tup_01[i])
# tuple 元组 只有一个值时，也会加一个逗号,以免你误解成数学计算意义上的括号：
# tup_02 = (1,)
# print(tup_02[0])

# if判断条件还可以简写，比如写：
# 只要x是非零数值、非空字符串、非空list等，就判断为True，否则为False
# x = 1
# if x:
#     print('true')
# else:
#     print('false')

# input  读取用户输入
# input_name = input("请输入姓名:")
# if input_name == "张三":
#     print("hello，张三")
# elif input_name == "李四":
#     print("hello,李四")
# else:
#     print("hello!!!")
# s = input("birth:")
# birth = int(s)
# if birth < 2000:
#     print("你是00前")
# else:
#     print("你是00后")


# 练习
# 小明身高1.75，体重80.5kg。
# 请根据BMI公式（体重除以身高的平方）帮小明计算他的BMI指数，并根据BMI指数：
# 低于18.5：过轻
# 18.5-25：正常
# 25-28：过重
# 28-32：肥胖
# 高于32：严重肥胖
# height = 1.75
# weight = 80.5
# bmi = (1.75/80.5)**2
# if bmi < 18.5:
#     print('过轻')
# elif 18.5 <= bmi < 25:
#     print('正常')
# elif 25 <= bmi < 28:
#     print('过重')
# elif 28 <= bmi < 32:
#     print('肥胖')
# elif 32 <= bmi:
#     print('严重肥胖')

# Python的循环有两种，一种是for...in循环，依次把list或tuple中的每个元素迭代出来
# lis_03 = ['sss', 'aaaa', '杀杀杀']  # List 列表
# for lis_03 in lis_03:
#     print(lis_03)
# tuple_03 = ('eee', '到底', 10, 102.1)
# for tuple_03 in tuple_03:
#     print(tuple_03)
# 再比如我们想计算1-10的整数之和，可以用一个sum变量做累加：
# sum = 0
# for x in [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]:
#     sum = sum + x
# print(sum)
# 如果要计算1-100的整数之和，从1写到100有点困难
# 幸好Python提供一个range()函数，可以生成一个整数序列，
# 再通过list()函数可以转换为list。比如range(5)生成的序列是从0开始小于5的整数
# sum_02 = 0
# for x in range(1, 101):
#     sum_02 = sum_02 + x
# print(sum_02)

# 输出乘法表
# for i in range(1, 10):
#     for j in range(1, i + 1):
#         print('%d * %d = %d ' % (i, j, i*j), end=" ")
#     i = i + 1
#     print()

# 输出乘法表 while 循环
# i = 1
# while i <= 9:
#     j = 1
#     while j <= i:
#         print('%d * %d = %d' % (i, j, i*j), end=" ")
#         j = j + 1
#     i = i + 1
#     print()

# # 在循环中，break语句可以提前退出循环。例如，本来要循环打印1～100的数字
# num = 1
# while num <= 100:
#     print(num)
#     num = num + 1
# # 面的代码可以打印出1~100
# # 如果要提前结束循环，可以用break语句：
#     if num == 50:
#         print("exit")
#         break

# continue
# 在循环过程中，也可以通过continue语句，跳过当前的这次循环，直接开始下一次循环。
# 如果我们想只打印奇数，可以用continue语句跳过某些循环：
# num_02 = 0
# while num_02 <= 100:
#     num_02 = num_02 + 1
#     if num_02 % 2 == 0:  # 除以二没有余数 为 偶数
#         print(num_02)
#         continue
#     else:
#         print('%d 是奇数' % num_02)
#


# dict
# Python内置了字典：dict的支持，dict全称dictionary，在其他语言中也称为map
# 使用键-值（key-value）存储，具有极快的查找速度
# dic = {'jack': 94, 'tom': 90, 'bob': 70}
# print(dic["tom"])
# 把数据放入dict的方法，除了初始化时指定外，还可以通过key放入
# dic['tom'] = 100
# print(dic["tom"])
# print(dic)  # 输出整个 dict 字典
# 如果key不存在，dict就会报错：
# 要避免key不存在的错误，有两种办法，一是通过in判断key是否存在：
# print('s' in dic)
# print(dic['s'])
# 二是通过dict提供的get()方法，如果key不存在，可以返回None，或者自己指定的value：
# print(dic.get('s'))
# print(dic.get('s', '不存在')
# 要删除一个key，用pop(key)方法，对应的value也会从dict中删除：
# print(dic)
# dic.pop('bob')
# print(dic)

# set
# set和dict类似，也是一组key的集合，但不存储value。
# 由于key不能重复，所以，在set中，没有重复的key。
# 要创建一个set，需要提供一个list作为输入集合：
# set_01 = set([1, 2, 3, 4, 5])
# print(set_01)
# 重复元素在set中自动被过滤：
# set_02 = set([1, 1, 2, 2, 3, 3])


# print(set_02)
# 通过add(key)方法可以添加元素到set中，可以重复添加，但不会有效果：
# set_02.add(9)
# set_02.add(1)  # 添加重复的 会被自动过滤掉
# print(set_02)
# 通过remove(key)方法可以删除元素：
# set_02.remove(9)
# print(set_02)
# set可以看成数学意义上的无序和无重复元素的集合，
# 因此，两个set 可以做数学意义上的交集、并集等操作
# print(set_01 & set_02)  # 交集
# print(set_01 | set_02)  # 并集
# set和dict的唯一区别仅在于没有存储对应的value，但是，set的原理和dict一样，
# 所以，同样不可以放入可变对象，因为无法判断两个可变对象是否相等，
# 也就无法保证set内部“不会有重复元素”。试试把list放入set，看看是否会报错。

# 在Python中，数值类型int 、float、 字符串str 、元祖tuple、boole 都是不可变对象
# 列表list、集合set、字典dict都是可变对象

# 上面我们讲了，str是不变对象，而list是可变对象。
# 对于可变对象，比如list，对list进行操作，list内部的内容是会变化的，比如：
# a = ['c', 'b', 'a']
# a.sort()
# print(a)
# 而对于不可变对象，比如str，对str进行操作呢：
# b = 'abc'
# c = b.replace('a', 'A')
# print(c)
# print(b.replace('a', 'A'))
# print(b)
#


# # 内置函数
# num_adb = -1
# print(num_adb)
# # abs()返回绝对值 参数整数 和 浮点数
# print(abs(num_adb))

# 定义函数
# print("求绝对值")
# num_abs = input("请输入数据:")
#
#
# def my_abs(x):
# isinstance 判断数据类型
#     if not isinstance(x, (int, float, )):
#         # 抛出异常
#         raise TypeError('Parameter type error')
#     if x <= 0:
#         return abs(x)
#     else:
#         return x
#
#
# res = my_abs(num_abs)
# print("%s 的绝对值是 %d " % (num_abs, res))


# 空函数
# 如果想定义一个什么事也不做的空函数，可以用pass语句：
# def d():
#     pass


# pass语句什么都不做，那有什么用？实际上pass可以用来作为占位符，
# 比如现在还没想好怎么写函数的代码，就可以先放一个pass，让代码能运行起来。
# pass还可以用在其他语句里，比如：
# if __name__ == '__main__':
#     pass

# 解一元二次方程
# import math
#
# a = []
# b = ['参数一', '参数二', '参数三']
# i = 0
# print("解一元二次方程")
# while i >= 0:
#     c = input("请输入参数{}:".format(b[i]))
#     if c.isdigit() is False:
#         # 抛出异常
#         raise TypeError('输入的内容不是数字!请输入数字!')
#         continue
#     num = int(c)
#     if num <= 0:
#         print("请输入数字")
#         continue
#     else:
#         a.append(num)
#         i = i + 1
#         if i == 3:
#             break
#
#
# def quadratic(a, b, c):
#     d = b ** 2 - 4 * a * c
#     if d == 0:
#         x = -(b / 2 * a)
#         print(x)
#     elif d > 0:
#         x = (-b - math.sqrt(d)) / 2 * a
#         x_02 = (-b + math.sqrt(d)) / 2 * a
#         print("x1=%f,x2=%f" % (x, x_02))
#     else:
#         print("无实数根!!!")
#
#
# quadratic(a[0], a[1], a[2])


# 计算x 的 n 次方
# def cal(x, n=2):  # n=2 默认参数
#     s = 1
#     while n > 0:
#         s = s * x
#         n = n - 1
#     return s
#
#
# a = cal(2)
# b = cal(2, 2)
#
# print(a, b)


# def add_end(L=[]):
#     L.append('END')
#     return L
#
#
# # print(add_end([1, 2, 3]))
# add_end()
# add_end()
# print(add_end())


# # 可变参数
# def calc(*number):
#     s = 0
#     for i in number:
#         s = s + i
#         print(s)
#
#
# a = [1, 2, 3]
# # calc([1, 2, 3, 4, 5, 6])
# # 可变参数允许你传入0个或任意个参数，
# # 这些可变参数在函数调用时自动组装为一个tuple *a
# calc(*a)

# 关键字参数 **kw 关键字参数
# 而关键字参数允许你传入0个或任意个含参数名的参数，
# 这些关键字参数在函数内部自动组装为一个dict (key-value)。
# def info(name, age, **kw):
#     print('name:', name, 'age:', age, 'other:', kw)
#
#
# # 关键字参数有什么用？它可以扩展函数的功能。
# # 比如，在person函数里，我们保证能接收到name和age这两个参数，
# # 但是，如果调用者愿意提供更多的参数，我们也能收到。试想你正在做一个用户注册的功能，
# # 除了用户名和年龄是必填项外，其他都是可选项，利用关键字参数来定义这个函数就能满足注册的需求。
# info('张', '22', city='北京', j=100)
# # 也可以先组装出一个dict，然后，把该dict转换为关键字参数传进去：
# extra = {'city': '北京', 'add': '朝阳门', 'phone': '18565451235'}
# info('王', '30', city=extra['city'], add=extra['add'], phone=extra['phone'])
# # 简易形式
# info('李', '32', **extra)


# 命名关键字参数
# 对于关键字参数，函数的调用者可以传入任意不受限制的关键字参数。
# 至于到底传入了哪些，就需要在函数内部通过 kw 检查。
# 检查是否有 city 和 job 参数
# def info(name, age, **kw):
#     if 'city' in kw:
#         # 检查 存在 city参数
#         pass
#     if 'job' in kw:
#         # 检查 存在 job参数
#         pass
#     print(name, age, kw)
#
#
# info('杨', 10, city='上海', job='Python')
# 命名关键字参数
# def info(name, age, *, city, job):
#     print(name, age, city, job)
#
#
# info('孙', 20, city='上海', job='java')
# 命名关键字参数 带有默认值
# def info(name, age, *, city='南京', job):
#     print(name, age, city, job)
#
#
# info('孙', 20, job='java')

# ·参数组合
# *args 是可变参数，args 接收的是一个 tuple;
# **kw 是关键字参数，kw 接收的是一个 dict;
# 在 Python 中定义函数，
# 必选参数name、默认参数age=20、可变参数*args、命名关键字参数**kw
# def info(name, age=20, *args, **kw):
#     print(name, age, args, kw)
#

# 必选参数name、默认参数age=20、关键字参数a和命名关键字参数**kw
# def info2(name, age=20, *, a, **kw):
#     print(name, age, a, kw)
#
#
# # info('朱', a='w', c='23')
# c = (1, 2, 3, 4)
# b = {'a': 200}
# info(*c, **b)


# 递归函数
# 在函数内部，可以调用其他函数。如果一个函数在内部调用自身本身，这个函数就是递归函数。
# 举个例子，我们来计算阶乘 n! = 1 x 2 x 3 x ... x n，用函数 fact(n)表示，可以看出：
# def fact(n):
#     print(n)
#     if n == 1:
#         return 1
#     return n * fact(n - 1)
#
#
# print(fact(5))

# 构造一个1, 3, 5, 7, ..., 99的列表，
# 通过循环实现：
# a = []
# j = 0
# for i in range(1, 101):
#     if i % 2 != 0:
#         a.append(i)
#         print(a[j], end=' ')
#         j = j + 1

# 切片
# 取一个 list 或 tuple 的部分元素是非常常见的操作：
# lis = ['搜索', '付费', '深度']
# n = len(lis)
# for i in range(n):
#     print(i, lis[i])

# # 切片函数
# def get_num(num, s):
#     lis = []
#     if s > len(num):
#         print("超出列表元素个数 ")
#     for i in range(s):
#         lis.append(num[i])
#         print('下标{}，切出元素{}'.format(i, lis[i]))
#
#
# a = [1, 2, 3, 4, 5, 6, 7, 8]
# get_num(a, 8)

# 这种经常取指定索引范围的操作，用循环十分繁琐，
# 因此，Python 提供了切片（Slice）操作符，能大大简化这种操
# 作。对应上面的问题，取前 3 个元素，用一行代码就可以完成切片
# list切片
# tuple切片 同list ,字符串也可以切片
# a = [1, 2, 3, 4, 's', '哈哈', 43.5]
# print(a[0:6])  # 取下标0-5，6是停止位置
# print(a[:6])  # 0 可省略
# print(a[-1])  # 取倒数第一个
# print(a[-5:-1])  # 取倒数第5个到最后一个
# print(a[:5:2])  # 前5个数每两个取一个
# b = (1, 2, 3, 4, 5)
# print(b[:2])
# print('asddasdjha'[:5])


# 迭代
# list 这种数据类型虽然有下标，但很多其他数据类型是没有下标的，
# 但是，只要是可迭代对象，无论有无下标，都可以迭代，比如 dict 就可以迭代
# 可迭代对象list,dict(keys(),values(),items()),tuple,str,set,range
# dict_num = {'s': 's', 'g': 100, 'f': 1.2}
# for key in dict_num.keys():
#     print(key)
# for value in dict_num.values():
#     print(value)
# for d in dict_num.items():
#     print(d)

# print(dir(str))  # 打印内置方法
# print('__iter__' in dir(str))  # 判断是否是可迭代对象
# 如何判断一个对象是可迭代对象，方法是通过 collections 模块的 Iterable 类型判断

# from collections.abc import Iterable
#
# print(isinstance('sss', Iterable))  # str 是否可迭代

# 如果要对 list 实现类似 Java 那样的下标循环怎么办？Python 内置的 enumerate 函数可以把一个
# list 变成索引-元素对，这样就可以在 for 循环中同时迭代索引和元素本身
# a = [1, 2, 34, 5, 6.56, 4.2, 5.234, 324, 324]
# for sign, i in enumerate(a):  # sign 下标，
#     print(sign, i)

# 列表生成式
# 列表生成式即 List Comprehensions，是 Python 内置的非常简单却强大的可以用来创建 list 的生成式。
# 举个例子，要生成 list [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]可以用 list(range(1, 11))
# a = list(range(1, 11))
# print(a)
# 但如果要生成[1x1, 2x2, 3x3, ..., 10x10]怎么做？方法一是循环：
# a = []
# for i in range(1, 11):
#     a.append(i * i)
#     print(a[i - 1], end=' ')
# 但是循环太繁琐，而列表生成式则可以用一行语句代替循环生成上面的 list：
# a = [x for x in range(1, 11)]
# print(a)
# b = [x * x for x in range(1, 11)]
# print(b)
# c = [x for x in range(1, 11) if x % 2 == 0]
# print(c)
# 还可以使用两层循环，可以生成全排列
# d = [m + n for m in 'ABC' for n in 'XYZ']
# print(d)
# 写列表生成式时，把要生成的元素 x * x 放到前面，后面跟 for 循环，就可以把 list 创建出来
# 很快就可以熟悉这种语法。

# 运用列表生成式，可以写出非常简洁的代码。
# 例如，列出当前目录下的所有文件和目录名，可以通过一行代码实现：
# os.listdir(path)返回指定目录下的所有文件和目录名。
# import os
#
# # print(os.listdir('D:\\desktop'))
# a = [x for x in os.listdir('D:\\desktop')]
# # os.remove('D:\\desktop\\新建文本文档.txt')
# print(a)

# for 循环其实可以同时使用两个甚至多个变量
# 比如 dict 的 items()可以同时迭代 key 和 value

# a = {'s': '3', 'a': '张'}
# b = [k + '=' + v for k, v in a.items()]
# print(b)
#
# # 修改列表中的所有字母为小写
# L = ['ASD', 'SSS', 'vvvv', 'pppPPPP']
# print([x.lower() for x in L])

# · 练习
# 如果 list 中既包含字符串，又包含整数
# 由于非字符串类型没有 lower()方法，所以列表生成式会报错：
# lis = ['age', 1, 'name']
# print([x.lower() for x in lis])  # 错的
# lis = ['age', 1, 'name']
# print([x.lower() for x in lis if isinstance(x, str)])  # 判断元素类型


# 生成器
# 要创建一个 generator，有很多种方法。
# 第一种方法很简单，只要把一个列表生成式的[]改成()，就创建了一个generator

# L = ['ASD', 'SSS', 'vasts', 'pppPPPPs']
# L_02 = (x for x in L)
# print(type(L_02))
# for x in range(4):
#     print(next(L_02))
#


# 斐波拉契数列
# def fib(max):
#     n = 0
#     a = 0
#     b = 1
#     while n < max:
#         print(b)
#         a, b = b, a + b
#         n = n + 1
#     return 'done'

# generator
# def fib(max):
#     n = 0
#     a = 0
#     b = 1
#     while n < max:
#         yield (b)
#         a, b = b, a + b
#         n = n + 1
#     return 'done'
#
#
# f = fib(6)
# for x in f:
#     print(x)

# while True:
#     try:
#         x = next(f)
#         print('x:', x)
#     except StopIteration:
#         print(StopIteration.value, end='by jun')
#         break

# 迭代器
# 可以被 next()函数调用并不断返回下一个值的对象称为迭代器：Iterator
# 可以使用 isinstance()判断一个对象是否是 Iterator 对象
# from collections.abc import Iterator, Iterable

# print(isinstance('sss', Iterator))
# 生成器都是 Iterator 对象，但 list、dict、str 虽然是 Iterable，却不是 Iterator。
# 把 list、dict、str 等 Iterable 变成 Iterator 可以使用 iter()函数：
# a = [1, 23, 0, 55, 6, 'ff']
# b = iter(a)  # 转换为迭代器
# # print(isinstance(b, Iterator))
# while True:
#     try:
#         x = next(b)
#         print('x:', x)
#     except StopIteration:
#         print('end...', end='by jun')
#         break

# 生成器：generator
# 迭代对象：Iterable
# 迭代器：Iterator


# 高阶函数
# 编写高阶函数，就是让函数的参数能够接收别的函数。
# 把函数作为参数传入，这样的函数称为高阶函数，函数式编程就是指这种高度抽象的编程范式。
# def higher_order(a, b, f):
#     res = f(a) + f(b)
#     return res
#
#
# # 求绝对值函数 abs 作为参数传入
# res = higher_order(-99, 100, abs)
# print(res)


# map/reduce
# map()函数接收两个参数，一个是函数，一个是 Iterable，map
# 将传入的函数依次作用到序列的每个元素，并把结果作为新的 Iterator 返回
# 比如我们有一个函数 f(x)=x 2 ，要把这个函数作用在一个 list [1, 2, 3, 4, 5, 6, 7, 8, 9]上，就可
# 以用 map()实现如下
# def f(x):
#     return x * x
#
#
# lis = [1, 2, 3, 4, 5, 6, 7, 8, 9]
# a = map(f, lis)
# print(a)
# print(list(a))

# reduce
# reduce 把一个函数作用在一个序列[x1, x2, x3, ...]上，这个函数必须接收两个参数，reduce
# 把结果继续和序列的下一个元素做累积计算
# 如果考虑到字符串 str 也是一个序列，对上面的例子稍加改动，配合 map()，我们就
# # 可以写出把 str 转换为 int 的函数

from functools import reduce

#
#
# def fn(x, y):
#     print('fn:', x, y)
#     return x * 10 + y
#
#
# def char2num(s):
#     print('char2num:', s)
#     return {'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9}[s]
#
#
# # print(list(map(char2num, '13579')))
# a = reduce(fn, map(char2num, '13579'))
# print(a)

# 整合
# def str2int(s):
#     def fn(x, y):
#         return x * 10 + y
#
#     def char2num(s):
#         return {'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9}[s]
#
#     return reduce(fn, map(char2num, s))
#
#
# a = str2int('98712345')
# print(a)

# lambda匿名函数的格式：冒号前是参数，可以有多个，用逗号隔开，冒号右边的为表达式。
# 其实lambda返回值是一个函数的地址，也就是函数对象。
# a = lambda x: x * x * x
# print(a(10))

# 把 str 转换为 int 的函数
# 用lambda进一步简化
# def str2int(s):
#     def char2num(s):
#         return {'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9}[s]
#
#     return reduce(lambda x, y: x * 10 + y, map(char2num, s))
#
#
# print(str2int('4564515245'))  # str 转int


# · 练习
# 1. 利用 map()函数，把用户输入的不规范的英文名字
# 变为首字母大写，其他小写的规范名字。输入：['adam', 'LISA'
# 'barT']，输出：['Adam', 'Lisa', 'Bart']：
# def normalize(s):
#     return s.capitalize()  # 该方法返回一个首字母大写的字符串
#
#
# name = ['adam', 'LISA', 'barT', 'jun']
#
# a = list(map(normalize, name))
#
# print(a)
# 2. Python提供的 sum()函数可以接受一个 list并求和
# 请编写一个 prod()函数，可以接受一个 list并利用 reduce()求积。
# def prod(x, y):
#     return x * y
#
#
# a = [1, 2, 3, 4, 5, 6]
# res = reduce(prod, a)
# print(res)
# 3. 利用 map 和 reduce 编写一个 str2float 函数
# 把字符串'123.456'转换成浮点数 123.456：
# def str2float(s):
#     def fn(x, y):
#         return x * 10 + y
#     n = s.index('.')
#     s1 = list(map(int, [x for x in s[:n]]))
#     s2 = list(map(int, [x for x in s[n + 1:]]))
#     return reduce(fn, s1) + reduce(fn, s2) / 10 ** len(s2)
#
#
# print(str2float('123.456'))

# filter
# Python 内建的 filter()函数用于过滤序列。
# 和 map()类似，filter()也接收一个函数和一个序列。
# 和 map()不同的是，filter()把传入的函数依次作用于每个元素
# 然后根据返回值是 True 还是 False 决定保留还是丢弃该元素
# 练习1
# 在一个 list 中，删掉偶数，只保留奇数
# def is_odd(n):
#     return n % 2 == 1
#
#
# res = list(filter(is_odd, list(range(1, 101))))
# print(res)
# 练习2
# 把一个序列中的空字符串删掉
# def del_str(s):
#     return s and s.strip()
#
#
# a = list(filter(del_str, ['22', '', None, '22 2']))
# print(a)

# 装饰器
# def now():
#     print('2020-7-18')
#
#
# f = now
# f()
# print(f.__name__)
#
# 装饰器
# import time
#
#
# def log(func):
#     def wrapper(*args, **kwargs):
#         print('call %s():' % func.__name__)
#         return func(*args, **kwargs)
#     return wrapper
#
#
# @log
# def now():
#     print(time.asctime())
#
#
# if __name__ == '__main__':
#     now()





