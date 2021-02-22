import unittest
from appium import webdriver
from appium_flutter_finder.flutter_finder import FlutterElement, FlutterFinder


class FlutterTest(unittest.TestCase):
    def setUp(self):
        desired_caps = {
            'platformName': 'Android',
            'platformVersion': '10',
            'deviceName': 'D3H7N17A11006378',
            'app': r'D:\desktop\app-debug.apk',
            'automationName': 'flutter'
        }
        self.driver = webdriver.Remote('http://127.0.0.1:4723/wd/hub', desired_caps)
        self.finder = FlutterFinder()

    def test_flutter(self):
        text_finder = self.finder.by_value_key("counter")
        button_finder = self.finder.by_value_key("increment")
        text_element = FlutterElement(self.driver, text_finder)
        button_element = FlutterElement(self.driver, button_finder)
        button_element.click()
        button_element.click()
        self.assertEqual('2', text_element.text)

    def tearDown(self):
        self.driver.quit()


if __name__ == '__main__':
    unittest.main()
