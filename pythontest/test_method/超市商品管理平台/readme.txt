工具: Python3,Pycharm,MySQL,Navicat for MySQL
创建MySQL数据库 :
	create database [数据库名称] charset utf8;
在config文件里面的DATABASE修改为自己创建的数据库名称
运行代码前需要运行MySQL
注意代码中一些步骤的注释含义
导入的包需要用pip安装或者用pycharm的Project Interpreter添加
第一次运行代码时，需要将index()函数下的 db.create_all()打开再运行
第一次进入网站界面时直接登录会报错，因为数据库此时为空
      需要点击注册或者自己在程序中或者数据库里添加
注意输入时间格式为XXXX-XX-XX的形式
