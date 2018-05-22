FROM postgres:9.6.1
MAINTAINER <miaolizhao@126.com>

ENV smartcooly_VERSION 0.1.2
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates wget && rm -rf /var/lib/apt/lists/* && \
    wget -q https://github.com/marstau/smartcooly/releases/download/v${smartcooly_VERSION}/smartcooly_linux_amd64.tar.gz && \
    tar -C /usr/src -xzf smartcooly_linux_amd64.tar.gz && \
    echo $'nohup /docker-entrypoint.sh postgres >/dev/null 2>&1 &\n\
sleep 30s\n\
cd /usr/src/smartcooly_linux_amd64/\n\
./smartcooly\n' >> /usr/src/smartcooly_linux_amd64/cmd.sh && \
    chmod +x /usr/src/smartcooly_linux_amd64/smartcooly && \
    chmod +x /usr/src/smartcooly_linux_amd64/cmd.sh && \
    rm -f smartcooly_linux_amd64.tar.gz && \
    apt-get purge -y --auto-remove wget

EXPOSE 9876

VOLUME /var/lib/postgresql/data

COPY config.ini /usr/src/smartcooly_linux_amd64/custom/config.ini

CMD ["/usr/src/smartcooly_linux_amd64/cmd.sh"]
