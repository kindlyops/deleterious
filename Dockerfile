FROM golang:1.12-alpine

LABEL "com.github.actions.name"="Deleterious CI Action"
LABEL "com.github.actions.description"="Enables CI/CD for Deleterious project"
# https://developer.github.com/actions/creating-github-actions/creating-a-docker-container/#supported-feather-icons
LABEL "com.github.actions.icon"="truck"
LABEL "com.github.actions.color"="gray-dark"

LABEL repository="https://github.com/kindlyops/deleterious"
LABEL homepage="https://kindlyops.com/knowledge-base/github-actions"
LABEL maintainer="Kindly Ops <support@kindlyops.com>"

RUN apk add --no-cache bash \
    curl \
    docker \
    git \
    mercurial \
    rpm

ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV GORELEASER_VERSION v0.111.0
ENV CGO_ENABLED=0

RUN wget https://github.com/goreleaser/goreleaser/releases/download/$GORELEASER_VERSION/goreleaser_Linux_x86_64.tar.gz \
    && tar -C /bin -xzvf goreleaser_Linux_x86_64.tar.gz \
    && rm goreleaser_Linux_x86_64.tar.gz
