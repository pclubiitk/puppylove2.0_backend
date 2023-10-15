FROM golang:1.21 AS build-stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /pupplovebackend

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /
COPY --from=build-stage /pupplovebackend /pupplovebackend
COPY ./.env .env

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT [ "/pupplovebackend" ]

