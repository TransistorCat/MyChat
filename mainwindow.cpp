#include "mainwindow.h"
#include "./ui_mainwindow.h"


MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    //创建一个CentralWidget, 并将其设置为MainWindow的中心部件
    login_dialog_ = new LoginDialog(this);
    setCentralWidget(login_dialog_);
    // login_dialog_->show();
    //创建和注册消息的链接
    connect(login_dialog_, &LoginDialog::switchRegister,
            this, &MainWindow::SlotSwitchReg);
    register_dialog_=new RegisterDialog(this);

    login_dialog_->setWindowFlags(Qt::CustomizeWindowHint|Qt::FramelessWindowHint);
    register_dialog_->setWindowFlags(Qt::CustomizeWindowHint|Qt::FramelessWindowHint);
    register_dialog_->hide();
}

MainWindow::~MainWindow()
{
    delete ui;
    //关闭时由于设置成CentralWidget，还在Qt事件循环所以删除会崩溃
    // if(login_dialog_){
    //     delete login_dialog_;
    //     login_dialog_=nullptr;
    // }
    // if(register_dialog_){
    //     delete register_dialog_;
    //     register_dialog_=nullptr;
    // }
}

void MainWindow::SlotSwitchReg(){
    setCentralWidget(register_dialog_);
    login_dialog_->hide();
    register_dialog_->show();
}
