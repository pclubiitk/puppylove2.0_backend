FROM ubuntu:20.04 AS base-stage

# Install CA certificates
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

FROM golang:1.21 AS build-stage

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /pupplovebackend

FROM base-stage AS final-stage
COPY --from=build-stage /pupplovebackend /pupplovebackend
COPY --from=build-stage /app/.env /.env
EXPOSE 8080
ENV PORT 8080
# set hostname to localhost
ENV HOSTNAME "0.0.0.0"
ENTRYPOINT ["./pupplovebackend"]
