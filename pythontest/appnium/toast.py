# coding:utf-8
from appium import webdriver

from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC


def get_toast(driver, text=None, timeout=5, poll_frequency=0.01):
    """
    get toast
    :param driver: driver
    :param text: toast text
    :param timeout: Number of seconds before timing out, By default, it is 5 second.
    :param poll_frequency: sleep interval between calls, By default, it is 0.5 second.
    :return: toast
    """
    print('传入:' + text)
    if text:
        toast_loc = (".//*[contains(@text, '%s')]" % text)
    else:
        toast_loc = "//*[@class='android.widget.Toast']"

    try:
        WebDriverWait(driver, timeout, poll_frequency).until(EC.presence_of_element_located(('xpath', toast_loc)))
        toast_elm = driver.find_element_by_xpath(toast_loc)

        return toast_elm

    except:
        return "Toast not found"
