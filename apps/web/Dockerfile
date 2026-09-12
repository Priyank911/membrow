FROM node:22-bookworm-slim AS build
WORKDIR /app
COPY apps/web/package.json apps/web/package-lock.json* ./
RUN npm install
COPY apps/web .
RUN npm run build

FROM node:22-bookworm-slim
WORKDIR /app
COPY --from=build /app/build ./build
COPY --from=build /app/package.json ./package.json
RUN npm install --omit=dev
ENV HOST=0.0.0.0
ENV PORT=5173
EXPOSE 5173
CMD ["node", "build"]
