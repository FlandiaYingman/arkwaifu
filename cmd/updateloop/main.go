package main

import "arkwaifu/internal/app/updateloop"

func main() {
	err := updateloop.UpdateResources()
	if err != nil {
		panic(err)
	}
}
