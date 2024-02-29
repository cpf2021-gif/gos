package main

import "github.com/cpf2021-gif/gos/tnet"

func main() {
	s := tnet.NewServer("[gos] Server v0.1")
	s.Serve()
}
