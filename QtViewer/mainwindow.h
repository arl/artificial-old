#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QPixmap>
#include <QLocalSocket>

namespace Ui {
class MainWindow;
}

// forward declarations
class QLabel;

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();

    void resizeEvent(QResizeEvent * evt);

public slots:
    void followEvolution();
    void showGenerationImage(int value);
    void showGenerationData(int value);

    void onSockError(QLocalSocket::LocalSocketError socketError);
    void onSockConnected();
    void onSockReadyRead();

private:

    QPixmap loadPixmap(QString fileName);
    void closeDatabase();
    void scaleQLabelPixmap(QLabel* lbl);

    Ui::MainWindow *ui;
    bool m_dbOpened;
    QString m_dir;
    QLocalSocket m_sock;

    QPixmap m_curPixmap;
    QPixmap m_refPixmap;
};

#endif // MAINWINDOW_H
