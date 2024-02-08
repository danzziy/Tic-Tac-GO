FROM node:lts-alpine

WORKDIR /app

copy ./frontend/ ./

RUN npm install

EXPOSE 8080

CMD ["npm", "run", "serve"]
