import os
import re

import pymysql
import traceback

from pymysql.cursors import DictCursor


def con_sql():
    sql = 'select id,pic_url from video where del_flag = 0 '
    try:
        conn = pymysql.connect(
            host='80.251.223.22',
            port=6033,
            charset='UTF8',
            passwd='Admin2020.',
            database='manage_db',
            user='readonly'
        )
        cmd = conn.cursor(DictCursor)
        cmd.execute(sql)
        sql_res = cmd.fetchall()
        print(len(sql_res))
        return sql_res
    except Exception:
        print('处理异常: ' + traceback.format_exc())
    finally:
        conn.close()


def pic_alter():
    video_path = '/home/resources/wwwroot/da_xiang/crypted/videos'
    # video_path = '/home/resources/resources/wwwroot/da_xiang/crypted/videos'
    res = con_sql()
    for info in res:
        video_id = str(info["id"])
        video_cover_url = info["pic_url"]
        if re.search(r'videos\/\w{6,32}\/covers\/cover-\d{1,}.webp\?resId=\w{2}', video_cover_url, re.I):
            reg_film_cover = re.compile(r"cover-\d{1,}.webp", re.I)
            res_film_cover = (re.findall(reg_film_cover, video_cover_url))[0]
        # try:
            # if os.path.exists(os.path.join(video_path, video_id, 'covers')):
            #     print('dir is exist')
            #     # 新文件夹路径
            #     new_dir_path = os.path.join(video_path, video_id, 'setcover')
            #     print(new_dir_path)
            #     # 新图片路径
            #     new_file_path = os.path.join(video_path, video_id, 'setcover/cover.webp')
            #     print(new_file_path)
            #     os.mkdir(new_dir_path)
            #     # 旧文件图片路径
            #     old_file_path = os.path.join(video_path, video_id, 'covers', res_film_cover)
            #     print(old_file_path)
            #     # 以二进制只读打开文件
            #     f_old_pic = open(old_file_path, 'rb', bufsize=1024 * 1024 * 5)
            #     pic_res = f_old_pic.read()
            #     f_new_pic = open(new_file_path, 'wb', bufsize=1024 * 1024 * 5)
            #     f_new_pic.write(pic_res)
            #
            #     f_old_pic.close()
            #     f_new_pic.close()
        # except Exception as e:
        #     print(e)


pic_alter()
