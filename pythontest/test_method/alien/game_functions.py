import sys
import pygame
from ship import Ship
from bullet import Bullet
from alien import Alien
from time import sleep


def check_events(ai_settings, screen, bullets, stats, play_button, aliens, ship_real, sb):
    """响应按键和鼠标事件"""
    for events in pygame.event.get():
        if events.type == pygame.QUIT:
            sys.exit()
        elif events.type == pygame.MOUSEBUTTONDOWN:
            mouse_x, mouse_y = pygame.mouse.get_pos()
            check_play_button(stats, play_button, mouse_x, mouse_y, aliens, ai_settings, screen, ship_real, bullets, sb)

        elif events.type == pygame.KEYDOWN:
            check_keydown_events(events, ship_real, ai_settings, screen, bullets)
        elif events.type == pygame.KEYUP:
            check_keyup_events(events, ship_real)


def check_play_button(stats, play_button, mouse_x, mouse_y, aliens, ai_settings, screen, ship_real, bullets, sb):
    """在玩家单机 play 时开始新游戏"""
    button_clicked = play_button.rect.collidepoint(mouse_x, mouse_y)
    if button_clicked and not stats.game_active:
        # 重置游戏设置
        ai_settings.initialize_dynamic_settings()
        # 开始游戏后 隐藏光标
        pygame.mouse.set_visible(False)

        # 重置 游戏信息
        stats.reset_stats()
        stats.game_active = True


        # 重置记分牌图像
        sb.prep_high_score()
        sb.prep_level()
        sb.prep_score()
        sb.prep_ships()

        # 清空外星人列表和子弹列表
        aliens.empty()
        bullets.empty()

        # 创建一批新外星人，并让飞船居中
        create_fleet(ai_settings, screen, aliens, ship_real)
        ship_real.center_ship()


def check_keydown_events(events, ship_real, ai_settings, screen, bullets):
    """响应 按键"""
    if events.key == pygame.K_RIGHT:
        # 向右移动飞船
        ship_real.moving_right = True
    elif events.key == pygame.K_LEFT:
        # 向右移动飞船
        ship_real.moving_left = True
    if events.key == pygame.K_SPACE:
        fire_bullet(bullets, ai_settings, screen, ship_real)


def check_keyup_events(events, ship_real):
    """响应 松开"""
    if events.key == pygame.K_RIGHT:
        ship_real.moving_right = False
    elif events.key == pygame.K_LEFT:
        ship_real.moving_left = False


def update_screen(ai_settings, screen, ship_real, bullets, aliens, stats, play_button, sb):
    """更新屏幕上的图像，并切换到新屏幕"""
    # 每次循环时重绘屏幕
    screen.fill(ai_settings.bg_color)
    # 绘制子弹
    for bullet in bullets.sprites():
        bullet.draw_bullet()
    # 绘制飞船
    ship_real.blitme()
    # 绘制外星人
    aliens.draw(screen)
    # 显示得分
    sb.show_score()
    if not stats.game_active:
        play_button.draw_button()

    # 让最近绘制的屏幕可见
    pygame.display.flip()


def update_bullets(ai_settings, screen, aliens, ship_real, bullets, stats, sb):
    """更新子弹的位置， 并删除已消失的子弹"""
    bullets.update()

    # 删除超出屏幕的子弹
    for bullet in bullets.copy():
        if bullet.rect.bottom <= 0:
            bullets.remove(bullet)
    # print(len(bullets))
    check_bullet_alien_collisions(ai_settings, screen, aliens, ship_real, bullets, stats, sb)
    check_aliens_bottom(ai_settings, screen, aliens, bullets, ship_real, stats,sb)


def check_bullet_alien_collisions(ai_settings, screen, aliens, ship_real, bullets, stats, sb):
    """响应子弹和外星人的碰撞"""
    # 检查是否有子弹击中了外星人，
    # 如果有，就删除相应的子弹和外星人
    collections = pygame.sprite.groupcollide(bullets, aliens, True, True)
    if collections:
        for aliens in collections.values():
            stats.score += ai_settings.alien_points * len(aliens)
            sb.prep_score()
        check_high_score(stats, sb)
    if len(aliens) == 0:
        # 清空子弹
        bullets.empty()
        # 如果整群外星人都被消灭就提高一个等级
        ai_settings.increase_speed()
        # 提高等级
        stats.level += 1
        sb.prep_level()
        # 再创建一批外星人
        create_fleet(ai_settings, screen, aliens, ship_real)


