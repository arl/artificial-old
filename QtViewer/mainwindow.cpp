#include <QDebug>
#include <QFileDialog>
#include <QSqlDatabase>
#include <QFileInfo>
#include <QPixmap>
#include <QSqlQuery>
#include <QIODevice>
#include <QLabel>

#include <algorithm>

#include "mainwindow.h"
#include "ui_mainwindow.h"

constexpr int GenerationFrequency = 100;
const QString socketName = "generation.sock";

int maxGenerationNumber()
{
    auto q = QSqlQuery("SELECT MAX(gen_number) FROM generations;", QSqlDatabase::database());
    int num = -1;
    if (q.next())
    {
        num = q.value(0).toInt();
    }
    return num;
}

QPixmap scalePixmapAspectRatio(QPixmap& pixmap, int w, int h)
{
    return pixmap.scaled(std::min(w, h), std::min(w, h), Qt::KeepAspectRatio);
}

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);

    ui->refImage->setMinimumSize(1, 1);
    ui->curImage->setMinimumSize(1, 1);

    m_dbOpened = false;
    QObject::connect(&m_sock, static_cast<void(QLocalSocket::*)(QLocalSocket::LocalSocketError)>(&QLocalSocket::error),
        this, &MainWindow::onSockError);
    QObject::connect(&m_sock, &QLocalSocket::connected, this, &MainWindow::onSockConnected);
    QObject::connect(&m_sock, &QLocalSocket::readyRead, this, &MainWindow::onSockReadyRead);
}

MainWindow::~MainWindow()
{
    closeDatabase();
    delete ui;
}

void MainWindow::resizeEvent(QResizeEvent * evt)
{
    ui->refImage->setPixmap(scalePixmapAspectRatio(m_refPixmap, ui->refImage->width(), ui->refImage->height()));
    ui->curImage->setPixmap(scalePixmapAspectRatio(m_curPixmap, ui->curImage->width(), ui->curImage->height()));
    QMainWindow::resizeEvent(evt);
}

void MainWindow::scaleQLabelPixmap(QLabel* lbl)
{
    const QPixmap * pixmap = lbl->pixmap();
    if (pixmap)
    {
        // get label dimensions
        int w = lbl->width();
        int h = lbl->height();

        // set a scaled pixmap to a w x h window keeping its aspect ratio
        lbl->setPixmap(pixmap->scaled(w, h, Qt::KeepAspectRatio));
    }
}

void MainWindow::closeDatabase()
{
    if (m_dbOpened)
    {
        auto db = QSqlDatabase::database();
        db.close();
    }
}

void MainWindow::followEvolution()
{
    // user must choose a folder containing a `generations.db` file
    auto fileName = QFileDialog::getOpenFileName(this,
        tr("Open evolution database"), "", tr("evolution db files(evolution.db)"));
    qInfo() << fileName;

    // open database
    closeDatabase();
    auto db = QSqlDatabase::addDatabase("QSQLITE");
    db.setDatabaseName(fileName);
    m_dbOpened = db.open();
    if (m_dbOpened)
    {
        qInfo() << "opened database";
        m_dir = QFileInfo(fileName).absolutePath();

        // configure slider
        ui->generationSlider->setMinimum(0);
        ui->generationSlider->setTickInterval(1);
        ui->generationSlider->setEnabled(true);
        ui->generationSlider->setSingleStep(1);
        ui->generationSlider->setValue(0);
        ui->generationSlider->setMaximum(maxGenerationNumber() / GenerationFrequency);

        // show generated images and data at generation 0
        showGenerationImage(0);
        showGenerationData(0);

        // show reference image
        m_refPixmap = loadPixmap("_ref.png");
        ui->refImage->setPixmap(scalePixmapAspectRatio(m_refPixmap, ui->refImage->width(), ui->refImage->height()));

        // open socket signaling arrival of new generation data
        m_sock.connectToServer(m_dir + QDir::separator() + socketName, QIODevice::ReadOnly);
    }
}

void MainWindow::showGenerationImage(int value)
{
    auto generation = value * GenerationFrequency;
    if (m_dbOpened)
    {
        m_curPixmap = loadPixmap(QString::number(generation) + ".png");
        ui->curImage->setPixmap(scalePixmapAspectRatio(m_curPixmap, ui->curImage->width(), ui->curImage->height()));
    }
}

QPixmap MainWindow::loadPixmap(QString fileName)
{
    auto imgPath = m_dir + QDir::separator() + fileName;
    if (QFileInfo(imgPath).exists())
        return QPixmap(imgPath);

    qInfo() << imgPath << " doesn't exist";
    return QPixmap();
}

void MainWindow::showGenerationData(int value)
{
    auto generation = value * GenerationFrequency;
    if (m_dbOpened)
    {
        auto q = QSqlQuery(QSqlDatabase::database());
        q.prepare("SELECT best_fitness, mean_fitness, fitness_stddev, elapsed FROM generations WHERE gen_number = ?");
        q.bindValue(0, generation);
        q.exec();
        if (q.next())
        {
            ui->genNumberLbl->setText(QString::number(generation));
            ui->bestFitnessLbl->setText(q.value(0).toString());
            ui->meanFitnessLbl->setText(q.value(1).toString());
            ui->standardDevLbl->setText(q.value(2).toString());
            ui->timeElapsedLbl->setText(q.value(3).toString());
        }
    }
}

void MainWindow::onSockError(QLocalSocket::LocalSocketError socketError)
{
    qWarning() << "socket error: " << socketError;
}

void MainWindow::onSockConnected()
{
    qInfo() << "socket connected";
}

void MainWindow::onSockReadyRead()
{
    qInfo() << "socket ready read";
    // the socket is use to signal us new data exists in the database, so
    // the data read can be discarded, we now have the information we want.
    m_sock.readAll();
    ui->generationSlider->setMaximum(maxGenerationNumber() / GenerationFrequency);
}

