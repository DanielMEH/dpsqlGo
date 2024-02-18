FROM golang

WORKDIR /go/src/github.com/DanielMEH/database

# Copy only necessary files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . ./

# Install air
RUN go install github.com/cosmtrek/air@latest && \
    air -v

# Build the application
RUN go install .

# Initialize air (if necessary)
RUN air init -d -c .air.toml
RUN mkdir tmp && chmod 777 tmp
# Set the default command to run the application
CMD ["/go/bin/air","-c", ".air.toml"]

#CMD ["/go/bin/database"]