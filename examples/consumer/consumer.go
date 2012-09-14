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

package main

import (
  "flag"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "github.com/jedsmith/kafka"
)

var hostname string
var topic string
var partition int
var offset uint64
var maxSize uint
var writePayloadsTo string
var consumerForever bool
var printmessage bool

func init() {
  flag.StringVar(&hostname, "hostname", "localhost:9092", "host:port string for the kafka server")
  flag.StringVar(&topic, "topic", "test", "topic to publish to")
  flag.IntVar(&partition, "partition", 0, "partition to publish to")
  flag.Uint64Var(&offset, "offset", 0, "offset to start consuming from")
  flag.UintVar(&maxSize, "maxsize", 1048576, "offset to start consuming from")
  flag.StringVar(&writePayloadsTo, "writeto", "", "write payloads to this file")
  flag.BoolVar(&consumerForever, "consumeforever", false, "loop forever consuming")
  flag.BoolVar(&printmessage, "printmessage", true, "print the message details to stdout")
}


func main() {
  flag.Parse()
  fmt.Println("Consuming Messages :")
  fmt.Printf("From: %s, topic: %s, partition: %d\n", hostname, topic, partition)
  fmt.Println(" ---------------------- ")
  broker := kafka.NewBrokerConsumer(hostname, topic, partition, offset, uint32(maxSize))

  var payloadFile *os.File = nil
  if len(writePayloadsTo) > 0 {
    var err error
    payloadFile, err = os.Create(writePayloadsTo)
    if err != nil {
      fmt.Println("Error opening file: ", err)
      payloadFile = nil
    }
  }

  consumerCallback := func(msg *kafka.Message) {
    if printmessage {
      msg.Print()
    }
    if payloadFile != nil {
      payloadFile.Write([]byte(fmt.Sprintf("Message at: %d\n", msg.Offset())))
      payloadFile.Write(msg.Payload())
      payloadFile.Write([]byte("\n-------------------------------\n"))
    }
  }

  if consumerForever {
    quit := make(chan bool, 1)
    go func() {
	  sigChan := make(chan os.Signal)
	  signal.Notify(sigChan)
      for {
        sig := <-sigChan
        if sig == syscall.SIGINT {
          quit <- true
        }
      }
    }()

    msgChan := make(chan *kafka.Message)
    go broker.ConsumeOnChannel(msgChan, 10, quit)
    for msg := range msgChan {
      if msg != nil {
        consumerCallback(msg)
      } else {
        break
      }
    }
  } else {
    broker.Consume(consumerCallback)
  }

  if payloadFile != nil {
    payloadFile.Close()
  }

}
