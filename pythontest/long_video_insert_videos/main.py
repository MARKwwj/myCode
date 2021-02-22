import json
import random
import re
import string
from datetime import datetime

import pymysql

i = 246


class videos_info:
    def __init__(self):
        self.video_title = self.get_name()
        self.video_intro = '这里是简介，这里是简介，这里是简介，这里是简介，这里是简介'
        video_num = random.randint(2, 8)
        self.video_cover = 'videos/{0}/setcover/cover.webp'.format(video_num)
        self.video_duration = random.randint(100, 5000)
        self.video_serial_num = '{0}-{1}'.format(random.choice(string.ascii_letters), random.randint(10, 99))
        self.video_tags, self.tags_id_list = self.get_tags()
        self.video_machine_id = 1
        self.video_url = 'videos/{0}/output.m3u8'.format(video_num)
        pay_type = [0, 1, 3]
        self.video_pay_type = (random.sample(pay_type, 1))[0]
        self.video_pay_coin = random.randint(5, 35)
        score = round(random.uniform(8, 9), 1)
        self.video_score = score
        self.video_play_count = random.randint(12999, 999999)
        self.video_praise_count = random.randint(12999, 999999)
        self.video_favorite_count = random.randint(12999, 999999)
        self.video_status = 1
        self.video_create_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        self.video_update_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        self.video_creator = '1002'
        global i
        self.video_id = i
        i = i + 1
        self.category_id = random.randint(1, 5)
        self.classify_id = self.get_classify_id(self.category_id)

    def get_classify_id(self, category_id):
        r = []
        if category_id == 1:
            r = [55, 60]
        elif category_id == 2:
            r = [56, 61]
        elif category_id == 3:
            r = [57, 62]
        elif category_id == 4:
            r = [58, 63]
        elif category_id == 5:
            r = [59, 64]
        classify_id = (random.sample(r, 1))[0]
        return classify_id

    def get_name(self):
        name_list = []
        with open("title.txt", 'r', encoding='utf-8') as fr:
            for f in fr.readlines():
                name_list.append(f.replace('\n', ''))
        name = (random.sample(name_list, 1))[0]
        return name

    def get_tags(self):
        tag_list = []
        with open('tags.txt', 'r', encoding='utf-8') as fr:
            for f in fr.readlines():
                tag_list.append(f.replace('\n', ''))
        tags = random.sample(tag_list, 3)
        tags = '[{0}, {1}, {2}]'.format(tags[0], tags[1], tags[2])
        tags_json = json.dumps(tags, ensure_ascii=False)
        tags_id_list = re.findall(r'tagId\\": (\w{1,3})', tags_json)
        return tags_json, tags_id_list


def main():
    """插入数据库"""
    conn = pymysql.connect(
        host='199.180.114.169',
        port=6033,
        user='root',
        password='ZsNice2020.',
        database='res_video_db',
    )
    cmd = conn.cursor()
    for count in range(1):
        info = videos_info()
        videos_insert_sql = "insert into long_video_video_info(" \
                            "video_title,video_intro,video_cover," \
                            "video_duration,video_serial_num,video_tags," \
                            "video_machine_id,video_url,video_pay_type," \
                            "video_pay_coin,video_score,video_play_count," \
                            "video_praise_count,video_favorite_count,video_status," \
                            "video_create_time,video_update_time,video_creator,video_id) values(" \
                            "'{0}','{1}','{2}',{3},'{4}',{5}," \
                            "{6},'{7}',{8},{9},{10},{11}," \
                            "{12},{13},{14},'{15}','{16}','{17}',{18})".format(
                                info.video_title, info.video_intro, info.video_cover,
                                info.video_duration, info.video_serial_num, info.video_tags,
                                info.video_machine_id, info.video_url, info.video_pay_type,
                                info.video_pay_coin, info.video_score, info.video_play_count,
                                info.video_praise_count, info.video_favorite_count, info.video_status,
                                info.video_create_time, info.video_update_time, info.video_creator, info.video_id)
        print(videos_insert_sql)
        cmd.execute(videos_insert_sql)
        for tags_id in info.tags_id_list:
            tags_insert_sql1 = "INSERT INTO long_video_tag_relation(tag_id, video_id) VALUES ({0}, {1})".format(tags_id, info.video_id)
            cmd.execute(tags_insert_sql1)

        category_insert_sql = "INSERT INTO long_video_category_relation(category_id, video_id) VALUES ({0}, {1})".format(info.category_id, info.video_id)
        cmd.execute(category_insert_sql)

        classify_insert_sql = "INSERT INTO long_video_classify_relation(classify_id, video_id) VALUES ({0}, {1})".format(info.classify_id, info.video_id)
        cmd.execute(classify_insert_sql)
    conn.commit()


main()
