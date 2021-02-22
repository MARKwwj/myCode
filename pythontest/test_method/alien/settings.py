class Settings:
    def __init__(self):
        """初始化游戏的设置"""
        """ 屏幕的设置"""
        self.screen_width = 1200
        self.screen_height = 800
        # 窗口背景颜色
        self.bg_color = (224, 224, 224)

        """飞船设置"""
        # 玩家拥有的飞船数
        self.ship_limit = 3

        """子弹设置"""
        # 子弹设置 创建宽 3 像素、高 15 像素的深灰色子弹
        self.bullet_width = 5
        self.bullet_height = 15
        self.bullet_color = (60, 60, 60)
        # 当前发射过的子弹数量
        self.bullets_num = 0
        # 允许发射的子弹数量
        self.bullets_allowed = 100

        """外星人设置"""
        # 外星人的纵向移动速度
        self.fleet_drop_speed = 10
        """以什么样的速度加快游戏节奏"""
        self.speed_scale = 1.1
        self.initialize_dynamic_settings()
        # 外星人 点数提高的速度
        self.score_scale = 1.1

    def initialize_dynamic_settings(self):
        """初始化随游戏进行而变化的设置"""
        # 飞船移动速度
        self.ship_speed_factor = 1
        # 子弹速度
        self.bullet_speed_factor = 3
        # 外星人的横向移动速度
        self.alien_speed_factor = 1

        # fleet_direction 为 1 表示向右移，为 -1 表示向左移（外星人）
        self.fleet_direction = 1
        # 击落一个外星人得分
        self.alien_points = 50

    def increase_speed(self):
        """提高速度设置"""
        self.ship_speed_factor *= self.speed_scale
        self.bullet_speed_factor *= self.speed_scale
        self.alien_speed_factor *= self.speed_scale
        self.alien_points = int(self.alien_points * self.score_scale)
