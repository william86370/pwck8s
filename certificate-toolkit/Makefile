.PHONY: all clean

all: ca-cert client-cert server-cert

ca-cert:
	./gen-ca.sh

client-cert:
	./gen-client.sh

server-cert:
	./gen-server.sh

clean:
	rm -rf ca
	rm -rf client
	rm -rf server