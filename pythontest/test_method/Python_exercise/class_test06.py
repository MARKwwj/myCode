# @property
class Student:
    @property  # 把一个 getter 方法变成属性
    def score(self):
        return self._score  # _socre 单下划线开头，表示私有不被导出
        # 单个下划线开头的名称只能在当前文件使用，不能导出到其他文件调用

    @score.setter  # 把一个 setter 方法变成属性赋值
    def score(self, score):
        if not isinstance(score, int):
            raise ValueError('score must be an integer!')
        elif score < 0 or score > 100:
            raise ValueError('score must between 0 ~ 100!')
        self._score = score


stu = Student()
# stu.set_score(-1)
# res = stu.get_score()
# print(res)
# stu.score = 99  # 实际转化为 s.set_score(60)
# print(stu.score)  # 实际转化为 s.get_score()


# 只定义 getter 方法，不定义 setter 方法就是一个只读属性
class Person:

    @property
    def name(self):
        return self._name

    @name.setter
    def name(self, name):
        self._name = name

    @property  # 只读
    def new_name(self):
        return 'new' + self._name


p = Person()
# p._name = '张'
# print(p._name)
# p.set_name('孙')
# print(p._name)
p.s
