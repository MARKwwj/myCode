from datetime import datetime
import random


class all_infos:
    def __init__(self):
        """视频信息"""
        self.title = ""  # 视频标题 ok
        self.introduce = ""  # 视频简介 ok
        self.pic_url = ""  # 视频封面路径 ok
        self.video_url = ""  # 视频路径(m3u8) ok
        self.video_times = 0  # 视频时长 ok
        self.cover_pic_quantity = 0  # 封面图数量 ok
        self.charge = 0  # 视频金币价格   ok
        self.video_shard_size = 0  # 视频分片总大小 ok

        """ 以下参数使用默认配置 不需要修改 """
        self.video_size = 0  # 视频大小
        self.chargeType = "3"  # 视频收费类型
        self.del_flag = "0"  # 删除标志 0代表存在
        self.res_id = "r8"  # 资源服务器ID
        self.status = "5"  # 视频状态 5 待审核
        self.score = round(8.8 + random.uniform(0.1, 0.7), 1)  # 视频评分
        self.view_count = 10000 + random.randint(1, 1490000)  # 视频播放量
        self.zan_count = 100 + random.randint(1, 900)  # 视频点赞数
        self.collect_count = 100 + random.randint(1, 900)  # 视频收藏数
        self.create_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')  # 创建时间
        self.update_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')  # 更新时间
        # self.videos_path_server = '/home/resources/wwwroot/da_xiang/crypted/videos'
        # self.videos_path_server_encrypt = '/home/resources/wwwroot/da_xiang/uncrypt/videos'

        # self.videos_path_server = 'D:/desktop/c1'
        # self.videos_path_server_encrypt = 'D:/desktop/c2'
        # self.videos_path_server = '/home/xiaowu/videos'
        # self.videos_path_server_encrypt = '/home/xiaowu/videos2'

        """数据库参数"""
        # 测试数据库
        # self.host = '64.64.241.63'
        # self.port = 6033
        # self.user = 'root'
        # self.password = 'ZsNice2020.'
        # self.database = 'manage_db'
        # self.host = '192.168.100.51'
        # self.port = 3306
        # self.user = 'root'
        # self.password = '123456'
        # self.database = 'manage_db'
        # self.host = '80.251.223.22'
        # self.port = 6033
        # self.user = 'root'
        # self.password = 'ZsNice2020.'
        # self.database = 'manage_db'
