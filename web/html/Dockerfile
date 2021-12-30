FROM alpine as build
WORKDIR /download
RUN wget https://github.com/gohugoio/hugo/releases/download/v0.91.2/hugo_0.91.2_Linux-64bit.tar.gz -O hugo.tar.gz && \
    tar xvf hugo.tar.gz

FROM scratch
COPY --from=build /download/hugo /bin/hugo
ENTRYPOINT ['/bin/hugo']
