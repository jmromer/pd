/*
Copyright © 2020 Jake Romer <mail@jakeromer.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// Exists returns true if the given file or directory exists, else false.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// TODO: Add doc
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", HomeDir(), 1)
	}
	return path
}

// TODO: Add doc
func HomeDir() string {
	home, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home
}

// TODO: Add doc
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
