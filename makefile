clean:
	@rm -rf bin/

# builds for standard docker installs
build-docker:
	@docker build .

# builds for your local system
build: 
	@go build -o bin/local/tesla-gen3wc-monitor .