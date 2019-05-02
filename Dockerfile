FROM alpine

EXPOSE 6831
ENV LISTEN 0.0.0.0:6831

COPY bin/jaeger-traefik /usr/local/bin/

RUN adduser -D -u 1000 jaeger
USER jaeger

CMD ["jaeger-traefik", "serve"]