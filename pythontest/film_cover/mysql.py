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