def fire_bullet(bullets, ai_settings, screen, ship_real):
    """飞船 开火"""
    # 创建一颗子弹，并将其加入编组bullets 中
    if len(bullets) < ai_settings.bullets_allowed:
        # if ai_settings.bullets_num < ai_settings.bullets_allowed:
        new_bullet = Bullet(ai_settings, screen, ship_real)
        bullets.add(new_bullet)
        # ai_settings.bullets_num += 1
        # print(ai_settings.bullets_num)


def get_number_aliens_x(ai_settings, alien_width):
    """ 计算每行可容纳多少个外星人 """
    # 一行可容纳外星人的总宽度
    available_space_x = ai_settings.screen_width - 2 * alien_width
    # 一行可容纳的外星人数量
    number_aliens_x = int(available_space_x / (2 * alien_width))
    return number_aliens_x


def get_number_rows(ai_settings, ship_height, alien_height):
    """计算屏幕可容纳多少行外星人"""
    available_space_y = (ai_settings.screen_height - (4 * alien_height) - ship_height)
    number_rows = int(available_space_y / (2.5 * alien_height))
    return number_rows


def create_alien(screen, ai_settings, aliens, alien_number, row_number):
    """ 创建一个外星人并将其放在当前行 """
    alien = Alien(screen, ai_settings)
    alien_width = alien.rect.width
    alien.x = alien_width + 2 * alien_width * alien_number
    alien.rect.x = alien.x
    alien.rect.y = alien.rect.height + 2 * alien.rect.height * row_number
    aliens.add(alien)


def create_fleet(ai_settings, screen, aliens, ship_real):
    """创建一群外星人"""
    alien = Alien(screen, ai_settings)
    number_aliens_x = get_number_aliens_x(ai_settings, alien.rect.width)
    number_rows = get_number_rows(ai_settings, ship_real.rect.height, alien.rect.height)

    for row_number in range(number_rows):
        for alien_number in range(number_aliens_x):
            # 创建一个外星人并加入当前行
            create_alien(screen, ai_settings, aliens, alien_number, row_number)


def check_fleet_edges(ai_settings, aliens):
    """检查外星人是否移动到边缘"""
    for alien in aliens.sprites():
        if alien.check_edges():
            change_fleet_direction(ai_settings, aliens)
            break


def change_fleet_direction(ai_settings, aliens):
    """将整个外星人群下移，并改变他们的方向"""
    for alien in aliens.sprites():
        alien.rect.y += ai_settings.fleet_drop_speed
    ai_settings.fleet_direction *= -1


def update_aliens(ai_settings, screen, aliens, bullets, ship_real, stats, sb):
    """检查是否有外星人位于屏幕边缘，并更新整群外星人的位置"""
    check_fleet_edges(ai_settings, aliens)
    aliens.update()
    # 检测外星人和飞船之间的碰撞
    if pygame.sprite.spritecollideany(ship_real, aliens):
        print("ship hit!!!")
        ship_hit(ai_settings, screen, aliens, bullets, ship_real, stats, sb)

    # 检查是否有外星人抵达屏幕底端
    check_aliens_bottom(ai_settings, screen, aliens, bullets, ship_real, stats, sb)


def ship_hit(ai_settings, screen, aliens, bullets, ship_real, stats, sb):
    """响应被外星人 撞到的飞船"""
    if stats.ships_left > 0:
        # 将 ship_left减1
        stats.ships_left -= 1

        # 更新记分牌
        sb.prep_ships()

        # 清空外星人列表和子弹列表
        aliens.empty()
        bullets.empty()
        # 创建 一群新的外星人，并将飞船重新放置底端中央
        create_fleet(ai_settings, screen, aliens, ship_real)
        ship_real.center_ship()
        # 暂停
        sleep(0.7)
    else:
        stats.game_active = False
        pygame.mouse.set_visible(True)


def check_aliens_bottom(ai_settings, screen, aliens, bullets, ship_real, stats,sb):
    """检查是否有外星人到达屏幕底端"""
    screen_rect = screen.get_rect()
    for alien in aliens.sprites():
        if alien.rect.bottom >= screen_rect.bottom:
            # 像飞船被撞到一样处理
            ship_hit(ai_settings, screen, aliens, bullets, ship_real, stats,sb)
            break


def check_high_score(stats, sb):
    """检查是否诞生了最高分"""
    if stats.score > stats.high_score:
        stats.high_score = stats.score
        sb.prep_high_score()
