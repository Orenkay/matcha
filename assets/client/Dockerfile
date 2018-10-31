FROM node:latest

COPY . /usr/app
WORKDIR /usr/app

RUN npm install
RUN npm run build
ENTRYPOINT [ "node", "serve.js" ]
