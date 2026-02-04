package util

import "log"

func Assert(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			log.Fatal(msg[0], ": ", err)
		} else {
			log.Fatal(err)
		}
	}
}
