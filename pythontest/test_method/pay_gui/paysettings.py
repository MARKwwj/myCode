class Paysettings:
    def __init__(self):
        """初始化 支付参数"""
        self.channel_id = '27'
        self.token = 'eyJhbGciOiJIUzUxMiJ9.eyJsb2dpbl91c2VyX2tleSI6IjUyZjVhYTQ4LTgwZDAtNDA1OC1iM2RlLTMwMTBjMTEyZWEyOSJ' \
                     '9.w0qu_6SygJ8D5JaXICu8z_JwxUhLlyJQ8KPFhxRW8u5OkZhxCE4G-PAhcNKqmVDJNhC049UWSsY2vOEyeKbuag'
        # 订单状态
        self.status = '1'
        # 支付类型
        self.type_pay = '1'
        # 商户id
        self.mch_id = '15520208888'
        # 商户密钥
        self.key = 'e07a7c0bedc507a53172e76d4b60af3c'
