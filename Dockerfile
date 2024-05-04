FROM golang:1.20

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN go build -o /backplate

RUN mkdir images inbox

# Copy frontend
COPY frontend/dist ./frontend/dist

EXPOSE 8090

# Run
CMD ["/backplate"]
