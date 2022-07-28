package storage

import (
	"bufio"
	"encoding/gob"
	"os"
	"io"
)

const basePath = "/tmp/" // XXX PLACEHOLDER

type Bag map[string]string

func dumpToDisk(bag Bag, name string) error {
	file, err := os.CreateTemp(basePath, name)
	if err != nil {
		return err
	}
	temp := file.Name()
	
	b := bufio.NewWriter(file)
	enc := gob.NewEncoder(b)
	if err = enc.Encode(bag); err != nil {
		return err
	}
	b.Flush()	
	file.Close()
	return os.Rename(temp, basePath+name+".toca")
}

func loadFromDisk(name string) (Bag, error) {
	bag := Bag{}
	file, err := os.Open(basePath+name+".toca")
	if err != nil {
		return bag, err
	}
	dec := gob.NewDecoder(bufio.NewReader(file))
	err = dec.Decode(&bag)
	if err != nil && err == io.EOF {
		err = nil
	}
	return bag, err
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
			if cmd.Op != CreateNs && cmd.Op != LoadNs  && !bagexists {
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
			case LoadNs:
				newbag, err := loadFromDisk(cmd.Ns)
				if err != nil {
					cmd.R <- Result{
						Val: Value{},
						Err: Failure,
					}
				} else {
					bags[cmd.Ns] = newbag
					cmd.R <- Result{
							Val: Value{},
							Err: Success,
						}
				}
			case SaveNs:
				r := Result{
					Val: Value{},
					Err: Success,
				}
				if dumpToDisk(bag, cmd.Ns) != nil {
					r.Err = Failure
				}
				cmd.R <- r
			}
		}
	}()

	return c
}
