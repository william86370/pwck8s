.PHONY: all clean

all: ca-cert client-cert

ca-cert:
	./gen-ca.sh

client-cert:
	./gen-client.sh

clean:
	rm -rf ca
	rm -rf client
