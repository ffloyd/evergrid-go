package worker

import "github.com/Sirupsen/logrus"

// Stats - это структура, которая содержит статистику об использовании ресурсов воркера
type Stats struct {
	UploadingTicks   int
	BuildingTicks    int
	CalculatingTicks int
	MoneySpent       float64
}

// StatsReport выводит суммарную статистику использования воркеров
func StatsReport(workers []*Worker, log *logrus.Entry) {
	totalUploadingTicks := 0
	totalBuildingTicks := 0
	totalCalculatingTicks := 0
	totalMoneySpent := 0.0

	statsCollection := make([]Stats, len(workers))
	for i, worker := range workers {
		statsCollection[i] = worker.Stats()
	}

	for _, stats := range statsCollection {
		totalUploadingTicks += stats.UploadingTicks
		totalBuildingTicks += stats.BuildingTicks
		totalCalculatingTicks += stats.CalculatingTicks
		totalMoneySpent += stats.MoneySpent
	}

	log.WithField("value", totalUploadingTicks).Info("Total uploading ticks")
	log.WithField("value", totalBuildingTicks).Info("Total building ticks")
	log.WithField("value", totalCalculatingTicks).Info("Total calculating ticks")
	log.WithField("value", totalMoneySpent).Info("Total money spent")
}
