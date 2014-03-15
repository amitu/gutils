package gutils

import (
	"fmt"
	"io/ioutil"
)

type intraPacket struct {
	id 	   int64
	packet []byte
}

type FiledChan struct {
	Prod   chan []byte 	// produce writes to this channel
	Cons   chan []byte 	// consumers read from this channel

	intra  chan *intraPacket

	idInFS chan int64 	// latest id writte on disc is available on this channel
	dir    string 		// in which dir to write droppping files
}

func (f *FiledChan) Init(cap int64, dir string) {
	f.Prod 	 = make(chan []byte)
	f.Cons 	 = make(chan []byte, cap)
	f.intra  = make(chan *intraPacket, cap)
	f.idInFS = make(chan int64)
	f.dir 	 = dir

	go f.goProducer()
	go f.goConsumer()
}

func (f *FiledChan) goProducer() {
	var id int64 = -1

	for {
		packet := <- f.Prod
		id += 1  	// we dont care abt id overflow for now, problem?
		ipacket := &intraPacket{
			id: id,
			packet: packet,
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
				f.Cons <- fpacket.packet
				idNext = fpacket.id + 1

			} else {
				// all packets till now should be in intra

				for id := idNext; id < idFromFS; id++ {

					// try to read as much as possible from intra
					select {
					case ipacket = <- f.intra:
						f.Cons <- ipacket.packet
						idNext = ipacket.id + 1
					default:
						break
					}

				}

				// remaining ones must be on disc
				for id := idNext; id <= idFromFS; id++ {
					fpacket2 := f.readPacketFromDisk(id)
					f.Cons <- fpacket2.packet
					idNext = fpacket2.id + 1
				}

				// we have read everything from intra, and we have a packet
				// so lets send it too

				idNext = idFromFS + 1
			}

		} else {

			// got packet.
			// packet is either in sequence, or out of sequence

			if ipacket.id == idNext {

				// packet is in sequence
				f.Cons <- ipacket.packet
				idNext += 1

			} else {

				// packet is out of sequence, meaning till this point
				// everything should be in file

				for id := idNext; id < ipacket.id; id++ {
					fpacket := f.readPacketFromDisk(id)
					f.Cons <- fpacket.packet
				}

				// we have read everything from disk, and we have a packet
				// so lets send it too

				f.Cons <- ipacket.packet
				idNext = ipacket.id + 1

			}
		}
	}
}

func (f *FiledChan) writeToDisk(ipacket intraPacket) {
	filename := fmt.Sprintf("%s/%d.ipacket", f.dir, ipacket.id)

	err := ioutil.WriteFile(filename, ipacket.packet, 0644)
	if err != nil {
		panic(err)
	}

	select {
	case f.idInFS <- ipacket.id:
		return
	default:
		// there already is something in idInFS, lets try to drain it
		select {
		case <- f.idInFS:
		default:
		}

		// finally this one can not block as there is no one else writing to it
		f.idInFS <- ipacket.id
	}
}

func (f *FiledChan) readPacketFromDisk(id int64) intraPacket {
	filename := fmt.Sprintf("%s/%d.ipacket", f.dir, id)

	packet, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return intraPacket{id: id, packet: packet}
}
