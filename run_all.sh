#!/bin/bash

echo "Running run_auth.sh..."
gnome-terminal -- /bin/sh -c './run_auth.sh; sleep infinity'
echo "Running run_backend.sh..."
gnome-terminal -- /bin/sh -c './run_backend.sh; sleep infinity'
echo "Running run_collab.sh..."
gnome-terminal -- /bin/sh -c './run_collab.sh; sleep infinity'
echo "Running run_frontend.sh..."
gnome-terminal -- /bin/sh -c './run_frontend.sh; sleep infinity'
echo "Running run_website.sh..."
gnome-terminal -- /bin/sh -c './run_website.sh; sleep infinity'
