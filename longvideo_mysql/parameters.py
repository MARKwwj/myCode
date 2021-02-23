from datetime import datetime
import random


class all_infos:
    def __init__(self):
        """视频信息"""
        self.video_id = 0  # 视频id                   ok
        self.video_title = ""  # 视频标题               ok
        self.video_intro = ""  # 视频简介               ok
        self.category_id = ""  # 视频类别               ok
        self.classify_id = ""  # 视频分类               ok
        self.video_tags_id = []  # 标签id                      ok
        self.video_tags = None  # 视频标签                    ok
        self.video_duration = 1  # 视频时长                     ok
        self.video_cover_quantity = 0  # 封面图数量              ok
        self.video_pay_coin = 0  # 视频金币价格                ok
        self.video_byte_size = 0  # 视频大小字节             ok

        """ 以下参数使用默认配置 不需要修改 """
        self.dict_videos_title = {}
        self.video_machine_id = 1
        self.video_machine_name = 'r1'
        self.tag_parent_id = 27  # 标签父级id
        self.video_pay_type = 3  # 视频收费类型
        self.video_status = 4  # 视频状态 4 待编辑
        self.video_creator = "reptile"  # 视频创建者
        self.video_score = round(8.8 + random.uniform(0.1, 0.7), 1)  # 视频评分
        self.video_total_score = self.video_score  # 视频总评分
        self.video_score_total_people = random.randint(1, 100)  # 视频评分认识
        self.video_play_count = random.randint(12999, 999999)  # 视频播放量
        self.video_praise_count = random.randint(1299, 9999)  # 视频点赞数
        self.video_favorite_count = random.randint(129, 999)  # 视频收藏数
        self.video_create_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')  # 创建时间
        self.video_update_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')  # 更新时间

        """数据库参数"""
        # 测试数据库
        self.host = '192.168.100.51'
        self.port = 3306
        self.user = 'root'
        self.password = '123456'
        self.database = 'res_video_db'
        # 生产数据库
        # self.host = '110.92.66.88'
        # self.port = 6033
        # self.user = 'root'
        # self.password = 'yxPvqJlbYBRrIs0z'
        # self.database = 'res_video_db'
        # 测试数据库
        # self.host = '199.180.114.169'
        # self.port = 6033
        # self.user = 'root'
        # self.password = 'ZsNice2020.'
        # self.database = 'res_video_db'

        """资源目录"""
        # 视频 资源目录  cur_video_path 爬虫下载的目录  res_video_path 服务器视频存放目录
        # self.cur_video_path = "/home/datadrive/resources/res9/Data"
        # self.res_video_path = "/home/datadrive/resources/res1/da_xiang/crypted/videos"
        self.cur_video_path = "D:\\desktop\\r1"
        self.res_video_path = "D:\\desktop\\r2"
