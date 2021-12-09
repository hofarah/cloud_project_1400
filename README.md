# cloud_project_1400
## docker microservices

here we have three web service that can use in docker or Separately :
- Auth
- DataAnalytic
- Global

all these services developed by golang
## Auth service
authentication of users

## Global service
fetch some data about games

## Data Analytic service
draw charts that you want

# Running
we offer two method for up and running these services.
- normal method with docker run
- fast method with docker-compose

but first we should build our own services
## build a new image
we need a dockerfile , now we create it:
```
docker build --no-cache --build-arg SERVICE={{SERVICE_NAME}} -t {{SERVICE_NAME}}:latest  . 

```
we provide better way to reach image :
```
./build.sh {{SERVICE_NAME}}

```
a Dockerfile must provided behind of build.sh file

## Dockerfile
Dockerfile contain following instructions :
```
FROM golang:1.17
```
this mean we want to have image with base golang version 1.17
```
ARG SERVICE=default
ENV SERVICEN=$SERVICE
```
get SERVICE variable that set in ``` docker build ``` command and set it to env variable
```
RUN git clone --depth 1 https://github.com/hofarah/cloud_project_1400.git /go/src/cloud_1400 --branch master
```
get files of repository
```
ENV GO111MODULE=on
```
we need set go modules to on for our build

```
WORKDIR /go/src/cloud_1400
```
set our workdir to /go/src/cloud_1400 that we set it in last step 
```
RUN go mod tidy
RUN go mod download
```
tidy and download our dependancy libraries.

```
RUN cd services/$SERVICE && go build -o main .
```
build our service
```
EXPOSE 7575
CMD /go/src/cloud_1400/services/${SERVICEN##*/}/main
```                                                           
expose on port 7575 and set command that run when container start



**now you have image of your service.**

## 1. docker run method
Follow the instructions below :

## 2. docker-compose method
in this method we can run multiple container with one command. we can use ```docker-compose up -d```. this command will automaticly run and manage the container we mentioned in a docker-compose.yml 

it works simply and if you want to remove the containers you can use ```docker-compose down```. it will stop and remove containers. docker-compose is like rum method but you can use it many times.


