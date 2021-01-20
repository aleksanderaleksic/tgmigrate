package config

type History struct {
	Storage HistoryStorage
}

type HistoryStorage struct {
	Type   string
	Config interface{}
}
