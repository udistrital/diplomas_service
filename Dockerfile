FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY main main
COPY conf/app.conf conf/app.conf
COPY assets/ assets/
COPY runtime/ /

ENTRYPOINT ["/main"]
