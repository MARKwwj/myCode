# coding=utf-8
import os
from appium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.wait import WebDriverWait

import toast

caps = {
    "platformName": "Android",
    "deviceName": "D3H7N17A11006378",
    "platformVersion": "10",
    "appPackage": "com.bshuai.tw",
    "appActivity": "com.bshuai.tw.MainActivity",
    'newCommandTimeout': "30000",
}

driver = webdriver.Remote("http://localhost:4723/wd/hub", caps)

adb_Shell = "adb -s emulator-5554 shell  {cmdstr}"


def execute(cmd):
    str = adb_Shell.format(cmdstr=cmd)
    print(str)
    os.system(str)


driver.implicitly_wait(50)
el1 = driver.find_element_by_id("com.android.packageinstaller:id/permission_allow_button")
el1.click()
el2 = driver.find_element_by_accessibility_id("确定")
el2.click()
el3 = driver.find_element_by_accessibility_id("我的, tab, 5 of 5")
el3.click()
el4 = driver.find_element_by_accessibility_id("取消")
el4.click()
el5 = driver.find_element_by_accessibility_id("绑定手机号")
el5.click()

el6 = driver.find_element_by_xpath(
    "/hierarchy/android.widget.FrameLayout/android.widget.LinearLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.view.View/android.view.View/android"
    ".view.View/android.view.View/android.widget.EditText[1]")
el6.click()
execute("input text '18567503926'")
print(el6.text)
assert el6.text == "18567503926, 请输入11位正确手机号"
# el7 = driver.find_element_by_accessibility_id("获取验证码")
# el7.click()

el8 = driver.find_element_by_xpath(
    "/hierarchy/android.widget.FrameLayout/android.widget.LinearLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.view.View/android.view.View/android"
    ".view.View/android.view.View/android.widget.EditText[2]")
el8.click()
execute("input text '888888'")


def load(driver):
    print("waitSpecificToast : load")
    toast_text = driver.find_element(By.XPATH, "//*[@class='android.widget.Toast']").text
    if toast_text.find('expected_text') != -1:
        return True
    else:
        return False


WebDriverWait(driver, 60).until(load)
res = load(driver)
print(res)
