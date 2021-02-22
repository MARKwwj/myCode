def replace(files, old_text, new_text):
    file_data = ""
    with open(files, 'r', encoding='utf-8') as file:
        for line in file:
            for text in old_text:
                if text in line:
                    line = line.replace(text, new_text)
            file_data += line
    with open(files, 'w', encoding='utf-8') as file:
        print(file_data)
        file.write(file_data)


if __name__ == '__main__':
    files = "d:\\desktop\\Env.dart"
    old_text = ['const currentEnv = Env.Release', 'const currentEnv = Env.Local', 'const currentEnv = Env.Local_Test']
    new_text = 'const currentEnv = Env.Test'
    replace(files, old_text, new_text)
