from types import MethodType


# · 使用__slots__
# 定义一个特殊的__slots__变量，来限制该 class 实例能添加的属性
class Person:
    __slots__ = ('name', 'age')  # __slots__ 使用tuple（元祖）的形式定义允许绑定的属性名称
    # __slots__定义的属性仅对当前类实例起作用，对继承的子类是不起作用的
    # 除非在子类中也定义__slots__，这样，子类实例允许定义的属性就是自身的__slots__加上父类的__slots__。


class chilren(Person):
    __slots__ = ('add')
    pass


per = Person()
# per.add = '上海'  # 绑定一个 允许范围外的属性
# print(per.add)  # 运行会报错 AttributeError
per.name = '张'
per.age = 23
print(per.name, per.age)

# __slots__定义的属性仅对当前类实例起作用，对继承的子类是不起作用的
chil = chilren()
chil.add = '上海'
print(chil.add)
chil.hobby = '旅游'
# 运行会报错 AttributeError
print(chil.hobby)  # 原因：子类实例允许定义的属性就是自身的__slots__加上父类的__slots__
