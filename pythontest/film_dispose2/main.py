import os
import shutil
import logging

fmt = '[%(asctime)s] %(levelname)s [%(funcName)s: %(filename)s, %(lineno)d] %(message)s'
logging.basicConfig(level=logging.DEBUG,  # log level
                    filename='film_dis.txt',  # log名字
                    format=fmt,  # log格式
                    datefmt='%Y-%m-%d %H:%M:%S',  # 日期格式
                    filemode='a')  # 追加模式


def ts_change_to_mp4(videos_path):
    file_dir_list = os.listdir(videos_path)
    for dir in file_dir_list:
        dir_path = os.path.join(videos_path, dir)
        # 删除 文件名字不是 视频id 的旧资源
        if len(dir) > 10:
            shutil.rmtree(dir_path)
        is_dir = os.path.isdir(dir_path)
        if is_dir:
            os.chdir(dir_path)
            try:
                os.stat("output.m3u8")
                res = os.system("ffmpeg -i output.m3u8 -y -acodec copy -vcodec copy -f mp4 output.mp4")
                if res == 0:
                    logging.info(dir_path + "-处理完成！")
                else:
                    logging.info(dir_path + "-处理失败！")
            except IOError:
                logging.info("{0}当前文件夹中不存在 m3u8".format(dir_path))
                os.chdir(videos_path)
                new_dir_path = os.path.join(videos_path, dir + "-m3u8IsNotExists")
                os.rename(dir_path, new_dir_path)


# /home/datadrive/resources/res4/da_xiang/crypted/videos
def main():
    for i in range(1, 7):
        if i == 3:
            continue
        videos_path = "/home/datadrive/resources/res{0}/da_xiang/crypted/videos_2".format(str(i))
        ts_change_to_mp4(videos_path)


main()