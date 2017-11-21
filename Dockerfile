FROM cpdpro/baseimage:alpine-narada4d

COPY crontab /app/crontab
COPY service /app/service
RUN set -ex -o pipefail; \
    ln -s /etc/sv/dcron /etc/service/cron; \
    install -m 0600 /app/crontab /etc/crontabs/app; \
    echo app >> /etc/crontabs/cron.update; \
    ln -nsf /app/service/* /etc/service/

COPY schema  /app/schema
COPY bin /usr/local/bin/
