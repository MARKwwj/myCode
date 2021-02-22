import random


# random.randint(3, 200)
class kou:
    def __init__(self):
        self.Buckle_quantity = 1
        self.proportion = 40
        self.cur_pro = 0

    def count(self, num_order):
        i = 0
        total_cur = 0
        order_money_list = [15, 30, 99, 168, 268]
        while i < num_order:
            order_money_num = random.randint(0, 4)
            order_money = order_money_list[order_money_num]
            total_cur = total_cur + order_money
            # 当前比例
            cur_pro = (self.Buckle_quantity / total_cur) * 100

            if cur_pro < self.proportion:
                self.Buckle_quantity = self.Buckle_quantity + order_money - 1

            i = i + 1
        print(self.Buckle_quantity)


if __name__ == '__main__':
    kou = kou()
    kou.count(10000)
