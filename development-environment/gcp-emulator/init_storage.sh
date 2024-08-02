#!/bin/bash

export STORAGE_EMULATOR_HOST=http://localhost:8080

gsutil mb -p test-project gs://my-bucket
