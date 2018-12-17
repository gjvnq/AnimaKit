package AnimaKit

import (
	"os"

	"github.com/robertkrimen/otto"
)

func LoadScriptFromFile(path string) (*otto.Otto, error) {
	vm := otto.New()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	_, err = vm.Run(file)
	return vm, err
}
