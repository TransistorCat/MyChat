#include "mainwindow.h"
#include "./ui_mainwindow.h"


MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    //创建一个CentralWidget, 并将其设置为MainWindow的中心部件
    login_dialog_ = new LoginDialog();
    setCentralWidget(login_dialog_);
    login_dialog_->show();
}

MainWindow::~MainWindow()
{
    delete ui;
}
