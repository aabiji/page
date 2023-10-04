#!/bin/bash
export EPUB_EXTRACT_DIRECTORY=""
export FILE_UPLOAD_DIRECTORY=""
cd backend &&
go build &&
./page &
cd frontend &&
npm run build &&
npm run preview