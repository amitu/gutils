package gutils

/* vim: set tabstop=4 */

import (
	"fmt"
	"io/ioutil"
	"os"
	"errors"
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

func (f *FiledChan) Init(cap int64) error {
	f.Prod 	 =  make(chan interface{})
	f.Cons 	 =  make(chan interface{})
	f.intra  =  make(chan *intraPacket, cap)
	f.idInFS =  make(chan int64, 1)
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

			if idFromFS == idNext {

				fpacket := f.readPacketFromDisk(idFromFS)
				f.Cons <- fpacket.Packet
				idNext = fpacket.ID + 1

			} else {
				// all packets till now should be in intra
				for id := idNext; id < idFromFS; id++ {

					// try to read as much as possible from intra
					select {
					case ipacket = <- f.intra:
						if ipacket.ID == idNext {
							f.Cons <- ipacket.Packet
							idNext = ipacket.ID + 1
							ipacket = nil
						} else {
							break
						}
					default:
						break
					}

				}

				// remaining ones must be on disc
				for id := idNext; id <= idFromFS; id++ {
					fpacket2 := f.readPacketFromDisk(id)
					f.Cons <- fpacket2.Packet
					idNext = fpacket2.ID + 1
				}

				// we have read everything from intra, and we have a packet
				// so lets send it too

				if ipacket != nil {
					f.Cons <- ipacket.Packet
					idNext = ipacket.ID + 1
					ipacket = nil
				}
			}

		} else {
			// got packet.
			// packet is either in sequence, or out of sequence

			if ipacket.ID == idNext {

				// packet is in sequence
				f.Cons <- ipacket.Packet
				idNext += 1

			} else {

				// packet is out of sequence, meaning till this point
				// everything should be in file

				for id := idNext; id < ipacket.ID; id++ {
					fpacket := f.readPacketFromDisk(id)
					f.Cons <- fpacket.Packet
				}

				// we have read everything from disk, and we have a packet
				// so lets send it too

				f.Cons <- ipacket.Packet
				idNext = ipacket.ID + 1

			}
		}
	}
}

func (f *FiledChan) writeToDisk(ipacket intraPacket) {
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

	select {
	case f.idInFS <- ipacket.ID:
		return
	default:
		// there already is something in idInFS, lets try to drain it
		select {
		case <- f.idInFS:
		default:
		}

		// finally this one can not block as there is no one else writing to it
		f.idInFS <- ipacket.ID
	}
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
