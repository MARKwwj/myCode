import unittest
from Python_exercise import name_function


class NamesTestCase(unittest.TestCase):
    """ 测试 name_function.py"""

    def test_first_last_name(self):
        """ 能够正确地处理像 Janis Joplin 这样的姓名吗？ """
        formatted_name = name_function.get_formatted_name('janis', 'joplin')
        self.assertEqual(formatted_name, 'Janis Joplin')


if __name__ == '__main__':
    unittest.main()
