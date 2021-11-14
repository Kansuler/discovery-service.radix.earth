# Build server stage
FROM golang:1.17.3 as Builder
WORKDIR /build
ENV CGO_ENABLED 0
ENV GOOS linux
COPY . .
RUN go build -a -installsuffix cgo -o binary ./cmd/service

# Production stage
FROM scratch
COPY --from=Builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=Builder /usr/share/zoneinfo /usr/share/zoneinfo/
COPY --from=Builder /build/binary ./
ENTRYPOINT ["./binary"]