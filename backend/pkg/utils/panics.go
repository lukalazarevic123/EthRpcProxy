package utils

import "log"

func Handle(err error) {
	if err != nil {
		log.Print(err.Error())
		panic(err)
	}
}
