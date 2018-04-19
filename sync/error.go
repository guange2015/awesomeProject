package sync

import "log"

func ce(err error, msg string) {
	if err != nil {
		log.Fatalf("%s error: %v", msg, err)
	}
}
