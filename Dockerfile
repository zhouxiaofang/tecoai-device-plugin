FROM alpine:latest
MAINTAINER zhoufang <zhoufang@qq.com>

COPY ./build/tecoai /root/tecoai

CMD ["/root/tecoai"]
