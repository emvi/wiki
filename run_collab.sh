#!/bin/bash

cd collab

export EMVI_WIKI_LOG_LEVEL=debug
export EMVI_WIKI_LOG_DISABLE_TIMESTAMP=true
export EMVI_WIKI_BACKEND_HOST=http://localhost.com:4003
export EMVI_WIKI_AUTH_HOST=http://localhost.com:4001
export EMVI_WIKI_AUTH_CLIENT_ID=xz5tN33UW6kZzIyHrO8x
export EMVI_WIKI_AUTH_CLIENT_SECRET=3DLvZUadCoy6xB9mHx5uaJkfsY3K7WjwVBW6pdAeTAEy37Bjf4u05By7e0QT0hBT

npm run compile && npm start
