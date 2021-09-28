package main

import (
	"os"
	"uff/scrap"
)

func main() {
	scrap.FetchFriends(os.Getenv("url"))
}
