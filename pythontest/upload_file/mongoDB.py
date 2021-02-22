import os

import pymongo


def ConnMongo(collection):
    """
    连接mongoDB
    :param collection: 集合名称
    :return: mycol
    """
    myclient = pymongo.MongoClient("mongodb://64.64.241.63:7777/")
    mydb = myclient["magic"]
    mycol = mydb[collection]
    return mycol


def insert_novel(name, id):
    mycol = ConnMongo('novel')
    mydict = {
        "_id": id,
        "title": "{0}".format(name),
        "introduce": "这是一本小说这是一本小说这是一本小说这是一本小说这是一本小说",
        "picUrl": "server_upload/novel/995ec04b1cfc87f88ea382acf734d9d2.jpg",
        "freeNum": 3,
        "coin": 6,
        "type": "连载中",
        "status": "启用",
        "delFlag": "0",
        "popularity": 158801,
        "score": 1,
        "readCount": 46,
        "wordNum": 456789,
        "tags": ["穿越",
                 "历史",
                 "仙侠"],
        "createBy": "admin",
        "updateBy": "admin",
        "_class": "com.zsproject.framework.novel.domain.Novel",
        "author": "琉璃猫",
        "payCount": 13,
        "baseReadCount": 45,
        "machineId": "r5",
        "totalNum": 20,
        "tatalNum": 20,
        "chargeType": "付费",
        "searchCount": 1231
    }
    mycol.insert_one(mydict)


def insert_chapter_detail_mdb(file_path, file_name, id):
    file_all = os.listdir(file_path)
    file_all.sort(key=lambda x: int(x[:-4]))
    collention = "novel_chapter_details"
    mycol = ConnMongo(collention)
    chapter = 1
    for file in file_all:
        mydict = {
            "novelId": id,
            "chapterName": "第 {0} 章".format(chapter),
            "content": "chapterDetails/{0}/{1}".format(file_name, file),
            "sort": chapter,
            "_class": "com.zsproject.framework.novel.domain.NovelChapterDetails"
        }
        chapter = chapter + 1
        mycol.insert_one(mydict)


def cal_id():
    mycol = ConnMongo('novel')
    id_list = []
    for x in mycol.find({}, {"_id": 1}):
        id_list.append(x['_id'])
    id_list.sort()
    id = int(id_list[len(id_list) - 1] + 1)
    print(id)
    return id


