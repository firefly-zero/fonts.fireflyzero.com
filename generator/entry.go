package generator

import (
	"fmt"
	"os"
)

func Run() error {
	files, err := os.ReadDir("bdf")
	if err != nil {
		return fmt.Errorf("open bdf dir: %v", err)
	}
	for _, file := range files {
		err1 := handleBDFFile("bdf/" + file.Name())
		if err1 != nil {
			return err1
		}
	}
	return nil
}
