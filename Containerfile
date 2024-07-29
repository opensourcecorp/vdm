# This Containerfile isn't mean for creating artifacts etc., it's just a way to
# perform portable, local CI checks in case there are workstation-specific
# issues a developer faces.
FROM docker.io/library/golang:1.20

RUN apt-get update && apt-get install -y \
      ca-certificates \
      make \
      sudo

COPY . /go/app
WORKDIR /go/app

RUN make ci
