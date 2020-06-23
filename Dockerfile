FROM golang AS build
WORKDIR /go/src/app
ADD *.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app .

FROM scratch
COPY --from=build /go/bin/app /
COPY *.html *.gif /
EXPOSE 8080
ENTRYPOINT ["/app"]
