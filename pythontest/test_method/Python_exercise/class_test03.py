# # # 判断一个对象是否是函数
# # import types
# #
# #
# # def fn():
# #     pass
# #
# #
# # print(type(fn))
# # print(type(fn) == types.FunctionType)
# # print(type(lambda x: x) == types.LambdaType)
# # print(type(abs) == types.BuiltinFunctionType)
# # print(type((x for x in range(10))) == types.GeneratorType)
#
# # 配合 getattr()获取属性 、setattr()设置一个属性 以及 hasattr()有属性xx?，我们可以直接操作一个对象的状态
# class My_object:
#     name = '张'
#
#     def __init__(self):
#         self.x = 9
#
#     def power(self):
#         return self.x * self.x
#
#
# if __name__ == '__main__':
#     obj = My_object()
#     # hasattr(obj, 'x')
#     # print(setattr(obj, 'y', 10))
#     # print(hasattr(obj, 'y'))
#     # print(getattr(obj, 'y'))
#     # print(obj.y)
#     #
#     # print(getattr(obj, 'z', 404))  # 可以传入一个 default 参数，如果属性不存在，就返回默认值
#     # print(hasattr(obj, 'power'))
#     # print(getattr(obj, 'power'))
#     # fn = getattr(obj, 'power')
#     # fn()
#     #
#     # num = obj.power()
#     # print(num)
#
#     # obj.name = '张'
#     # print(obj.name)
#     # 千万不要把实例属性和类属性使用相同的名字，因为相同名称的实例属性
#     # 将屏蔽掉类属性，但是当你删除实例属性后，再使用相同的名称，访问到的将是类属性
#     obj.name = 's'
#     del obj.name  # 删除实例的name属性
#     print(obj.name)
#     print(My_object.name)
