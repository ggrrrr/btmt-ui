# ARG SERVICE
FROM node:lts-alpine as builder

WORKDIR /app

COPY package*.json ./

RUN npm install
COPY . .

RUN npm run build


FROM nginx:stable-alpine as production-stage
COPY --from=builder /app/dist /usr/share/nginx/html
COPY ./nginx/prod.conf  /etc/nginx/nginx.conf
EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
