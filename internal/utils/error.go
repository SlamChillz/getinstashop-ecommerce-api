package utils

import "fmt"

func ConcatenateErrors(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	return fmt.Errorf("%v; %v", err1.Error(), err2.Error())
}
