FROM golang:1.11-alpine AS build

WORKDIR /go/src/app
COPY . .

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

## List project dependencies with Gopkg.toml and Gopkg.lock
## These layers are only re-built when Gopkg files are updated
##COPY app/Gopkg.lock Gopkg.toml /go/src/project/
#ADD app /go/src/app
#WORKDIR /go/src/

# Install library dependencies
RUN dep ensure -vendor-only

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
#COPY app /go/src/project/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/places


# This results in a single layer image
FROM golang:1.11-alpine
COPY --from=build /bin/places /bin/places
EXPOSE 8081
ENTRYPOINT ["/bin/places"]
#CMD ["--help"]



