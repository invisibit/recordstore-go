# syntax=docker/dockerfile:1.7

FROM golang:1.25 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/recordstore-go .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/recordstore-go /recordstore-go

ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/recordstore-go"]
