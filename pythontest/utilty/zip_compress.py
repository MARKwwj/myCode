import os
import zipfile


def zip_txt(file_path):
    # file_path = 'D:\\desktop\\秦时世界里的主神余孽'
    # 切换到需要压缩的文件 的路径 （压缩后的文件会生成到此目录下）
    os.chdir(file_path)
    # zipfile.ZipFile(filename,mode)  filename 是压缩后的名字
    z = zipfile.ZipFile('uncrypt.zip', mode="w", compression=zipfile.ZIP_DEFLATED)
    # f 需要压缩的文件
    z.write("app.json")
    z.write("conf")
    z.write("zip")
    z.close()


zip_txt("D:\\desktop\\uncrypt (12)")
