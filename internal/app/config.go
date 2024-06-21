package app

type Config struct {
	App Application `yaml:"app"`
}

type Application struct {
	UpdateProcessorWorkerCount int `yaml:"update_processor_worker_count"`
}
