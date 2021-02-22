# coding:utf8

from flask import Flask, render_template, session, request, flash, redirect, url_for
from flask_wtf import FlaskForm
from wtforms import StringField, PasswordField, SubmitField, IntegerField, DateField, BooleanField, ValidationError
from wtforms.validators import DataRequired, EqualTo
import config
from exts import db
import datetime

app = Flask(__name__)
app.config.from_object(config)
db.init_app(app)
app.config["SECRET_KEY"] = "12345678"


class User(db.Model):
    __tablename__ = 'user'
    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    username = db.Column(db.String(10), nullable=False)
    telephone = db.Column(db.String(11), nullable=False)
    address = db.Column(db.String(50), nullable=False)
    password = db.Column(db.String(20), nullable=False)
    password2 = db.Column(db.String(20))

    # def __repr__(self):
    #     return "用户信息: 姓名=%s" % self.用户名


class Provider(db.Model):
    __tablename__ = "provider"
    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    provider = db.Column(db.String(100), nullable=False)
    provide_time = db.Column(db.Date, nullable=False)

    # def __repr__(self):
    #     return "供货商信息: 公司名称=%s" % self.供货商


class Goods(db.Model):
    __tablename__ = "goods"
    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    title = db.Column(db.String(100), nullable=False)
    type = db.Column(db.String(100), nullable=False)
    num = db.Column(db.Integer, nullable=False)
    price = db.Column(db.String(10), nullable=False)
    produce_date = db.Column(db.Date, nullable=False)
    guarantee_time = db.Column(db.Integer, nullable=False)
    # 以下操作为数据表的属性关联
    provider_id = db.Column(db.Integer, db.ForeignKey('provider.id'))
    provider = db.relationship('Provider', backref=db.backref('goods'))
    # def __repr__(self):
    #     return "商品信息: 名称=%s" % self.名称


class RegisterForm(FlaskForm):
    # 自定义的注册表单模型
    username = StringField(u'用户名', validators=[DataRequired(u'用户名不能为空')])
    telephone = StringField(u'手机号', validators=[DataRequired(u'手机号不能为空')])
    address = StringField(u'邮箱', validators=[DataRequired(u'邮箱不能为空')])
    password = PasswordField(u'密码', validators=[DataRequired(u'密码不能为空')])
    password2 = PasswordField(u'确认密码', validators=[DataRequired(u'确认密码不能为空'), EqualTo("password", u'两次密码不一致')])
    submit = SubmitField(label=u'立即注册')

    # 自定义用户名验证器
    def validate_username(self, field):
        if User.query.filter_by(username=field.data).first():
            raise ValidationError('用户名已注册，请选用其它名称')

    # 自定义邮箱验证器
    def validate_address(self, field):
        if User.query.filter_by(address=field.data).first():
            raise ValidationError('该邮箱已注册使用，请选用其它邮箱')


class LoginForm(FlaskForm):
    username = StringField(label=u'用户名', validators=[DataRequired(u"用户名不能为空")])
    password = PasswordField(label=u'密码', validators=[DataRequired(u"密码不能为空")])
    submit = SubmitField(label=u'登录')


class AddForm(FlaskForm):
    id = IntegerField(label=u"供应商编号", validators=[DataRequired(u'供应商编号必填')])
    name = StringField(label=u"供货商", validators=[DataRequired(u'供货商必填')])
    type = StringField(label=u"类型", validators=[DataRequired(u'类型必填')])
    prv_time = DateField(label=u"供货时间", validators=[DataRequired(u'供货时间必填')])
    title = StringField(label=u"商品名称", validators=[DataRequired(u'名称必填')])
    pro_time = DateField(label=u"生产日期", validators=[DataRequired(u'生产日期必填')])
    guarantee_time = IntegerField(label=u"保质期", validators=[DataRequired(u"保质期必填")])
    num = IntegerField(label=u"数量", validators=[DataRequired(u"数量必填")])
    price = IntegerField(label=u"价格", validators=[DataRequired(u"价格必填")])
    submit = SubmitField(label=u'保存数据')


