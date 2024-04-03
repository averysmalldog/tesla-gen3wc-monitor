clean:
	@rm -rf bin/

# builds for standard docker installs
build-docker:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/docker/tesla-gen3wc-monitor .

# builds for your local system
build: 
	@go build -o bin/local/tesla-gen3wc-monitor .