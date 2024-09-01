FROM golang:1.22-alpine AS builder

RUN apk add gcc g++ sqlite-dev

RUN rm -rf /var/cache/apk/

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /app/harvest ./cmd/harvest/main.go

FROM alpine:3.18

RUN apk add --no-cache git bash curl

COPY --from=builder /app/harvest /app/harvest
COPY --from=builder /app/public /app/public
COPY --from=builder /app/run.sh /app/run.sh
COPY --from=builder /app/.env /app/.env

COPY crontab /etc/crontabs/root

RUN chmod +x /app/run.sh

ENV GITHUB_TOKEN=${GITHUB_TOKEN}
ENV GITHUB_USERNAME=${GITHUB_USERNAME}
ENV GITHUB_REPOSITORY=${GITHUB_REPOSITORY}
ENV GIT_COMMITTER_NAME=${GIT_COMMITTER_NAME}
ENV GIT_COMMITTER_EMAIL=${GIT_COMMITTER_EMAIL}

CMD ["crond", "-f"]
