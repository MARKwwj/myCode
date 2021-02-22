# try:
#     f = open('D:/desktop/a.txt', 'r')
#     a = f.read()
#     print(a)
# except Exception as e:
#     print('Error:', e)
# finally:
#     if f:
#         f.close()
# 但是每次都这么写实在太繁琐，
# 所以，Python 引入了 with 语句来自动帮我们调用 close()方法
# 可以反复调用 read(size)
# 方法，每次最多读取 size 个字节的内容。
# 另外，调用 readline()可以每次读取一行内容，
# 调用 readlines()一次读取所有内容并按行返回 list。因此，要根据需要决定怎么调用
# with open('D:/desktop/a.txt', 'r+') as f:
#     # for x in f.readlines():
#     #     print(x)
#     # print(f.read())
#     f.write('made in shanghai by jun...')
#

# StringIO
# 很多时候，数据读写不一定是文件，也可以在内存中读写。
# StringIO 顾名思义就是在内存中读写 str。

from io import StringIO
f = StringIO()
f.write('hello world!')
f.write(' \n')
f.write('by j')
print(f.getvalue())











































