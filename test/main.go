package main

import (
	"fmt"
	"io/ioutil"
)

func main() {

	dat, err := ioutil.ReadFile("/Users/wenzhenxi/work/sunmi/KMS/key/rsacert.der")
	check(err)
	fmt.Printf("%s",[]byte(dat))
	fmt.Println("test")
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}