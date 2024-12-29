#!/bin/bash

# Generate handlers from the OpenAPI specification using ogen-go.
go run github.com/ogen-go/ogen/cmd/ogen@latest --target api --clean openapi.yaml

# The browser Beacon API sends all content in text/plain even if the content is in JSON.
# As ogen-go validates the Content-Type header from our OpenAPI specification, this is a
# workaround to patch api/oas_request_decoders_gen.go to add text/plain to the list of
# accepted content types for the /event/hit endpoint and still use the same JSON decoder.

# This runs perl and replaces a line number with a new line.
# This should be verified after each update to the OpenAPI specification or ogen-go
# in case the line number changes.
#
# perl is more portable across different systems compared to sed.
# Line 256
perl -i -pe 's/^.*$// if $. == 256; $. == 256 and print "    case ct == \"application/json\", ct == \"text/plain\":"' ./api/oas_request_decoders_gen.go
