build: bin vendor
	go build -o bin/jaeger-lite .

bin:
	mkdir -p bin

vendor:
	dep ensure

clean:
	rm -rf bin vendor