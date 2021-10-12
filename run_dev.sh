#!/bin/bash

echo "Running node-sass and dev..."
cd frontend/public/
gnome-terminal -- /bin/sh -c 'npm run style; sleep infinity'
gnome-terminal -- /bin/sh -c 'npm run dev; sleep infinity'
cd ..
cd ..
echo "Running design node-sass..."
cd design/
gnome-terminal -- /bin/sh -c 'npm run dev; sleep infinity'
cd ..
echo "Running build_frontend.sh..."
gnome-terminal -- /bin/sh -c './build_frontend.sh; sleep infinity'
echo "Running build_website.sh..."
gnome-terminal -- /bin/sh -c './build_website.sh; sleep infinity'
echo "Running build_website.sh..."
gnome-terminal -- /bin/sh -c './build_website.sh; sleep infinity'
