package main

import (
	"fmt"

	"github.com/cpf2021-gif/gos/utils"
)

func main() {
	fmt.Printf("%+v\n", utils.GlobalConfig)

	utils.LoadConfig("./demo/v0.2/")

	fmt.Printf("%+v\n", utils.GlobalConfig)
}
