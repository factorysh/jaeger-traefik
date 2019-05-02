# jaeger-traefik

POCing [jaeger](https://www.jaegertracing.io/) (and [opentracing](http://opentracing.io/)), without [Cassandra](https://cassandra.apache.org/) and [ElasticSearch](https://www.elastic.co/).

Main target are traefik traces.

## TL;DR

Jaeger is an implementation of tracing.
Tracing can be seen as light logging, with typed tags, timestamp and parent reference, for composed services.
Tracing is non blocking, and can explain what happens in a request compound by services (parallelizeds or sequentials).

## Demo time

    +-------------+   +---------+   +--------------+
    | http client +-->| traefik +-->| web services |
    +-------------+   +----+----+   +--------------+
                           |
                           v
                       +--------+
                       | jaeger |
                       +--------+

Do it

    cd demo

Launch backround services

    docker-compose up -d traefik

Watch _jaeger-traefik_ logs

    docker-compose logs jaeger

In another window, trigger some curl action

    docker-compose up client

Wait for consolidation, and ask prometheus probe

    docker-compose up promclient | grep apdex

You can do all that steps with one command :

    make demo

## Licence

3 terms BSD licence, ©2018 Mathieu Lecarme
