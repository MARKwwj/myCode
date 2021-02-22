import re


def replace():
    regex = re.compile(r'    title:(.*[\n].*){14}')
    file = ''
    txt = ''
    with open('D:\\desktop\\mdb1.txt', 'r', encoding='utf-8', buffering=1024 * 1024 * 10) as fr:
        txt = fr.read()

    file = re.sub(regex, 'chapters: [', txt)

    with open('D:\\desktop\\mdb1.txt', 'w', encoding='utf-8', buffering=1024 * 1024 * 10) as fw:
        fw.write(file)


replace()
