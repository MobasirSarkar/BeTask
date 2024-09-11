FROM node:22-alpine3.20 AS development

RUN mkdir /app
WORKDIR /app


COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm@latest
RUN pnpm install

COPY . .
EXPOSE 5173

CMD [ "pnpm","run","dev" ]
