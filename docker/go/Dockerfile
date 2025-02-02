ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk add --update --no-cache git fish curl less build-base tzdata vim \
  && go install github.com/air-verse/air@latest \
  && go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app
VOLUME /app

ENV PORT=3000

ARG GROUP_ID
ARG USER_ID
ARG user=app
RUN addgroup -g $GROUP_ID -S $user \
  && adduser --disabled-password --gecos '' --uid $USER_ID --ingroup $user $user \
  && chown -R $user:$user /app /go

USER $user
EXPOSE $PORT
CMD ["air"]


