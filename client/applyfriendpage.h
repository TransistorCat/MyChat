#ifndef APPLYFRIENDPAGE_H
#define APPLYFRIENDPAGE_H

#include <QWidget>

namespace Ui {
class ApplyFriendPage;
}

class ApplyFriendPage : public QWidget
{
    Q_OBJECT

public:
    explicit ApplyFriendPage(QWidget *parent = nullptr);
    ~ApplyFriendPage();
protected:
    void paintEvent(QPaintEvent *event);
private:
    Ui::ApplyFriendPage *ui;
signals:
    void sig_show_search(bool);
};

#endif // APPLYFRIENDPAGE_H
