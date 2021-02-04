package test

//go:generate mockgen -destination ./mock/history_mock.go -source ../history/history.go -package mock
//go:generate mockgen -destination ./mock/state_mock.go -source ../state/state.go -package mock