@app.route('/index')
def index():
    db.create_all()  # 创建数据库
    # 第一次运行可以开启一下demo字段，下次运行可以关闭，防止字段混乱
    goods1 = Goods(title='雪碧', price='3.0', num=250, type='饮料', produce_date='2018-9-25', guarantee_time=330,
                   provider_id=2)
    goods2 = Goods(title='可乐', price='3.0', num=200, type='饮料', produce_date='2018-9-26', guarantee_time=300,
                   provider_id=2)
    goods3 = Goods(title='牛奶', price='3.5', num=180, type='饮料', produce_date='2018-9-28', guarantee_time=45,
                   provider_id=2)
    goods4 = Goods(title='泡面', price='4.0', num=150, type='食品', produce_date='2018-10-8', guarantee_time=180,
                   provider_id=1)
    goods5 = Goods(title='火腿肠', price='1.5', num=320, type='食品', produce_date='2018-10-10', guarantee_time=150,
                   provider_id=1)
    goods6 = Goods(title='果汁', price='4.0', num=300, type='饮料', produce_date='2018-10-9', guarantee_time=200,
                   provider_id=1)
    goods7 = Goods(title='薯片', price='3.5', num=250, type='零食', produce_date='2018-8-20', guarantee_time=250,
                   provider_id=3)
    goods8 = Goods(title='饼干', price='2.5', num=235, type='零食', produce_date='2018-8-22', guarantee_time=240,
                   provider_id=3)

    company1 = Provider(provider='食品公司', provide_time='2018-10-15')
    company2 = Provider(provider='饮料公司', provide_time='2018-10-3')
    company3 = Provider(provider='零食公司', provide_time='2018-10-21')

    db.session.add_all([goods1,goods2,goods3,goods4,goods5,goods6,goods7,goods8])
    db.session.add_all([company1, company2, company3])
    # 如果需要在代码中创建用户则手动添加代码
    # db.session.drop(); 删除全部数据库
    db.session.commit()
    username = session.get("username", "请先登录")
    return render_template('index.html', username=username)


@app.route('/check/')  # 检测临期商品视图
def check():
    # 获取商品数据库里生产日期
    goods_time = Goods.produce_date
    days = Goods.guarantee_time
    # 预计到期时间=生产日期+保质期
    Pre_date = goods_time + datetime.timedelta(days=30)
    # 当前时间
    now = datetime.datetime.now()
    warn_date = Pre_date - datetime.timedelta(days=20)  # 在到期前20天警告
    goods_list = Goods.query.filter(Goods.produce_date.between(now, warn_date)).all()
    # goods_list = Goods.query.all()
    return render_template('check.html', goods=goods_list, now=now)


@app.route('/information/')  # 商品信息视图
def information():
    goods_list = Goods.query.all()
    provider_list = Provider.query.all()
    print(goods_list)
    print(provider_list)
    return render_template('information.html', goods=goods_list, providers=provider_list)


@app.route('/', methods=['GET', 'POST'])
def login():
    form = LoginForm()
    if form.validate_on_submit():
        username = form.username.data
        password = form.password.data
        user = User.query.first()
        if user:
            if username == user.username and password == user.password:
                flash("成功登录！")
                print("用户登录" + "  " + username, password)
            session['username'] = username
            return redirect(url_for("index"))
        else:
            print("未知账号")
    return render_template('login.html', form=form)


@app.route('/logout')
# 退出账号
def logout():
    if RegisterForm.username:
        session.pop('username')
        return redirect(url_for('login'))
    else:
        return render_template('login.html')


@app.route('/add/', methods=["GET", "POST"])  # 添加商品信息
def add():
    form = AddForm()
    if form.validate_on_submit():
        '''表单验证'''
        # 提取表单数据
        id = form.id.data
        name = form.name.data
        prv_time = form.prv_time.data
        title = form.title.data
        type = form.type.data
        pro_time = form.pro_time.data
        guarantee_time = form.guarantee_time.data
        num = form.num.data
        price = form.price.data
        # 保存数据库
        company = Provider(id=id, provider=name, provide_time=prv_time)
        goods = Goods(title=title, price=price, num=num, type=type, produce_date=pro_time,
                      guarantee_time=guarantee_time,
                      provider_id=id)
        db.session.add(company)
        db.session.add(goods)
        db.session.commit()
    return render_template('add.html', form=form)


@app.route('/register/', methods=["GET", "POST"])  # 用户信息注册界面
def register():
    # post请求的情况下，前端发送的数据会在构造form对象时，存放到对象里
    # 如果form中的数据完全满足所有的验证
    form = RegisterForm()
    if form.validate_on_submit():
        username = form.username.data
        telephone = form.telephone.data
        address = form.address.data
        password = form.password.data
        password2 = form.password2.data
        user = User(username=username, telephone=telephone, address=address, password=password)
        db.session.add(user)
        db.session.commit()
        flash('注册成功')
        print(username, telephone, address, password, password2)
        session['username'] = username
        return redirect(url_for("index"))
    return render_template('register.html', form=form)


@app.route('/search/')  # 查找页面
def search():
    q = request.args.get('q')
    goods1 = Goods.query.filter(Goods.title.contains(q))
    return render_template('search.html', goods1=goods1)


@app.route('/user_info/')  # 用户信息
def user_info():
    user_list = User.query.all()
    print(user_list)
    return render_template('user_info.html', users=user_list)


if __name__ == '__main__':
    app.run(
        debug=True,
        host='127.0.0.1',
        port=8000
    )
