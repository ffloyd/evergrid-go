package worker

// Stats - это структура, которая содержит статистику об использовании ресурсов воркера
type Stats struct {
	UploadingTicks   int
	BuildingTicks    int
	CalculatingTicks int
}
