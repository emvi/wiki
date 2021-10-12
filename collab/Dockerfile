FROM node:14-alpine
RUN apk update && apk upgrade
WORKDIR /app
COPY collab /app

# build
RUN npm install && npm run compile && \
	rm -r src .gitignore package-lock.json package.json

RUN addgroup -S appuser && \
    adduser -S -G appuser appuser && \
    chown -R appuser:appuser /app
USER appuser

# default config
ENV EMVI_WIKI_LOGLEVEL=info

# run
EXPOSE 4004
CMD ["node", "dist/main.js"]
