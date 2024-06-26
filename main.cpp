#include "mainwindow.h"

#include <QApplication>
/******************************************************************************
 *
 * @file       main.cpp
 * @brief      main Function
 *
 * @author     Liao
 * @date       2024/06/26
 * @history
 *****************************************************************************/

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    MainWindow w;
    w.show();
    return a.exec();
}
