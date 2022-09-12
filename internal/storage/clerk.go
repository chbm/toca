package storage

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"os"
	"log"

	"github.com/chbm/toca/internal/server"
	. "github.com/chbm/toca/internal/types"
)

const basePath = "/tmp/" // XXX PLACEHOLDER
const peerHost = "localhost"

type Namespaces struct {
	c chan Command
        url string
}

type bagInit int8
const (
	NOLOAD bagInit = iota
	TRYLOAD bagInit = iota
	MUSTLOAD bagInit = iota
)

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


var portCounter int 
func nextPeerAddress() string {
	portCounter += 1
	return fmt.Sprintf("%s:%d", peerHost, portCounter)
}

func runPeers(h chan Command, address string)  {
	peers := httpserver.Start(h)
	http.ListenAndServe(address , peers)	
}


func Start() chan Command {
	portCounter = 60000
	c := make(chan Command)

	bags := map[string]Namespaces{
		"default": bagHolder("default", TRYLOAD),
	}

	go func() {
		for {
			cmd := <-c

			bag, bagexists := bags[cmd.Ns]
			switch cmd.Op {
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
					bags[cmd.Ns] = bagHolder(cmd.Ns, NOLOAD)
					cmd.R <- Result{
						Val: Value{
							V: "",
							Exists: false,
						},
						Err: Success,
					}
				}
			case LoadNs:
				newbag := bagHolder(cmd.Ns, MUSTLOAD)
				if newbag.c == nil {
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
			case GetURL:
				if !bagexists {
				cmd.R <- Result{
					Val: Value{V: "", Exists: false	},
					Err: NoNS,
				}
				} else {
					cmd.R <- Result{
						Val: Value{V: bag.url, Exists: true},
						Err: Success,
					}
				}
			default:
				if !bagexists {
					cmd.R <- Result{
						Val: Value{
							V: "",
							Exists: false,
						}, 
						Err: NoNS,
					}
				} else {
					bag.c <- cmd	
				}
			}
		}
	}()
	return c
}


func bagHolder(name string, policy bagInit) Namespaces {
	var c chan Command
	var bag Bag
	

	if policy == NOLOAD {
		bag = map[string]string{}
	} else {
		b, e := loadFromDisk(name)
		if e != nil {
			if policy == TRYLOAD {
				b = map[string]string{}
			} else {
				return Namespaces{
					c: c,
					url: "",
				}
			}
		}
		bag = b
	}
	
	h := make(chan Command)
	url := nextPeerAddress()
	go runPeers(h, url)	
	logger := log.Default()
	logger.Printf("%s -- http listening on %s", name, url)

	c = make(chan Command)
	go func() {
		for {
			var cmd Command
			select {
			case msg := <-h:
				if msg.Ns != name {
					msg.R <- Result{
						Val: Value{V: "", Exists: false},
						Err: Success,
					}
					continue
				}
				cmd = msg
			case msg := <- c:
				cmd = msg
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
	return Namespaces{
		c: c, 
		url: url,
	}
}
