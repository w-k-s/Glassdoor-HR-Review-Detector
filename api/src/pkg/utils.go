package pkg

import "log"

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
