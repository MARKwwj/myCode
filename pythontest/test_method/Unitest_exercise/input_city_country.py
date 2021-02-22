import city_functions

print('请登记你的城市和国家！ 退出请按 q ')
while True:
    city = input('请输入你的城市：')
    if city == 'q':
        break
    country = input('请输入你的国家：')
    if city == 'q':
        break
    cc = city_functions.city_country(city, country)
    print(cc)
