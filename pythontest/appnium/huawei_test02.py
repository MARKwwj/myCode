# This sample code uses the Appium python client
# pip install Appium-Python-Client
# Then you can paste this into a file and simply run with Python
from toast import get_toast
from appium import webdriver

caps = {}
caps["platformName"] = "Android"
caps["deviceName"] = "D3H7N17A11006378"
caps["platformVersion"] = "10"
caps["appPackage"] = "com.vmall.client"
caps["appActivity"] = "com.vmall.client.splash.fragment.StartAdsActivity"
caps["ensureWebviewsHavePages"] = True

driver = webdriver.Remote("http://localhost:4723/wd/hub", caps)

el1 = driver.find_element_by_id("com.vmall.client:id/button_positive")
el1.click()
driver.implicitly_wait(20)
el2 = driver.find_element_by_id("com.android.permissioncontroller:id/permission_allow_button")
el2.click()
driver.back()
driver.back()
toast = get_toast(driver, text='再按一次退出')
print(toast.text)

