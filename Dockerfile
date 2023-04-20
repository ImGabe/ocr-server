FROM alpine:3.14 as base
LABEL maintainer="otiai10 <otiai10@gmail.com>" \
	 updatedBy="RocSun <oldsixa@163.com>"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add tesseract-ocr \
    && apk add tesseract-ocr-dev \
    && apk add tesseract-ocr-data-chi_sim

FROM base AS build

RUN apk add go \
      g++

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn,direct"
ENV GOPATH=/go
ENV PATH=${PATH}:${GOPATH}/bin

ADD . $GOPATH/src/ocr-server
WORKDIR $GOPATH/src/ocr-server
RUN go get -v ./... && go install .

FROM base AS OS
RUN apk add bash
COPY --from=build /go /go
ENV PORT=8080
CMD ["/go/bin/ocr-server"]