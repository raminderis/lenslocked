FROM node:latest

WORKDIR /tailwind

RUN npm init -y
RUN npm install @tailwindcss/cli

CMD npx @tailwindcss/cli \
    -i /src/styles.css \
    -o /dst/output.css \
    --watch --poll