import paramiko

import re
import multiprocessing

import pymysql
import traceback


def con_sql():
    sql = 'select pic_url from video '
    pic_url_list = []
    try:
        conn = pymysql.connect(
            host='80.251.223.22',
            port=6033,
            charset='UTF8',
            passwd='Admin2020.',
            database='manage_db',
            user='readonly'
        )
        cmd = conn.cursor()
        cmd.execute(sql)
        sql_res = cmd.fetchall()
        for r in sql_res:
            if re.search(r'videos\/\w{32}\/covers\/cover-\d{1,}.webp\?resId=\w{2}', str(r), re.I):
                reg_film_name = re.compile(r"\w{32}", re.I)
                reg_film_cover = re.compile(r"cover-\d{1,}.webp", re.I)
                res_filmname = re.findall(reg_film_name, str(r))
                res_filmcover = re.findall(reg_film_cover, str(r))
                pic_url_list.append(res_filmname + res_filmcover)
        return pic_url_list
    except Exception:
        print('处理异常: ' + traceback.format_exc())
    finally:
        conn.close()


def upload_file(ip, port, password, process_name):
    # 用户名和密码
    transport = paramiko.Transport((ip, port))
    transport.connect(username='root', password=password)
    sftp = paramiko.SFTPClient.from_transport(transport)
    pic_url_list = con_sql()
    pic_url_list.reverse()
    for l in pic_url_list:
        try:
            sftp.chdir('/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/covers/')
            print(process_name + " " + "dir is exists")
            # 新文件夹路径
            new_dir_path = '/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/setcover'
            # 新图片路径
            new_file_path = '/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/setcover/' + l[1]
            # 创建新文件夹
            sftp.mkdir(new_dir_path)
            # 旧文件图片路径
            # /wwwroot/da_xiang/uncrypt/videos/6b91b079491ba665f662ed4e1d7bf65b/covers/cover-6.webp
            old_file_path = '/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/covers/' + l[1]
            print(process_name + " " + old_file_path)
            # 以二进制只读打开文件
            f_old_pic = sftp.open(old_file_path, 'rb', bufsize=1024 * 1024 * 5)
            pic_res = f_old_pic.read()
            f_new_pic = sftp.open(new_file_path, 'wb', bufsize=1024 * 1024 * 5)
            f_new_pic.write(pic_res)

            f_old_pic.close()
            f_new_pic.close()
        except Exception as e:
            print("当前进程{0}，{1}".format(process_name, e))


if __name__ == '__main__':
    # cdn1
    process1 = multiprocessing.Process(target=upload_file, args=('80.251.223.68', 27565, 'cgsVsEW5F5Xp', 'process1'))
    # cdn3
    process2 = multiprocessing.Process(target=upload_file, args=('216.24.190.133', 26827, 'XDcCCAb73AEr', 'process2'))
    # cdn5
    process3 = multiprocessing.Process(target=upload_file, args=('216.24.185.123', 28457, 'aX29h4cNVAue', 'process3'))

    process1.start()
    process2.start()
    process3.start()
