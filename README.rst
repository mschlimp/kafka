=========================
github.com/jedsmith/kafka
=========================

Go bindings for Kafka
=====================

This is **kafka.go**, a repository containing Go_ bindings for the
`Apache Kafka`_ pub/sub messaging system developed by LinkedIn. The bindings
were originally authored by `@jdamick`_ but have largely been abandoned; I
brought them up to date for Go 1 for my needs at work. Since that implied
reorganizing the repository, I have elected to maintain a separate repository
instead of issuing a pull request and disrupting that codebase.

.. _Go: http://golang.org/
.. _`Apache Kafka`: http://incubator.apache.org/kafka/
.. _`@jdamick`: https://github.com/jdamick/kafka.go


Getting Started
---------------

The import path for kafka.go is **github.com/jedsmith/kafka**.

.. code-block:: text

    $ export GOPATH=`pwd`
    $ go get github.com/jedsmith/kafka
    $ go test github.com/jedsmith/kafka
    ok     github.com/jedsmith/kafka      0.037s

To play around with publishers and consumers, there are also example binaries.
Assuming a running Kafka server on your local machine:

.. code-block:: text

    $ go get github.com/jedsmith/kafka/examples/consumer
    $ go get github.com/jedsmith/kafka/examples/producer

    $ bin/producer -topic test -message "sup"
    $ bin/consumer -topic test -consumeforever

    2012/09/14 15:28:00 ----- Begin Message ------
    2012/09/14 15:28:00 magic: 0
    2012/09/14 15:28:00 checksum: ABBBF394
    2012/09/14 15:28:00 payload: 737570
    2012/09/14 15:28:00 payload(string): sup
    2012/09/14 15:28:00 offset: 0
    2012/09/14 15:28:00 ----- End Message ------

API Usage
---------

**A complete refactor is coming for the API. Don't get too attached!**

Full documentation is available with godoc_, but some quick examples follow.

.. _godoc: http://go.pkgdoc.org/github.com/jedsmith/kafka

Publishing
~~~~~~~~~~

.. code-block:: go

    broker := kafka.NewBrokerPublisher("localhost:9092", "mytesttopic", 0)
    broker.Publish(kafka.NewMessage([]byte("tesing 1 2 3")))

Consumer
~~~~~~~~

.. code-block:: go

    broker := kafka.NewBrokerConsumer("localhost:9092", "mytesttopic", 0, 0, 1048576)
    broker.Consume(func(msg *kafka.Message) { msg.Print() })

Or, the consumer can use a channel-based approach:

.. code-block:: go

    broker := kafka.NewBrokerConsumer("localhost:9092", "mytesttopic", 0, 0, 1048576)
    go broker.ConsumeOnChannel(msgChan, 10, quitChan)

Consuming Offsets
~~~~~~~~~~~~~~~~~

.. code-block:: go

    broker := kafka.NewBrokerOffsetConsumer("localhost:9092", "mytesttopic", 0)
    offsets, err := broker.GetOffsets(-1, 1)

TODO
----

- Documentation cleanup to godoc standards.
- Code cleanup and optimization.
- Optional Zookeeper integration.

License
-------

Copyright |copy| 2000-2011 NeuStar, Inc.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

- Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

- Redistributions in binary form must reproduce the above copyright notice, this
  list of conditions and the following disclaimer in the documentation and/or
  other materials provided with the distribution.

- Neither the name of NeuStar, Inc., Jed Smith, nor the names of their
  contributors may be used to endorse or promote products derived from this
  software without specific prior written permission.

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

.. |copy| unicode:: U+000A9 .. COPYRIGHT SIGN

