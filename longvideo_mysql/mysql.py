import hashlib

import pymysql
from pymysql.cursors import DictCursor

from parameters import all_infos


def conn_mysql(info):
    """连接数据库"""
    conn = pymysql.connect(
        host=info.host,
        port=info.port,
        user=info.user,
        passwd=info.password,
        database=info.database,
    )
    cmd = conn.cursor(DictCursor)
    return cmd, conn


def close_mysql(cmd, conn):
    """关闭连接"""
    conn.commit()
    cmd.close()
    conn.close()


def select_category_id(info, category_name):
    """视频类别"""
    sql = "select category_id from long_video_category_info where category_name ='{0}'".format(category_name)
    cmd, conn = conn_mysql(info)
    cmd.execute(sql)
    result = (cmd.fetchall())[0]["category_id"]
    info.category_id = result
    close_mysql(cmd, conn)


def select_classify_id(info, classify_name):
    """视频分类"""
    sql = "select classify_id from long_video_classify_info where classify_name ='{0}'".format(classify_name)
    cmd, conn = conn_mysql(info)
    cmd.execute(sql)
    result = (cmd.fetchall())[0]["classify_id"]
    info.classify_id = result
    close_mysql(cmd, conn)


def select_tag_info(info):
    """视频标签"""
    cmd, conn = conn_mysql(info)
    tags_new_list = []
    tags_new_id_list = []
    if info.video_tags:
        for tag in info.video_tags:
            sql = "select tag_id,tag_name from long_video_tag_info where tag_name ='{0}'".format(tag)
            cmd.execute(sql)
            result = cmd.fetchall()
            # 如果标签不存在
            if not result:
                if tag != "":
                    sql = "INSERT INTO long_video_tag_info (tag_name,parent_id) VALUES('{0}', {1} ) ".format(tag, info.tag_parent_id)
                    cmd.execute(sql)
                    conn.commit()

        for tag in info.video_tags:
            sql = "select tag_id,tag_name from long_video_tag_info where tag_name ='{0}'".format(tag)
            cmd.execute(sql)
            result = cmd.fetchall()
            if result:
                tags_new_list.append((list(result))[0])
                tags_new_id_list.append(result[0]["tag_id"])
        info.video_tags = (((str(tags_new_list)).replace("'", "\"")).replace("tag_id", "tagId")).replace("tag_name", "tagName")
        info.video_tags_id = tags_new_id_list
    close_mysql(cmd, conn)


def insert_videos_info(info):
    """插入视频信息"""
    cmd, conn = conn_mysql(info)
    # sql = "insert into long_video_video_info (" \
    #       "video_title,video_intro,video_cover_quantity," \
    #       "video_duration,video_tags,video_pay_type," \
    #       "video_pay_coin,video_score,video_play_count," \
    #       "video_praise_count,video_favorite_count,video_status," \
    #       "video_create_time,video_update_time,video_byte_size," \
    #       "video_creator,video_machine_id,video_machine_name) values(" \
    #       "'{0}','{1}',{2}," \
    #       "{3},'{4}',{5}," \
    #       "{6},{7},{8}," \
    #       "{9},{10},{11}," \
    #       "'{12}','{13}',{14}," \
    #       "'{15}',{16},'{17}')".format(info.video_title, info.video_intro, info.video_cover_quantity,
    #                                    info.video_duration, info.video_tags, info.video_pay_type,
    #                                    info.video_pay_coin, info.video_score, info.video_play_count,
    #                                    info.video_praise_count, info.video_favorite_count, info.video_status,
    #                                    info.video_create_time, info.video_update_time, info.video_byte_size,
    #                                    info.video_creator, info.video_machine_id, info.video_machine_name)
    sql = """
        INSERT INTO long_video_video_info (
        video_title,
        video_intro,
        video_cover_quantity,
        video_duration,
        video_tags,
        video_pay_type,
        video_pay_coin,
        video_score,
        video_play_count,
        video_praise_count,
        video_favorite_count,
        video_status,
        video_create_time,
        video_update_time,
        video_byte_size,
        video_creator,
        video_machine_id,
        video_machine_name,
        video_total_score,
        video_score_total_people
    )
    VALUES(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) 
    """
    cmd.execute(sql, (info.video_title, info.video_intro, info.video_cover_quantity,
                      info.video_duration, info.video_tags, info.video_pay_type,
                      info.video_pay_coin, info.video_score, info.video_play_count,
                      info.video_praise_count, info.video_favorite_count, info.video_status,
                      info.video_create_time, info.video_update_time, info.video_byte_size,
                      info.video_creator, info.video_machine_id, info.video_machine_name,
                      info.video_total_score, info.video_score_total_people))
    video_id = conn.insert_id()
    print("video_id:", video_id)
    info.video_id = video_id
    close_mysql(cmd, conn)


def insert_category_relation(info):
    """插入类别关系表"""
    cmd, conn = conn_mysql(info)
    sql = "insert into long_video_category_relation (category_id,video_id) VALUES ({0}, {1})".format(info.category_id, info.video_id)
    cmd.execute(sql)
    close_mysql(cmd, conn)


def insert_classify_relation(info):
    """插入分类关系表"""
    cmd, conn = conn_mysql(info)
    sql = "insert into long_video_classify_relation (classify_id,video_id) VALUES ({0}, {1})".format(info.classify_id, info.video_id)
    cmd.execute(sql)
    close_mysql(cmd, conn)


def insert_tag_relation(info):
    """插入标签关系表"""
    cmd, conn = conn_mysql(info)
    for id in info.video_tags_id:
        sql = "insert into long_video_tag_relation (tag_id,video_id) VALUES ({0}, {1})".format(id, info.video_id)
        cmd.execute(sql)
    close_mysql(cmd, conn)


def select_video_title(info):
    """获取数据库中所有视频标题"""
    cmd, conn = conn_mysql(info)
    sql = "select video_title from long_video_video_info"
    cmd.execute(sql)
    result = cmd.fetchall()
    dict_videos_title = {}
    with open('title_record.txt', 'w+', encoding='utf-8') as fr:
        dict_videos_title = eval("{" + ((fr.read()).replace('{', '')).replace('}', ',') + "}")
        fr.close()
    for dict_title in result:
        title = str(dict_title["video_title"]).strip()
        key = hashlib.md5(title.encode(encoding='utf-8')).hexdigest()
        dict_videos_title[key] = title
    with open('title_record.txt', 'w', encoding='utf-8') as fw:
        fw.write(str(dict_videos_title))
    print("title_record.txt 已生成！")
    close_mysql(cmd, conn)


select_video_title(info=all_infos())
