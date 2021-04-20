# build stage
FROM golang:1.15-alpine AS build-env
RUN apk add --update make
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app
RUN CGO_ENABLED=0 GOOS=linux make

# final stage
FROM alpine:3.9
LABEL maintainer="m.vorobev"
WORKDIR /app
COPY --from=build-env /go/src/app/bin/app /app/
COPY --from=build-env /go/src/app/assets/ /app/assets/
COPY --from=build-env /go/src/app/template/ /app/template/
RUN chmod +x /app/app
CMD ["/app/app"]
