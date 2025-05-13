package config

import "sync"

type Environment string

const (
	EnvironmentLocal Environment = "local"
	EnvironmentProd  Environment = "prod"
)

var (
	EnvironmentValue Environment
	once             sync.Once
)

func loadEnvironment(value string) {
	once.Do(func() {
		switch value {
		case string(EnvironmentLocal):
			EnvironmentValue = EnvironmentLocal
		case string(EnvironmentProd):
			EnvironmentValue = EnvironmentProd
		default:
			EnvironmentValue = EnvironmentLocal
		}
	})
}
