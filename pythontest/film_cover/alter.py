import os
import re

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


def pic_alter():
    pic_url_list = con_sql()
    for l in pic_url_list:
        try:
            if os.path.exists('/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/covers/'):
                print('dir is exist')
                # 新文件夹路径
                new_dir_path = '/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/setcover'
                # 新图片路径
                new_file_path = '/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/setcover/cover.webp'
                os.mkdir(new_dir_path)
                # 旧文件图片路径
                old_file_path = '/home/resources/wwwroot/da_xiang/crypted/videos/' + l[0] + '/covers/' + l[1]
                # 以二进制只读打开文件
                f_old_pic = open(old_file_path, 'rb', bufsize=1024 * 1024 * 5)
                pic_res = f_old_pic.read()
                f_new_pic = open(new_file_path, 'wb', bufsize=1024 * 1024 * 5)
                f_new_pic.write(pic_res)

                f_old_pic.close()
                f_new_pic.close()
        except Exception as e:
            print(e)


pic_alter()
