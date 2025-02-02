ARG NODEJS_VERSION
FROM node:${NODEJS_VERSION}-alpine

RUN apk add --update --no-cache git fish curl less build-base tzdata vim yarn yq \
  && npm i -g swagger2openapi

WORKDIR /app
VOLUME /app

USER node
CMD ["yarn", "run", "dev"]
