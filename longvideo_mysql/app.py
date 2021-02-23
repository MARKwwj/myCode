import os

import flask
from flask import Flask, jsonify

from film_disp_func import run
from mysql import *
from parameters import all_infos

app = Flask(__name__)


@app.route('/videos', methods=["post"])
def videos():
    info = all_infos()
    if flask.request.method == "POST":
        result = flask.request.json
        info.video_title = result["video_title"]
        if result["only_query"]:
            # 防插入重复视频校验
            mark = repeat_check(info)
            print("not_repeat：", mark)
            return_json_repeat_msg = {
                "not_repeat": mark,
            }
            return jsonify(return_json_repeat_msg)
        else:
            print(result)
            # 防插入重复视频校验
            mark = repeat_check(info)
            if not mark:
                print("not_repeat：", mark)
                return_json_repeat_msg = {
                    "not_repeat": mark,
                }
                return jsonify(return_json_repeat_msg)
            info.video_intro = result["video_intro"]
            if info.video_intro == "":
                info.video_intro = result["video_title"]
            info.video_duration = result["video_duration"]
            info.video_byte_size = result["video_byte_size"]
            info.video_pay_coin = cal_charge(result["video_duration"])
            tags_string = result["video_tags"]
            video_classify_string = result["video_classify"]
            if tags_string != "":
                # 标签string  转 list
                tags_list = tags_string.split(',')
                info.video_tags = tags_list
            select_tag_info(info)
            if video_classify_string != "":
                # 类别-分类 拆分字符串
                category_name = (video_classify_string.split('-'))[0]
                classify_name = (video_classify_string.split('-'))[1]
            else:
                return_json_classify_msg = {
                    "msg": "video_classify 不能为空！"
                }
                return jsonify(return_json_classify_msg)
            select_category_id(info, category_name)
            select_classify_id(info, classify_name)
            # 数据库插入数据
            insert_videos_info(info)
            insert_category_relation(info)
            insert_classify_relation(info)
            insert_tag_relation(info)
            # 名字写入文本
            write_video(info)
            return_json_success_msg = {
                "video_id": info.video_id,
            }
            return jsonify(return_json_success_msg)


@app.route('/videoAlter', methods=["post"])
def videoAlter():
    info = all_infos()
    if flask.request.method == "POST":
        result = flask.request.json
        if result["video_id"] == "":
            return_json_success_msg = {
                "msg": "视频Id是空",
            }
            return jsonify(return_json_success_msg)
        video_id = str(result["video_id"])
        run(info.cur_video_path, info.res_video_path, video_id)
        return_json_success_msg = {
            "msg": "视频:{0} 处理结束!".format(video_id),
        }
        return jsonify(return_json_success_msg)


def cal_charge(video_times):
    """计算视频金币"""
    charge = 0
    if video_times < 60 * 30:
        charge = 10
    elif video_times < 60 * 60:
        charge = 20
    elif video_times < 60 * 120:
        charge = 30
    else:
        charge = 40
    return charge


def repeat_check(info):
    mark = True
    title = info.video_title
    key = hashlib.md5(title.encode(encoding='utf-8')).hexdigest()
    with open('title_record.txt', 'r', encoding='utf-8') as fr:
        info.dict_videos_title = eval("{" + ((fr.read()).replace('{', '')).replace('}', ',') + "}")
        fr.close()
    if info.dict_videos_title.__contains__(key):
        mark = False
    return mark


def write_video(info):
    title = info.video_title
    key = hashlib.md5(title.encode(encoding='utf-8')).hexdigest()
    info.dict_videos_title[key] = title
    with open('title_record.txt', 'w', encoding='utf-8') as fw:
        fw.write(str(info.dict_videos_title))
        fw.close()


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8001, debug=True)
