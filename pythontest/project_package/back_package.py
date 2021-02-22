import os
import pexpect

apps_server_project = '/Users/mini/Documents/gitauto/zsproject'
back_server = '/Users/mini/Documents/gitauto/app_back_server'
back_web = '/Users/mini/Documents/gitauto/app_back_web'


def update_code():
    os.chdir(apps_server_project)
    os.system('git restore .')
    os.system('git checkout .')
    # 自动输入git 密码
    child = pexpect.spawnu('git pull origin master')
    child.expect("git@192.168.100.51's password:")
    child.sendline('git\r')
    child.close()


def sftp_package():
    ...







