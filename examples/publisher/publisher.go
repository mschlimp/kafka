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
  "github.com/jedsmith/kafka"
)

var hostname string
var topic string
var partition int
var message string
var messageFile string

func init() {
  flag.StringVar(&hostname, "hostname", "localhost:9092", "host:port string for the kafka server")
  flag.StringVar(&topic, "topic", "test", "topic to publish to")
  flag.IntVar(&partition, "partition", 0, "partition to publish to")
  flag.StringVar(&message, "message", "", "message to publish")
  flag.StringVar(&messageFile, "messagefile", "", "read message from this file")
}

func main() {
  flag.Parse()
  fmt.Println("Publishing :", message)
  fmt.Printf("To: %s, topic: %s, partition: %d\n", hostname, topic, partition)
  fmt.Println(" ---------------------- ")
  broker := kafka.NewBrokerPublisher(hostname, topic, partition)

  if len(message) == 0 && len(messageFile) != 0 {
    file, err := os.Open(messageFile)
    if err != nil {
      fmt.Println("Error: ", err)
      return
    }
    stat, err := file.Stat()
    if err != nil {
      fmt.Println("Error: ", err)
      return
    }
    payload := make([]byte, stat.Size())
    file.Read(payload)
    timing := kafka.StartTiming("Sending")
    broker.Publish(kafka.NewMessage(payload))
    timing.Print()
    file.Close()
  } else {
    timing := kafka.StartTiming("Sending")
    broker.Publish(kafka.NewMessage([]byte(message)))
    timing.Print()
  }
}
