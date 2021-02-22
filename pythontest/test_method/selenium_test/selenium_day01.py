import selenium
from selenium import webdriver

bro = webdriver.Chrome()
bro.get("https://www.baidu.com/")
bro.find_element_by_id("kw")
