FROM golang:1.18-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh build-base

# Set the Current Working Directory inside the container
WORKDIR /app

COPY src/go.mod ./src/
COPY src/go.sum ./src/
RUN cd src && go mod download

COPY ./src ./src
COPY ./config/config.yml ./config/config.yml

# Build the project
RUN cd src && go build -o websocket -tags musl

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
WORKDIR /app/src

CMD ["./websocket", "serve"]
