# python 天生自带多态效果（继承object类）
class Animal:
    def run(self):
        print('Animal is running...')

    def run_twice(self):
        self.run()


class Dog(Animal):  # 继承 Animal

    def run(self):
        print('Dog is running...')


class Cat(Animal):  # 继承 Animal

    def run(self):
        print('Cat is running...')


def run_twice(animal):
    animal.run()


if __name__ == '__main__':
    dog = Dog()
    # dog.run()
    cat = Cat()
    # cat.run()
    ani = Animal()
    # ani.run()
