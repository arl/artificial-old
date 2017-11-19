#include <QDebug>
#include <QFileDialog>
#include <QSqlDatabase>
#include <QFileInfo>
#include <QPixmap>
#include <QSqlQuery>

#include "mainwindow.h"
#include "ui_mainwindow.h"

constexpr int GenerationFrequency = 100;

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

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    m_dbOpened = false;
}

MainWindow::~MainWindow()
{
    delete ui;
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
    auto db = QSqlDatabase::addDatabase("QSQLITE");
    db.setDatabaseName(fileName);
    m_dbOpened = db.open();
    if (m_dbOpened)
    {
        m_dir = QFileInfo(fileName).absolutePath();

        // configure slider
        ui->generationSlider->setMinimum(0);
        ui->generationSlider->setMaximum(maxGenerationNumber() / GenerationFrequency);
        ui->generationSlider->setTickInterval(1);
        ui->generationSlider->setEnabled(true);
        ui->generationSlider->setSingleStep(1);
        ui->generationSlider->setValue(0);

        // show images
        showGenerationImage(0);
        showGenerationData(0);
        ui->refImage->setPixmap(loadPixmap("_ref.png"));
    }
    qInfo() << "db.open() -> " << m_dbOpened;
}

void MainWindow::showGenerationImage(int value)
{
    auto generation = value * GenerationFrequency;
    if (m_dbOpened)
        ui->curImage->setPixmap(loadPixmap(QString::number(generation) + ".png"));
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
