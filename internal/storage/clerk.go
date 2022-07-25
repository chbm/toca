package clerk

import "commands"

func Start() chan {
	c := make(chan Command)

	bag := make(map[string]string)

	go func() {
		for {
		cmd := <-c

		switch cmd.op {
			case commands.Get:
				v, e := bag[cmd.key]
				if e {
					cmd.r <- v
				} else {
					cmd.r <- nil
				}
			case comamnds.Put:
				bag[cmd.key] = cmd.value
		}
	}
	}

	return c
}
