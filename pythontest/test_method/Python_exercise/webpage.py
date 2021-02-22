import webbrowser
import sys
from pyQt5.QtCore import QUrl
from pyQt5.QtWebEngineWindgets import QWebEnginePage, QWebEngineView
from pyQt5.QtWidgets import QApplication, Qwidget, QPushBUtton, \
    QDesktopWidget, QLabel, QGridLayout
import webbrowser, sys

# webbrowser.open('http://baidu.com/')
app = QApplication(sys.argv)
browser = QWebEngineView()
browser.load(QUrl("http://baidu.com/"))
browser.show()
app.exec_()


class Ui_MaimWind(QWindget):
    item_name = "PyQt打开外部链接"

    def __init__(self):
        super().__init__()
        self.initUI()

    def initUI(self):
        self.tips_1 = QLabel("网站：<a href='http://baidu.com/'>http://baidu.com/</a>");
        self.tips_1.setOpenEXternalLinks(True)
        self.btn_webbrowser = QPushBUtton('webbrowser效果', self)
        self.btn_webbrowser.clicked.connect(self.btn_webbrowser_Clicked)
        grid = QGridLayout()
        grid.setSpacing(10)
        grid.addWidget(self.btn_webbrowser, 1, 0)
        grid.addWidget(self.tips_1, 2, 0)
        self.setLayout(grid)
        self.resize(250, 150)
        self.setMinmumSize(266, 304);
        self.setMaximumSize(266, 304);
        self.center()
        self.setWindowTitle(self.item_name)
        self.show()

    def btn_webbrowser_Clicked(self):
        webbrowser.open('http://www.baidu.com/')

    def center(self):
        qr = self.frameGeometry()
        cp = QDesktopWidget().availableGeometty().center()
        qr.moveCenter(cp)
        self.move(qr.topLeft())


if __name__ == "__main__":
    app = QApplication(sys.argv)
    a = Ui_MaimWindow()
    sys.exit(app.exec_())
