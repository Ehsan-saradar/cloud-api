# Build Image
FROM golang:1.16 AS build

WORKDIR /tmp/cloud-api

# Cache Go dependencies like this:
COPY go.mod go.sum ./
RUN go mod download

COPY  . .

# Compile.
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo ./cmd/cloud-api

# Main Image
FROM busybox

RUN mkdir -p openapi/generated
COPY --from=build /tmp/cloud-api/openapi/generated/doc.html ./openapi/generated/doc.html
COPY --from=build /tmp/cloud-api/cloud-api .
COPY config/config.json .

CMD [ "./cloud-api", "config.json" ]
