import pygame
from pygame.sprite import Sprite


class Ship(Sprite):
    def __init__(self, screen, ai_settings):
        """初始化飞船，并设置初始位置"""
        super().__init__()
        self.screen = screen
        self.ai_settings = ai_settings
        # 加载飞船图像并获取飞船外界矩形
        self.image = pygame.image.load('D:/python_test/test_method/alien/images/shipmini2.png')
        # 获取图片的矩形
        self.rect = self.image.get_rect()
        # 获取屏幕的矩形
        self.screen_rect = screen.get_rect()
        # 将每艘新飞船 放置屏幕底部中央
        self.rect.centerx = self.screen_rect.centerx
        self.rect.bottom = self.screen_rect.bottom
        # 在飞船的属性center中存储小数值
        self.center = float(self.rect.centerx)

        # 移动标志
        self.moving_right = False
        self.moving_left = False

    def blitme(self):
        """在指定位置绘制飞船"""
        self.screen.blit(self.image, self.rect)

    def update(self):
        """根据移动标志调整飞船方向"""
        if self.moving_right and self.rect.centerx < self.screen_rect.right:
            self.center += self.ai_settings.ship_speed_factor
        elif self.moving_left and self.rect.centerx > 0:
            self.center -= self.ai_settings.ship_speed_factor

        self.rect.centerx = self.center

    def center_ship(self):
        """让飞船再屏幕上居中"""
        self.center = self.screen_rect.centerx
