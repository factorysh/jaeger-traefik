Jaeger-lite
===========

POCing [jaeger](https://www.jaegertracing.io/) (and [opentracing](http://opentracing.io/)), without [Cassandra](https://cassandra.apache.org/) and [ElasticSearch](https://www.elastic.co/).

Main target is traefik traces.

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