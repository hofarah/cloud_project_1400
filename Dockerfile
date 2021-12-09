FROM golang:1.17


ARG SERVICE=default
ENV SERVICEN=$SERVICE

RUN git clone --depth 1 https://github.com/hofarah/cloud_project_1400.git /go/src/cloud_1400 --branch master


ENV GO111MODULE=on

WORKDIR /go/src/cloud_1400

RUN go mod tidy
RUN go mod download


RUN cd services/$SERVICE && go build -o main .

EXPOSE 7575
CMD /go/src/cloud_1400/services/${SERVICEN##*/}/main
