Jaeger-lite
===========

POCing [jaeger](https://www.jaegertracing.io/) (and [opentracing](http://opentracing.io/)), without [Cassandra](https://cassandra.apache.org/) and [ElasticSearch](https://www.elastic.co/).

Main target are traefik traces.

TL;DR
-----

Jaeger is one implementation of tracing.
Tracing can be seen as light logging, with typed tags, timestamp and parent reference.
Tracing is non blocking, and can explain what happens in a request compound by services (parallelizeds or sequentials).

Demo time
---------

    +--------+   +---------+   +-----+
    | client +-->| traefik +-->| web |
    +--------+   +----+----+   +-----+
                      |
                      v
                 +--------+
                 | jaeger |
                 +--------+

Do it

    cd demo

Launch backround services

    docker-compose up -d traefik

Watch *jaeger-lite* logs

    docker-compose logs jaeger

In another window, trigger some curl action

    docker-compose up client

Licence
-------

3 terms BSD licence, ©2018 Mathieu Lecarme