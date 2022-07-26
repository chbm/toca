package storage


func Start() chan Command {
	c := make(chan Command)

	bag := make(map[string]string)

	go func() {
		for {
			cmd := <-c

			switch cmd.Op {
			case Get:
				v, e := bag[cmd.Key]
				if e {
					cmd.R <- Value{
						V: v,
						Exists: true,
					}
				} else {
					cmd.R <- Value{
						V: "",
						Exists: false,
					}
				}
			case Put:
				bag[cmd.Key] = cmd.Value
				cmd.R <- Value{
					V: "",
					Exists: true,
				}
			case Delete:
				_, e := bag[cmd.Key]
				delete(bag, cmd.Key)
				cmd.R <- Value{
					V: "",
					Exists: e,
				}
			}
		}
	}()

	return c
}
