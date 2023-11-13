run:
	@PORT=8443 \
	TLS_CERT=/home/ubuntu/workspace/pwck8s/certificate-toolkit/server/server.crt \
	TLS_KEY=/home/ubuntu/workspace/pwck8s/certificate-toolkit/server/server.key \
	CA_CERT=/home/ubuntu/workspace/pwck8s/certificate-toolkit/ca/ca.crt \
	PROXY_URL=localhost:8080 \
	DEBUG=true \
	go run main.go




run-headers:
	@go run headers/headers.go