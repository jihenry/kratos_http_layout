FROM d-harbor.lepass.cn/public/go:1.17.12
LABEL maintainer="qasim"
ENV LANG C.UTF-8
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone
RUN mkdir /app
WORKDIR /app
ADD . .
ENTRYPOINT ["sh","/app/entrypoint.sh"]
