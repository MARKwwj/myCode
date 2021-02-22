# This sample code uses the Appium python client
# pip install Appium-Python-Client
# Then you can paste this into a file and simply run with Python
from appium import webdriver
import os
from appium import webdriver
from appium.webdriver.common.touch_action import TouchAction

from toast import get_toast

adb_Shell = "adb -s D3H7N17A11006378 shell  {cmdstr}"


def execute(cmd):
    str = adb_Shell.format(cmdstr=cmd)
    print(str)
    os.system(str)


caps = {"platformName": "Android",
        "deviceName": "D3H7N17A11006378",
        "platformVersion": "10",
        "appPackage": "com.bshuai.tw",
        "appActivity": "com.bshuai.tw.MainActivity",
        "ensureWebviewsHavePages": True,
        'app': r'D:\desktop\app-debug.apk',
        'automationName': 'flutter'
        }

driver = webdriver.Remote("http://localhost:4723/wd/hub", caps)
driver.implicitly_wait(20)
el1 = driver.find_element_by_id("com.android.permissioncontroller:id/permission_allow_button")
el1.click()
el2 = driver.find_element_by_accessibility_id("确定")
el2.click()
el3 = driver.find_element_by_accessibility_id("取消")
el3.click()
el4 = driver.find_element_by_xpath(
    "/hierarchy/android.widget.FrameLayout/android.widget.LinearLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.view.View/android.view.View/android.view.View[1]/android.view.View/android.widget.ImageView[2]")
el4.click()
# el33 = driver.find_element_by_accessibility_id("我的, tab, 5 of 5")
# el33.click()
# el44 = driver.find_element_by_accessibility_id("绑定手机号")
# el44.click()
#
# el6 = driver.find_element_by_xpath(
#     "/hierarchy/android.widget.FrameLayout/android.widget.LinearLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.view.View/android.view.View/android"
#     ".view.View/android.view.View/android.widget.EditText[1]")
# el6.click()
# execute("input text '18567503926'")
# print(el6.text)
# assert el6.text == "18567503926, 请输入11位正确手机号"
#
# el7 = driver.find_element_by_accessibility_id("获取验证码")
# el7.click()
#
# toast = get_toast(driver, text='验证码已发送！')
# print(toast.text)
el11 = driver.find_element_by_accessibility_id("任务, tab, 4 of 5")
el11.click()
TouchAction(driver).tap(x=1223, y=1435).perform()
toast = get_toast(driver, text='复制成功')
print(toast)
driver.quit()
