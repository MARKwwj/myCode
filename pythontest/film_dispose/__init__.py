# def process():
#     info = all_infos()
#     all_films = os.listdir(info.videos_path_server)
#
#     all_films_index = len(all_films)
#     one_part = int(all_films_index / 2)
#     list_one = []
#     list_another = []
#
#     for i in range(one_part):
#         list_one.append(all_films[i])
#     for i in range(one_part, all_films_index):
#         list_another.append(all_films[i])
#
#     process1 = multiprocessing.Process(target=run, args=(list_one, info))
#     process2 = multiprocessing.Process(target=run, args=(list_another, info))
#
#     process1.start()
#     process2.start()
#
#
# process()
import datetime
import os
import random
import time
import uuid
