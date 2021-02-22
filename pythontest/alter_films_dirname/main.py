import os
import pymysql
import traceback
from pymysql.cursors import DictCursor

videos_files_path = '/home/resources/wwwroot/da_xiang/crypted/videos'
# videos_files_path = '/home/resources/resources/wwwroot/da_xiang/crypted/videos'



def con_sql():
    try:
        conn = pymysql.connect(
            host="80.251.223.22",
            port=6033,
            user="root",
            passwd="ZsNice2020.",
            database="manage_db"
        )
        return conn
    except Exception:
        print('处理异常: ' + traceback.format_exc())


def get_video_info(conn):
    """获取 video_id 和 video_cover"""
    try:
        cmd = conn.cursor(DictCursor)
        sql = 'select id,video_url from video  where del_flag=0 order by id ASC'
        cmd.execute(sql)
        result = cmd.fetchall()
        return result
    except Exception as e:
        print(e)


def alter_dir_name(videos_info):
    for video in videos_info:
        video_dir_name = ((video['video_url']).split('/'))[1]
        video_id = str(video['id'])
        video_dir_path = os.path.join(videos_files_path, video_dir_name)
        print(video_dir_path)
        isExists = os.path.exists(video_dir_path)
        print(isExists)
        os.chdir(videos_files_path)
        if isExists:
            if video_id != video_dir_name:
                cmd_alter_video_dir_name = 'mv {0} {1}'.format(video_dir_name, video_id)
                print(cmd_alter_video_dir_name)
                os.system(cmd_alter_video_dir_name)
            else:
                print("此视频id 对应的文件夹已经存在！")
            """修改 setcover 内 图片名称"""
            set_cover_path = os.path.join(videos_files_path, video_dir_name, 'setcover')
            isExists_setcover = os.path.exists(set_cover_path)
            if isExists_setcover:
                os.chdir(set_cover_path)
                cmd_alter_cover_name = 'mv *.webp cover.webp'
                print(cmd_alter_cover_name)
                os.system(cmd_alter_cover_name)
        else:
            # os.system('cp -r backups {0}'.format(video_id))
            # print("复制 backups 文件夹！")
            print("该文件夹不在此资源服务器！")


conn = con_sql()
videos_info = get_video_info(conn=conn)
alter_dir_name(videos_info)
