FROM alpine

COPY ./app ./

EXPOSE 8000
ENTRYPOINT ["./app"]