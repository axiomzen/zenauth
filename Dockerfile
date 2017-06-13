# Load a base image
FROM scratch

# Entry point command
CMD ["/zenauth"]

# Ports
EXPOSE 5000
EXPOSE 5001

# Get SSL root certificates
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt

ADD https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem /etc/ssl/certs/

# Always add the binary last to maximize caching
ADD zenauth /zenauth
