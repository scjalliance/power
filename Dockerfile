# --------
# Stage 1: Build
# -------
FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /go/src/github.com/scjalliance/power
COPY . .

WORKDIR /go/src/github.com/scjalliance/power/cmd/power

ENV CGO_ENABLED=0

RUN go-wrapper download
RUN go-wrapper install

# --------
# Stage 2: Release
# --------
FROM gcr.io/distroless/base

COPY --from=builder /go/bin/power /

CMD ["/power"]
