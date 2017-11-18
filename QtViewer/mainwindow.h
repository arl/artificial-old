#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = 0);
    ~MainWindow();

public slots:
    void watchEvolution();
    void showGenerationImage(int value);
    void showGenerationData(int value);

private:

    QPixmap loadPixmap(QString fileName);

    Ui::MainWindow *ui;
    bool m_dbOpened;
    QString m_dir;
};

#endif // MAINWINDOW_H
