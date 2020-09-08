# Inspration https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i
FROM golang:1.15.1 as builder

WORKDIR /build
# Let's cache modules retrieval - those don't change so often
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code necessary to build the application
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o emailtocal_cli ./cmd/cli/main.go

# Let's create a /dependancies folder containing just the files necessary for runtime.
# Later, it will be copied as the / (root) of the output image.
WORKDIR /dependancies
RUN cp /build/emailtocal_cli ./emailtocal_cli

# Optional: in case your application uses dynamic linking (often the case with CGO), 
# this will collect dependent libraries so they're later copied to the final image
# NOTE: make sure you honor the license terms of the libraries you copy and distribute
RUN ldd emailtocal_cli | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

# Docker scratch x509 fix
# https://gist.github.com/michaelboke/564bf96f7331f35f1716b59984befc50
RUN apt update && apt --assume-yes upgrade && apt install --assume-yes  ca-certificates
RUN update-ca-certificates

# ------------------------------------------------------------------------------
# Create the minimal runtime image
# ------------------------------------------------------------------------------
FROM scratch

# Docker scratch x509 fix
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy dependancies into root
COPY --chown=0:0 --from=builder /dependancies /

# Set up the app to run as a non-root user inside the /data folder
# User ID 65534 is usually user 'nobody'. 
# The executor of this image should still specify a user during setup.
COPY --chown=65534:0 --from=builder /build/config /config
USER 65534

ENTRYPOINT ["/emailtocal_cli"]
VOLUME /config