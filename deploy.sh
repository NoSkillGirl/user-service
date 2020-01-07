#!/bin/bash
echo "Starting Deployment"

echo "1. Make a go build which is compatible with ubuntu"
echo "--------------------------------------------------"
GOOS=linux GOARCH=amd64 go build -v .

echo "2. Make a docker image and push it to docker hub"
echo "--------------------------------------------------"
VERSION=$(cat VERSION) 
let "NEW_VERSION=VERSION+1"

echo "Current Version: $VERSION"
echo "New Version: $NEW_VERSION"

echo $NEW_VERSION > VERSION

docker build -t ajoop/user-service:v$NEW_VERSION .
docker push ajoop/user-service:v$NEW_VERSION

echo "3. Make a manifest file for new version"
echo "--------------------------------------------------"
/bin/bash -c "sed -i '' 's/v$VERSION/v$NEW_VERSION/g' deploy-app.yaml"

cat deploy-app.yaml

echo "4. kubectl apply"
echo "--------------------------------------------------"
kubectl apply -f deploy-app.yaml

echo "5. Removing the generated go build"
echo "--------------------------------------------------"
rm user-service

echo "Successfully Deployed"