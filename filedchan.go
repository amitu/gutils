package gutils

/* vim: set tabstop=4 */

import (
	"os"
	"fmt"
	"time"
	"errors"
	"io/ioutil"
	"encoding/gob"
)

type intraPacket struct {
	ID 	   int64
	Packet interface{}
}

type FiledChan struct {
	Prod   chan interface{} 	// produce writes to this channel
	Cons   chan interface{} 	// consumers read from this channel

	intra  chan *intraPacket

	idInFS chan int64 	// latest id writte on disc is available on this channel
	Dir    string 		// in which dir to write droppping files
}

func (f *FiledChan) Init(cap, dcap int64) error {
	f.Prod 	 =  make(chan interface{})
	f.Cons 	 =  make(chan interface{})
	f.intra  =  make(chan *intraPacket, cap)
	f.idInFS =  make(chan int64, dcap + 1)
	dir, err := ioutil.TempDir("", "filedchan")
	if err != nil {
		return err
	}
	f.Dir = dir

	err = f.checkDir()
	if err != nil {
		return err
	}

	go f.goProducer()
	go f.goConsumer()

	return nil
}

func (f *FiledChan) Quit() error {
	// TODO: convert it to Close()
	// implement proper close semantics
	return os.RemoveAll(f.Dir)
}

func (f *FiledChan) checkDir() error {
	// todo, dir must be writable and empty

	filename := fmt.Sprintf("%s/test.ipacket", f.Dir)
	err := ioutil.WriteFile(filename, []byte("test"), 0644)

	if err != nil {
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(f.Dir)
	if err != nil {
		return err
	}

	if len(files) != 0 {
		errors.New("Dir is not empty")
	}

	return nil
}

func (f *FiledChan) goProducer() {
	var id int64 = -1
	for {
		id += 1  	// we dont care abt id overflow for now, problem?
		packet := <- f.Prod
		ipacket := &intraPacket{
			ID: id,
			Packet: packet,
		}

		select {
		case f.intra <- ipacket:
			continue
		default:
		}

		// looks like intra is full! write to disc
		f.writeToDisk(*ipacket)
	}
}

func (f *FiledChan) goConsumer() {
	idNext := int64(0)
	for {
		var ipacket *intraPacket
		var idFromFS int64

		select {
		case ipacket = <- f.intra:
		case idFromFS = <- f.idInFS:
		}

		// either ipacket is not nil or idFromFS is set

		if ipacket == nil {

			// all packets till now should be in intra

			for id := idNext; id < idFromFS; id++ {

				// try to read as much as possible from intra
				select {
				case ipacket = <- f.intra:
					if idNext != ipacket.ID {
						panic("logic error")
					}
					f.Cons <- ipacket.Packet
					idNext = ipacket.ID + 1
				default:
					break
				}

			}

			// we have read everything from intra, and we have a packet
			// so lets send it too

			if idNext == idFromFS {
				fpacket := f.readPacketFromDisk(idFromFS)
				f.Cons <- fpacket.Packet
				idNext += 1
			}

		} else {
			// got packet.
			// packet is either in sequence, or out of sequence

			// packet is out of sequence, meaning till this point
			// everything should be in file

			for id := idNext; id < ipacket.ID; id++ {
				select {
				case i := <- f.idInFS:
					if i != idNext {
						panic("logic error")
					}
				case <- time.After(1e8):
					panic("logic error default")
				}
				fpacket := f.readPacketFromDisk(id)
				f.Cons <- fpacket.Packet
				idNext += 1
			}

			// we have read everything from disk, and we have a packet
			// so lets send it too

			if idNext != ipacket.ID {
				panic("logic error")
			}

			f.Cons <- ipacket.Packet
			idNext = ipacket.ID + 1

		}
	}
}

func (f *FiledChan) writeToDisk(ipacket intraPacket) {
	gob.Register(S3Upload{})
	filename := fmt.Sprintf("%s/%d.ipacket", f.Dir, ipacket.ID)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	enc := gob.NewEncoder(file)

	err = enc.Encode(ipacket)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	f.idInFS <- ipacket.ID
}

func (f *FiledChan) readPacketFromDisk(id int64) intraPacket {
	filename := fmt.Sprintf("%s/%d.ipacket", f.Dir, id)

	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	dec := gob.NewDecoder(file)

	var ipacket intraPacket
	err = dec.Decode(&ipacket)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}

	return ipacket
}
