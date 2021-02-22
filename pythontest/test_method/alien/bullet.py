import pygame
from pygame.sprite import Sprite


class Bullet(Sprite):
    """一个对飞船发出的子弹 进行管理的类"""

    def __init__(self, ai_settings, screen, ship_real):
        """在飞船所处的位置创建一个子弹对象"""

        """
        在super机制里可以保证公共父类仅被执行一次
        执行的顺序，是按照MRO：方法解析顺序进行的
        执行完当前类按照mro顺序执行下个类
        """
        # 调用父类 Sprite
        super().__init__()
        self.screen = screen
        # 在(0，0)处创建一个表示子弹的矩形，在设置其位置 rect（left,top,width,height）
        # 子弹并非基于图像的，因此我们必须使用 pygame.Rect() 类从空白开始创建一个矩形
        self.rect = pygame.Rect(0, 0, ai_settings.bullet_width, ai_settings.bullet_height)
        # 设置子弹的 中央x轴坐标 等于飞船的centerx坐标
        self.rect.centerx = ship_real.rect.centerx
        # 设置子弹的 top 的坐标 等于飞船的top坐标
        self.rect.top = ship_real.rect.top
        # 存储用小数表示的子弹位置
        self.y = float(self.rect.y)
        # 颜色
        self.color = ai_settings.bullet_color
        # 速度
        self.bullet_factor = ai_settings.bullet_speed_factor

    def update(self):
        """子弹是不断向上移动的"""
        self.y -= self.bullet_factor
        self.rect.y = self.y

    def draw_bullet(self):
        """在屏幕上绘制子弹"""
        pygame.draw.rect(self.screen, self.color, self.rect)
