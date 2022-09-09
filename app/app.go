package app

import (
	"errors"
	"fmt"
)

func Run(config string, c string, debug string) error {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	fmt.Println(config, c, debug)
	err := errors.New("text string")
	if err != nil {
		return err
	}
	return nil

}
