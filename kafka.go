/*
 * Copyright 2000-2011 NeuStar, Inc. All rights reserved.
 * NeuStar, the Neustar logo and related names and logos are registered
 * trademarks, service marks or tradenames of NeuStar, Inc. All other 
 * product names, company names, marks, logos and symbols may be trademarks
 * of their respective owners.  
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
