FROM golang:1.24-alpine AS build-stage

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/overseerr-auto-decline

FROM gcr.io/distroless/base-debian12 AS release-page
WORKDIR /
COPY --from=build-stage /app/overseerr-auto-decline /app/overseerr-auto-decline
USER nonroot:nonroot

CMD ["/app/overseerr-auto-decline"]