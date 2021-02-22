import os
import zipfile


def zip_uncompress(file_path):
    zip_file_name = os.path.basename(file_path)
    file_name_list = zip_file_name.split('.')
    # 文件名称 去掉后缀的
    file_name = file_name_list[0]
    # 文件后缀
    file_suffix = file_name_list[len(file_name_list) - 1]
    # 判断如果文件后缀不是zip 就修改为zip
    if file_suffix != 'zip':
        new_file_path = file_path.replace(file_suffix, 'zip')
        os.rename(file_path, new_file_path)
        file_path = new_file_path
    print(file_path)
    file = zipfile.ZipFile(file_path)
    # 解压到哪个目录
    file.extractall()
    file.close()


if __name__ == '__main__':
    path = "D:\\desktop\\222\\v\\2439.txt"
    zip_uncompress(path)
