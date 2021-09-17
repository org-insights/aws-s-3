#!/bin/sh

set -x

mc alias set s3 "${S3_ENDPOINT}" "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}" --api S3v4

mc mb s3/"${BUCKET}" --ignore-existing;

for day in 01 02 03 04 05 06 07 08 09 10
do
  echo "$day" > /temp_file.txt
  mc cp /temp_file.txt s3/"${BUCKET}"/client=1000/date=2021-09-$day/temp_file.txt
  rm /temp_file.txt
done
