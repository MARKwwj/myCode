import sys
import pygame
from settings import Settings
from ship import Ship
from pygame.sprite import Group
import game_functions as gf
from bullet import Bullet
from alien import Alien
from game_stats import GameStats
from button import Button
from scoreboard import Scoreboard


def run_game():
    # 初始化游戏并创建一个屏幕对象
    pygame.init()
    # 实例化
    ai_settings = Settings()
    # 创建一个显示窗口,指定了游戏窗口的尺寸
    screen = pygame.display.set_mode((ai_settings.screen_width, ai_settings.screen_height))
    # 给窗口定义一个标题
    pygame.display.set_caption("Alien Invasion")
    # 创建一艘飞船
    ship_real = Ship(screen, ai_settings)
    # 创建一个编组用来存储子弹
    bullets = Group()
    # 创建一个外星人编组
    aliens = Group()
    # 创建外星人群
    gf.create_fleet(ai_settings, screen, aliens, ship_real)
    # 创建一个 用于存储游戏统计信息的实例
    stats = GameStats(ai_settings)
    # 创建一个“play”按钮
    play_button = Button(ai_settings, screen, 'play')
    # 创建存储游戏统计信息的实例
    sb = Scoreboard(ai_settings, screen, stats)

    # 开始游戏的主循环
    while True:
        # 响应按键和鼠标事件
        gf.check_events(ai_settings, screen, bullets, stats, play_button, aliens, ship_real, sb)
        if stats.game_active:
            # 根据移动标志调整飞船方向
            ship_real.update()
            # 更新子弹位置，并清除超出屏幕外的子弹
            gf.update_bullets(ai_settings, screen, aliens, ship_real, bullets, stats, sb)
            # 更新外星人位置
            gf.update_aliens(ai_settings, screen, aliens, bullets, ship_real, stats, sb)
        # 更新屏幕上的图像，并切换到新屏幕
        gf.update_screen(ai_settings, screen, ship_real, bullets, aliens, stats, play_button, sb)


if __name__ == '__main__':
    run_game()
