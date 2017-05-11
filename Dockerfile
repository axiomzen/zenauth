# Load a base image
FROM ubuntu:14.04

# Add Golang
RUN apt-get update && apt-get install -y golang

# Set ENV variables
# for testing, the env for the process is set by the test
# ENV ENV=TEST
# ENV API_TOKEN="LocalApiToken"
# ENV HASH_SECRET="LocalHashSecret"
# ENV PORT=5000
# ENV $GO_HATCH_DB_ENV_URI=$GO_HATCH_DB_URI
# ENV DOMAIN=localhost:5000

# ports
EXPOSE 5000

# Set working directory
WORKDIR /gopath/src/github.com/axiomzen/zenauth

# Entry point command
CMD ["/gopath/src/github.com/axiomzen/zenauth/authentication"]
