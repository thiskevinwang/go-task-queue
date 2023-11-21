FROM node:18-bookworm as builder
WORKDIR /app
COPY /apps/web /app/apps/web
COPY package.json /app
COPY package-lock.json /app
RUN npm install

FROM node:18-bookworm
WORKDIR /app
COPY --from=builder /app/node_modules /app/node_modules
COPY . /app

EXPOSE 3000

CMD ["npm", "run", "start", "-w", "apps/web"]
