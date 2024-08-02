#!/bin/bash

echo "GCP Emulator initializing..."
gcloud beta emulators pubsub start --host-port=0.0.0.0:8686 &
PUBSUB_PID=$!

export PUBSUB_EMULATOR_HOST=localhost:8686
gcloud config configurations create emulator
gcloud config set auth/disable_credentials true
gcloud config set project test-project
gcloud config set api_endpoint_overrides/pubsub http://localhost:8686/

echo "Pub/Sub initializing..."
/scripts/init_pubsub.sh

echo "Storage initializing..."
/scripts/init_storage.sh

echo "GCP emulator started"

wait $PUBSUB_PID
