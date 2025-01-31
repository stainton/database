/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/stainton/database/cmd"

func main() {
	svr := cmd.NewDbServer()
	svr.Execute()
}
