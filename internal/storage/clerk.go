package storage

import (
	"bufio"
	"encoding/gob"
	"os"
)

const basePath = "/tmp/" // XXX PLACEHOLDER

type Bag map[string]string

func dumpToDisk(bag Bag, name string) error {
	file, err := os.CreateTemp(basePath, name)
	if err != nil {
		return err
	}
	temp := file.Name()

	enc := gob.NewEncoder(bufio.NewWriter(file))
	if err = enc.Encode(bag); err != nil {
		return err
	}

	file.Close()
	return os.Rename(temp, basePath+name+".toca")
}

func loadFromDisk(name string) (Bag, error) {
	var bag Bag
	file, err := os.Open(basePath+name+".toca")
	if err != nil {
		return bag, err
	}
	dec := gob.NewDecoder(bufio.NewReader(file))
	err = dec.Decode(&bag)
	if err != nil {
		return bag, err 
	}
	return bag, nil
}

func Start() chan Command {
	c := make(chan Command)

	bags := map[string]map[string]string{
		"default": {},
	}

	go func() {
		for {
			cmd := <-c

			bag, bagexists := bags[cmd.Ns]
			if cmd.Op != CreateNs && !bagexists {
				cmd.R <- Result{
					Val: Value{
						V: "",
						Exists: false,
					}, 
					Err: NoNS,
				}
				continue
			}

			switch cmd.Op {
			case Get:
				v, e := bag[cmd.Key]
				if e {
					cmd.R <- Result{
						Val: Value{
							V: v,
							Exists: true,
						},
						Err: Success,
					}
				} else {
					cmd.R <- Result{
						Val: Value{
							V: "",
							Exists: false,
						},
						Err: Success,
					}
				}
			case Put:
				oldval, e := bag[cmd.Key]
				bag[cmd.Key] = cmd.Value
				cmd.R <- Result{
					Val: Value{
						V: oldval,
						Exists: e,
					}, 
					Err: Success,
				}
			case Delete:
				_, e := bag[cmd.Key]
				delete(bag, cmd.Key)
				cmd.R <- Result{
					Val: Value{
						V: "",
						Exists: e,
					},
					Err: Success,
				}
			case CreateNs:
				if bagexists {
					cmd.R <- Result{
						Val: Value{
							V: "",
							Exists: true,
						},
						Err: Conflict,
					}
				} else {
					bags[cmd.Ns] = make(map[string]string)
					cmd.R <- Result{
						Val: Value{
							V: "",
							Exists: false,
						},
						Err: Success,
					}
				}
			}
		}
	}()

	return c
}
