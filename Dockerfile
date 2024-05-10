FROM golang:1.22
WORKDIR /usr/src/app
COPY . .
RUN sh install-amd64.sh
RUN echo $(ls -a)
RUN sh build.sh
EXPOSE 8080
CMD ["tmp/main"]
