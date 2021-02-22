import uuid


class All_infos():
    def __init__(self):
        """数据库参数"""
        self.host = '199.180.114.169'
        self.port = 6033
        self.user = 'root'
        self.password = 'ZsNice2020.'
        self.database = 'long_video_db'

        """支付参数"""
        self.channel_id = 0  # 53 支付宝 54 微信
        self.trade_status = ""  # 支付状态

        self.user_whether_exists = 0  # 用户是否存在 1存在 0不存在
        self.app_name = ""
        self.app_names_list = ['long_video', 'novel']
        self.app_type = 0  # app类型
        self.user_id = 0  # 用户ID

        self.commodity_vip_info = []  # list  id name money
        self.commodity_gold_info = []  # list  id name money
        self.commodity_index = 0  # list 下标
        self.card_name_list = []

        self.commodity_type = 0  # 充值类型 1vip  2金币
        self.commodity_id = 0  # 1 月卡 2 季卡 3 年卡 4 永久卡
        self.commodity_amount = 0  # 套餐金额

        # 更换用户的ID
        self.need_changed_user_id = 0
        self.machine_code = ""


