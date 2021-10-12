#!/bin/bash

VERSION=${1-latest}

echo "Building docker images version $VERSION"
docker build -t emviwiki-auth -f auth/Dockerfile .
docker build -t emviwiki-backend -f backend/Dockerfile .
docker build -t emviwiki-batch -f batch/Dockerfile .
docker build -t emviwiki-collab -f collab/Dockerfile .
docker build -t emviwiki-dashboard -f dashboard/Dockerfile .
docker build -t emviwiki-frontend -f frontend/Dockerfile .
docker build -t emviwiki-website -f website/Dockerfile .

echo "Tagging docker images version $VERSION"
docker tag emviwiki-auth "registry.emvi-integration.com/emviwiki-auth:$VERSION"
docker tag emviwiki-backend "registry.emvi-integration.com/emviwiki-backend:$VERSION"
docker tag emviwiki-batch "registry.emvi-integration.com/emviwiki-batch:$VERSION"
docker tag emviwiki-collab "registry.emvi-integration.com/emviwiki-collab:$VERSION"
docker tag emviwiki-dashboard "registry.emvi-integration.com/emviwiki-dashboard:$VERSION"
docker tag emviwiki-frontend "registry.emvi-integration.com/emviwiki-frontend:$VERSION"
docker tag emviwiki-website "registry.emvi-integration.com/emviwiki-website:$VERSION"

echo "Pushing docker images version $VERSION"
docker push "registry.emvi-integration.com/emviwiki-auth:$VERSION"
docker push "registry.emvi-integration.com/emviwiki-backend:$VERSION"
docker push "registry.emvi-integration.com/emviwiki-batch:$VERSION"
docker push "registry.emvi-integration.com/emviwiki-collab:$VERSION"
docker push "registry.emvi-integration.com/emviwiki-dashboard:$VERSION"
docker push "registry.emvi-integration.com/emviwiki-frontend:$VERSION"
docker push "registry.emvi-integration.com/emviwiki-website:$VERSION"
