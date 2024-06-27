#include "mainwindow.h"
#include <QFile>
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

    QFile qss(":/style/stylesheet.qss");

    if( qss.open(QFile::ReadOnly))
    {
        qDebug("open success");
        QString style = QLatin1String(qss.readAll());
        a.setStyleSheet(style);
        qss.close();
    }else{
        qDebug("Open failed");
    }

    MainWindow w;
    w.show();

    return a.exec();
}
