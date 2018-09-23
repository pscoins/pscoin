FROM golang:1.10.1 AS builder

ENV COCKROACH_VERSION cockroach-v2.0.5.linux-amd64
ENV COCKROACH_URL https://binaries.cockroachdb.com/$COCKROACH_VERSION.tgz

ENV PATH /usr/local/go/bin:/go/bin:/usr/local/bin:$PATH


RUN wget -qO- $COCKROACH_URL | tar xvz
RUN cp -i $COCKROACH_VERSION/cockroach /usr/local/bin
#RUN cockroach version


COPY pscoin /usr/local/bin
COPY scripts/config-roach.sh /scripts/config-roach.sh
COPY scripts/roach-run.sh /scripts/roach-run.sh

RUN chmod +x /scripts/config-roach.sh
RUN chmod +x /scripts/roach-run.sh


CMD ["pscoin"]
