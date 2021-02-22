class Student:
    # 在 Python 中，实例的变量名如果以__开头，
    # 就变成了一个私有变量（private），只有内部可以访问，外部不能访问
    def __init__(self, name, age, add):
        self.__age = age
        self.__name = name
        self.__add = add

    def print_age(self):
        print('age:{}'.format(self.__age))

    def print_status(self):
        if self.__age >= 10 & self.__age < 18:
            print('青少年')
        elif self.__age >= 18 & self.__age < 30:
            print('年轻')
        elif self.__age >= 30:
            print('oh!@')

    # 修改name
    def set_name(self, name):
        self.__name = name

    # 获取name
    def get_name(self):
        return self.__name


if __name__ == '__main__':
    stu1 = Student('张', 21, '上海')

    # print(stu1.age, stu1.name)
    # print(stu1)
    # stu1.print_age()
    # stu1.print_status()
    # stu1.set_name('哈哈')
    # print(stu1.get_name())
    # stu1.print_status()
    # stu1.print_age()
    # 不能直接访问__name(private,私有的属性) 是因为 Python 解释器对外
    # 把__name 变量改成了_Student__name，所以，仍然可以通过_Student__name 来访问__name 变量
    # 但是强烈建议你不要这么干，不同版本的 Python 解释器可能会把__name 改成不同的变量名
    # print(stu1._Student__name)
    # stu1.__name = '吴'
    # print(stu1.__name)  # 内部的__name 变量已经被 Python 解释器自动改成了_Student__name
    # print(stu1.get_name())  # 内部的__name 变量已经被 Python 解释器自动改成了_Student__name
