import time

from telegram.ext import Updater
from telegram.ext import CommandHandler
# import telegram
import logging
import jenkins

from telegram_my.tele import automessage

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

token = '1445191367:AAHCq-M_2cOCzx4GKFSDgHDGHnvrlx4akds'


# bot = telegram.Bot(token)
# print(bot.get_me())


def start(update, context):
    context.bot.send_message(chat_id=update.effective_chat.id, text="打包中...")
    result = jenkins_run()
    if result == "SUCCESS":
        text = """
        打包成功！
        """
    else:
        text = """
        打包异常！
        """
    context.bot.send_message(chat_id=update.effective_chat.id, text=text)


def help(update, context):
    text = """
    命令提示：
    打包长视频： /longvideo 
    """
    context.bot.send_message(chat_id=update.effective_chat.id, text=text)


def telegram_run():
    updater = Updater(token, use_context=True)
    dispatcher = updater.dispatcher

    start_handler = CommandHandler('longvideo', start)
    help_handler = CommandHandler('help', help)

    dispatcher.add_handler(start_handler)
    dispatcher.add_handler(help_handler)

    updater.start_polling()


def jenkins_run():
    jenkins_server_url = "http://192.168.100.241:8080/jenkins/"
    user_id = "flutter_mini"
    api_token = "11b35f30d496e11ba83a5e50b6d422c6ed"
    job_name = "new_long"
    # 实例化jenkins对象，连接远程的jenkins master server
    server = jenkins.Jenkins(jenkins_server_url, username=user_id, password=api_token)
    # 构建job名为job_name的job（不带构建参数）
    server.build_job(job_name)
    # 获取job名为job_name的job的相关信息
    job_info = server.get_job_info(job_name)
    # 获取job名为job_name的job的最后次构建号
    next_build_number = server.get_job_info(job_name)['nextBuildNumber']
    print("build_number:", next_build_number)
    while True:
        try:
            # 判断job名为job_name的job的某次构建是否还在构建中
            job_whether_is_build = server.get_build_info(job_name, next_build_number)['building']
            # print("job_whether_is_build:", job_whether_is_build)
            if job_whether_is_build is False:
                print("构建结束！")
                break
        except Exception as e:
            print(e)
            continue
    # 获取job名为job_name的job的某次构建的执行结果状态
    job_result = server.get_build_info(job_name, next_build_number)['result']
    print("job_result:", job_result)
    return job_result


telegram_run()
