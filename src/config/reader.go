package config

import "io/ioutil"

type ConfigReader interface {
	ReadFile(path string) ([]byte, error)
}

type FileReader struct{}

/**
* Reads a file and returns an bytes array
 */
func (*FileReader) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
