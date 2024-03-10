package repo

import (
	"fmt"
	"log"

	"github.com/go-backend-test/infra"
)

type LineRepo struct {
	LineKey string `env:"LINE_KEY" mapstructure:"LINE_KEY"`
}

func NewLineRepo() ILine {
	lr := LineRepo{}
	initEnv := func() {
		if err := infra.GetEnv(&lr); err != nil {
			log.Fatal(err)
		}
	}

	initEnv()
	fmt.Println(lr)
	return &lr
}
