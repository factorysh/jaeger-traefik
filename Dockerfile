FROM alpine

EXPOSE 6831
ENV LISTEN 0.0.0.0:6831

COPY bin/jaeger-lite /usr/local/bin/

RUN adduser -D -u 1000 jaeger
USER jaeger

CMD ["jaeger-lite"]