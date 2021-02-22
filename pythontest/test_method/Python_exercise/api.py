# python练习：发送请求
# import requests
#
# r = requests.get('https://www.baidu.com/')
# # r.status_code 返回状态码
# code = r.status_code
# if code == 200:
#     print("success")
#
# # print(r.content)
# # 使用str(string[, encoding])对数组进行转换
# html = str(r.content, 'utf-8')
# print(html)


# python练习：检查字符串是否ip
# def is_ip(ip):
#     num_list = ip.split(".")
#     for num in num_list:
#         if not num.isdigit() or not 0 <= int(num) <= 255:
#             return False
#     return True
#
#
# print(is_ip("101.1.0.201"))
#
# a = "s"
# print(a.isdigit())

# def bubble_sort(m_list):
#     for i in range(0, len(m_list)):
#         for j in range(i + 1, len(m_list)):
#             if m_list[i] > m_list[j]:
#                 temp = m_list[i]
#                 m_list[i] = m_list[j]
#                 m_list[j] = temp
#                 print()
#     print(m_list)
#
#
# bubble_sort([1, 32, 3, 4, 7, 1, 34, 23, 45])
