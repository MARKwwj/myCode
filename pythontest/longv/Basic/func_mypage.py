import os

adb_Shell = "adb -s emulator-5554 shell  {cmdstr}"


def execute(cmd):
    """执行 adb 命令"""
    str = adb_Shell.format(cmdstr=cmd)
    print(str)
    os.system(str)


def bind_phone(driver):
    """绑定手机号"""
    # 切换到我的页面
    my_page = driver.find_element_by_accessibility_id("我的, tab, 5 of 5")
    my_page.click()
    # 点击绑定手机号
    bind_f = driver.find_element_by_accessibility_id("绑定手机号")
    bind_f.click()
    # 点击手机号输入框
    input_phone_number = driver.find_element_by_xpath(
        "/hierarchy/android.widget.FrameLayout/android.widget.LinearLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.view.View/android.view.View"
        "/android.view.View/android.view.View/android.widget.EditText[1]")
    input_phone_number.click()
    # 输入手机号
    execute("input text '18567503926'")
    security_code_get = driver.find_element_by_accessibility_id("获取验证码")
    security_code_get.click()
    # 点击验证码输入框
    input_security_code = driver.find_element_by_xpath(
        "/hierarchy/android.widget.FrameLayout/android.widget.LinearLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.widget.FrameLayout/android.view.View/android.view.View"
        "/android.view.View/android.view.View/android.widget.EditText[2]")
    input_security_code.click()
    # 输入错误的验证码
    execute("input text '888888'")


