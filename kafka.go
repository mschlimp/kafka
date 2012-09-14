/*

github.com/jedsmith/kafka: Go bindings for Kafka

Copyright 2000-2011 NeuStar, Inc. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of NeuStar, Inc., Jed Smith, nor the names of
	  contributors may be used to endorse or promote products derived from
	  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL NEUSTAR OR JED SMITH BE LIABLE FOR ANY DIRECT,
INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

NeuStar, the Neustar logo and related names and logos are registered
trademarks, service marks or tradenames of NeuStar, Inc. All other 
product names, company names, marks, logos and symbols may be trademarks
of their respective owners.  

*/

package kafka

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	MAGIC_DEFAULT = 0
	NETWORK       = "tcp"
)

type Broker struct {
	topic     string
	partition int
	hostname  string
}

func newBroker(hostname string, topic string, partition int) *Broker {
	return &Broker{topic: topic,
		partition: partition,
		hostname:  hostname}
}

func (b *Broker) connect() (conn *net.TCPConn, err error) {
	raddr, err := net.ResolveTCPAddr(NETWORK, b.hostname)
	if err != nil {
		log.Println("Fatal Error: ", err)
		return nil, err
	}
	conn, err = net.DialTCP(NETWORK, nil, raddr)
	if err != nil {
		log.Println("Fatal Error: ", err)
		return nil, err
	}
	return conn, err
}

// returns length of response & payload & err
func (b *Broker) readResponse(conn *net.TCPConn) (uint32, []byte, error) {
	reader := bufio.NewReader(conn)
	length := make([]byte, 4)
	lenRead, err := io.ReadFull(reader, length)
	if err != nil {
		return 0, []byte{}, err
	}
	if lenRead != 4 || lenRead < 0 {
		return 0, []byte{}, errors.New("invalid length of the packet length field")
	}

	expectedLength := binary.BigEndian.Uint32(length)
	messages := make([]byte, expectedLength)
	lenRead, err = io.ReadFull(reader, messages)
	if err != nil {
		return 0, []byte{}, err
	}

	if uint32(lenRead) != expectedLength {
		return 0, []byte{}, errors.New(fmt.Sprintf("Fatal Error: Unexpected Length: %d  expected:  %d", lenRead, expectedLength))
	}

	errorCode := binary.BigEndian.Uint16(messages[0:2])
	if errorCode != 0 {
		return 0, []byte{}, errors.New(fmt.Sprintf("error %d", errorCode))
	}
	return expectedLength, messages[2:], nil
}
