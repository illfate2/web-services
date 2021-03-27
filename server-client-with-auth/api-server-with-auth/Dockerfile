FROM golang:alpine3.12 as build
RUN apk add build-base
WORKDIR server
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -o /server/bin/server ./cmd/

FROM scratch
COPY --from=build /server/bin/server .
CMD ["./server"]