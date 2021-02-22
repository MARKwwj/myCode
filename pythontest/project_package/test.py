import json
import time

import paramiko

from project_package.new_upload import encrypt


def test():
    while True:
        try:
            transport = paramiko.Transport(('199.180.114.169', 26091))
            transport.connect(username='root', password='VzyC3RmT4wRY')
            transport.banner_timeout = 60
            sftp = paramiko.SFTPClient.from_transport(transport)
            ssh = paramiko.SSHClient()
            ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            ssh.connect('199.180.114.169', 26091, 'root', 'VzyC3RmT4wRY', timeout=60)
            print('资源服务器 连接成功！！！')
            break
        except Exception:
            print('资源服务器 连接失败！！！')
            time.sleep(3)
            print('三秒后重连！！！')
            continue
    sftp.chdir("/home/resources/wwwroot/da_xiang/crypted")
    sftp.remove("ver.json")
    ssh.exec_command("sh {0}test.sh".format("/home/resources/wwwroot/da_xiang/crypted/"))
    encrypt("ver.json", "ver_uncrypt.json", sftp)
    with sftp.open('ver_uncrypt.json') as fr:
        result = fr.read()
        json_content = json.loads(result.decode())
    print("mv app.json app_{0}.json".format(json_content["app"]))
    print("mv api.json api_{0}.json".format(json_content["api"]))
    print("mv zip.zip zip_{0}.zip".format(json_content["zip"]))

    app_json_name = "app_{0}.json".format(json_content["app"])
    api_json_name = "api_{0}.json".format(json_content["api"])
    zip_zip_name = "zip_{0}.zip".format(json_content["zip"])
    try:
        sftp.stat(app_json_name)
        print(app_json_name, "文件已存在！！跳过")
        sftp.remove("app.json")
    except IOError:
        sftp.rename("app.json", app_json_name)

    try:
        sftp.stat(api_json_name)
        print(api_json_name, "文件已存在！！跳过")
        sftp.remove("api.json")
    except IOError:
        sftp.rename("api.json", api_json_name)

    try:
        sftp.stat(zip_zip_name)
        print(zip_zip_name, "文件已存在！！跳过")
        sftp.remove("zip.zip")
    except IOError:
        sftp.rename("zip.zip", zip_zip_name)

    sftp.close()
    ssh.close()


test()
