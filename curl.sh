#!/bin/bash

for i in {1..100};
do
echo "Sending $i"
sleep 2
curl --location --request POST 'http://localhost:1323/packages/location/123' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": "Sylhet",
    "to": "Gafargaon",
    "vehicle_id": "123"
}'
done