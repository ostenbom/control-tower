#!/bin/bash

[ "$VERBOSE" ] && { set -x; export BOSH_LOG_LEVEL=debug; }
set -euo pipefail

deployment="systest-github-$RANDOM"

cleanup() {
  status=$?
  ./cup-new --non-interactive destroy $deployment
  exit $status
}
set +u
if [ -z "$SKIP_TEARDOWN" ]; then
  trap cleanup EXIT
else
  trap "echo Skipping teardown" EXIT
fi
set -u

cp "$BINARY_PATH" ./cup
chmod +x ./cup

echo "DEPLOY WITH GITHUB FLAGS"

./cup deploy $deployment \
  --github-auth-client-id "$GITHUB_AUTH_CLIENT_ID" \
  --github-auth-client-secret "$GITHUB_AUTH_CLIENT_SECRET" \
  --domain cup.engineerbetter.com

config=$(./cup info --json $deployment)
domain=$(echo "$config" | jq -r '.config.domain')
username=$(echo "$config" | jq -r '.config.concourse_username')
password=$(echo "$config" | jq -r '.config.concourse_password')
echo "$config" | jq -r '.config.concourse_ca_cert' > generated-ca-cert.pem

fly --target system-test login \
  --ca-cert generated-ca-cert.pem \
  --concourse-url "https://$domain" \
  --username "$username" \
  --password "$password"

fly --target system-test set-team \
  --team-name=git-team \
  --github-user=EngineerBetterCI

fly --target system-test login \
  --team-name=git-team 2>&1 >fly_out &

pkill -9 fly

url="$(grep redirect fly_out)"

curl -sL "$url" | grep -q '/sky/issuer/auth/github'

echo "TEST SUCCESSFUL"