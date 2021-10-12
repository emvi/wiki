#!/bin/bash

rm -r website/public/static/*
#rm -r frontend/public/static/*
rm -r auth/static/*
cp design/dist/website.css website/public/static/website.css
cp design/dist/blog.css website/public/static/blog.css
#cp design/dist/frontend.css frontend/public/static/frontend.css
cp design/dist/website.css auth/static/website.css
cp -a design/static/ website/public/
#cp -a design/static/ frontend/public/
cp -a design/static/ auth/
