package storage

func Start() chan Command {
	c := make(chan Command)

	bags := map[string]map[string]string{
		"default": map[string]string{},
	}

	go func() {
		for {
			cmd := <-c

			bag, bagexists := bags[cmd.Ns]
			if cmd.Op != CreateNs && !bagexists {
				cmd.R <- Value{
					V: "",
					Exists: false,
				}
				continue
			}

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
			case CreateNs:
				if bagexists {
					cmd.R <- Value{
						V: "",
						Exists: true,
					}
				} else {
					bags[cmd.Ns] = make(map[string]string)
					cmd.R <- Value{
						V: "",
						Exists: false,
					}
				}
			}
		}
	}()

	return c
}
