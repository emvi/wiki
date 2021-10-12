FROM golang AS build
RUN apt-get update && \
	apt-get upgrade -y && \
	apt-get install -y curl && \
	curl -sL https://deb.nodesource.com/setup_14.x -o nodesource_setup.sh && bash nodesource_setup.sh && \
	apt-get install -y nodejs
WORKDIR /go/src/emviwiki
COPY frontend ./frontend
COPY shared ./shared
COPY go.mod ./go.mod
COPY go.sum ./go.sum

# build server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" /go/src/emviwiki/frontend/main.go && \
	mkdir /app && mv /go/src/emviwiki/main /app/server

# build frontend
RUN cd /go/src/emviwiki/frontend/public && bash -c "npm i" && bash -c "npm rebuild node-sass" && bash -c "npm run build"

# cleanup
RUN rm -r /go/src/emviwiki/frontend/public/node_modules \
	/go/src/emviwiki/frontend/public/src \
	/go/src/emviwiki/frontend/public/.babelrc \
	/go/src/emviwiki/frontend/public/package-lock.json \
	/go/src/emviwiki/frontend/public/package.json \
	/go/src/emviwiki/frontend/public/webpack.config.js && \
	mv /go/src/emviwiki/frontend/public /app/public && \
    rm -r /go

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=build /app /app
WORKDIR /app

RUN addgroup -S appuser && \
    adduser -S -G appuser appuser && \
    chown -R appuser:appuser /app
USER appuser

# default config
ENV EMVI_WIKI_HOST=0.0.0.0:4003
ENV EMVI_WIKI_LOGLEVEL=info
ENV EMVI_WIKI_ALLOWED_ORIGINS=*
ENV EMVI_WIKI_TLS_ENABLE=false
ENV EMVI_WIKI_AUTH_HOST=https://auth.emvi.com
ENV EMVI_WIKI_WEBSITE_HOST=https://emvi.com
ENV EMVI_WIKI_BACKEND_HOST=https://api.emvi.com
ENV EMVI_WIKI_COLLAB_HOST=https://api.emvi.com
ENV EMVI_WIKI_HTTP_WRITE_TIMEOUT=5
ENV EMVI_WIKI_HTTP_READ_TIMEOUT=5

# run
EXPOSE 4003
VOLUME ["/logs"]
ENTRYPOINT ["/app/server"]
