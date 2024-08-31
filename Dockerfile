FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git bash curl gcc g++ sqlite-dev

RUN rm -rf /var/cache/apk/

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o /app/harvester ./cmd/harvest

FROM alpine:3.18

RUN apk add --no-cache git bash curl

COPY --from=builder /app/harvester /app/harvester
COPY --from=builder /app/public /app/public
COPY --from=builder /app/harvest_and_push.sh /app/harvest_and_push.sh
COPY --from=builder /app/.env /app/.env

COPY crontab /etc/crontabs/root

RUN chmod +x /app/harvest_and_push.sh

ENV GITHUB_TOKEN=${GITHUB_TOKEN}
ENV GITHUB_USERNAME=${GITHUB_USERNAME}
ENV GITHUB_REPOSITORY=${GITHUB_REPOSITORY}
ENV GIT_COMMITTER_NAME=${GIT_COMMITTER_NAME}
ENV GIT_COMMITTER_EMAIL=${GIT_COMMITTER_EMAIL}

WORKDIR /app

CMD ["crond", "-f"]
