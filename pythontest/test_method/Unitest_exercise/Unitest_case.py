import unittest
import city_functions


class city_country_mytestcase(unittest.TestCase):
    """测试 city_function.py """

    def test_citycountry(self):
        """测试"""
        cc = city_functions.city_country('sss', 'ss')
        self.assertEqual(cc, '上海,china')


if __name__ == '__main__':
    unittest.main()
