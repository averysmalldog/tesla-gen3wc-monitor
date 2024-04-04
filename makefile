clean:
	@rm -rf bin/

# Builds for standard docker installs
build-docker:
	@docker build -t averysmalldog/tesla-gen3wc-monitor:latest .

# Builds the image locally and then spins up the compose suite without needing
# to pull from Docker Hub. I originally wanted to call this `docker-compose-up`
# but that was just too much to remember/type.
up:
	@make build-docker
	@docker compose up -d

# Builds for your local system
build: 
	@go build -o bin/local/tesla-gen3wc-monitor .