package app

import "time"

type Config struct {
	App   Application `yaml:"app"`
	Store `yaml:"store"`
}

type Application struct {
	UpdateProcessorWorkerCount int `yaml:"update_processor_worker_count"`
}

type Store struct {
	TTL         time.Duration `yaml:"ttl"`
	MaxCapacity int64         `yaml:"max_capacity"`
}
