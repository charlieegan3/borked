FROM golang:1.10 as build

WORKDIR /go/src/github.com/charlieegan3/borked

RUN go get -u github.com/gobuffalo/packr/packr

COPY . .

RUN CGO_ENABLED=0 packr build -o borked


FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/charlieegan3/borked/borked /

CMD ["/borked"]
