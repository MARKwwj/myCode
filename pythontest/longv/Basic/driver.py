from appium import webdriver


def driver_init():
    """
    设备参数
    :return: driver
    """
    caps = {
        "platformName": "Android",
        "deviceName": "emulator-5554",
        "platformVersion": "7.1.2",
        "appPackage": "com.bshuai.tw",
        "appActivity": "com.bshuai.tw.MainActivity",
        'newCommandTimeout': "30000",
    }

    driver = webdriver.Remote("http://localhost:4723/wd/hub", caps)
    return driver
