#!/bin/sh

set -x

mc alias set s3 "${S3_ENDPOINT}" "${MINIO_ACCESS_KEY}" "${MINIO_SECRET_KEY}" --api S3v4

mc mb s3/"${BUCKET}" --ignore-existing;

for day in 01 02 03 04 05 06 07 08 09 10
do
  dd if=/dev/zero of=temp_file.data  bs=${day}K  count=1
  mc cp /temp_file.data s3/"${BUCKET}"/client=1000/date=2021-09-${day}/temp_file.data
  rm /temp_file.data
done

for day in 01 02 03 04 05 
do
  for hour in 00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20 21 22 23
  do
    dd if=/dev/zero of=temp_file.data  bs=1K  count=1
    mc cp /temp_file.data s3/"${BUCKET}"/client=2000/date=2021-09-${day}/hour=${hour}/temp_file.data
    rm /temp_file.data
  done
done