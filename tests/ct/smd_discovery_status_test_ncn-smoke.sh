#!/bin/bash -l

# MIT License

# (C) Copyright [2021] Hewlett Packard Enterprise Development LP

# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

# HMS test metrics test cases: 3
# 1. GET /Inventory/RedfishEndpoints API response code
# 2. GET /Inventory/RedfishEndpoints API response body
# 3. Verify Redfish endpoint discovery statuses

# initialize test variables
SMOKE_TEST_LIB="/opt/cray/tests/ncn-resources/hms/hms-test/hms_smoke_test_lib_ncn-resources_remote-resources.sh"
TARGET="api-gw-service-nmn.local"

trap ">&2 echo \"recieved kill signal, exiting with status of '1'...\" ; \
    exit 1" SIGHUP SIGINT SIGTERM

# source HMS smoke test library file
if [[ -r ${SMOKE_TEST_LIB} ]] ; then
    . ${SMOKE_TEST_LIB}
else
    >&2 echo "ERROR: failed to source HMS smoke test library: ${SMOKE_TEST_LIB}"
    exit 1
fi

# check for jq dependency
WHICH_CMD="which jq"
WHICH_OUT=$(eval ${WHICH_CMD})
WHICH_RET=$?
if [[ ${WHICH_RET} -ne 0 ]] ; then
    echo "${WHICH_OUT}"
    >&2 echo "ERROR: '${WHICH_CMD}' failed with status code: ${WHICH_RET}"
    exit 1
fi

echo "Running smd_discovery_status_test..."

# retrieve Keycloak authentication token for session
TOKEN=$(get_auth_access_token)
TOKEN_RET=$?
if [[ ${TOKEN_RET} -ne 0 ]] ; then
    exit 1
fi

# query SMD for the Redfish endpoint discovery statuses
CURL_CMD="curl -s -k -H \"Authorization: Bearer ${TOKEN}\" https://${TARGET}/apis/smd/hsm/v1/Inventory/RedfishEndpoints"
timestamp_print "Testing '${CURL_CMD}'..."
CURL_OUT=$(eval ${CURL_CMD})
CURL_RET=$?
if [[ ${CURL_RET} -ne 0 ]] ; then
    >&2 echo "ERROR: '${CURL_CMD}' failed with status code: ${CURL_RET}"
    exit 1
elif [[ -z "${CURL_OUT}" ]] ; then
    >&2 echo "ERROR: '${CURL_CMD}' returned an empty response."
    exit 1
fi

# parse the SMD response
JQ_CMD="jq '.RedfishEndpoints[] | { ID: .ID, LastDiscoveryStatus: .DiscoveryInfo.LastDiscoveryStatus}' -c | sort -V | jq -c"
timestamp_print "Processing response with: '${JQ_CMD}'..."
PARSED_OUT=$(echo "${CURL_OUT}" | eval "${JQ_CMD}" 2> /dev/null)
if [[ -z "${PARSED_OUT}" ]] ; then
    echo "${CURL_OUT}"
    >&2 echo "ERROR: '${CURL_CMD}' returned a response with missing endpoint IDs or LastDiscoveryStatus fields"
    exit 1
fi

# sanity check the response body
while read LINE ; do
    ID_CHECK=$(echo "${LINE}" | grep -E "\"ID\"")
    if [[ -z "${ID_CHECK}" ]] ; then
        echo "${LINE}"
        >&2 echo "ERROR: '${CURL_CMD}' returned a response with missing endpoint ID fields"
        exit 1
    fi
    STATUS_CHECK=$(echo "${LINE}" | grep -E "\"LastDiscoveryStatus\"")
    if [[ -z "${STATUS_CHECK}" ]] ; then
        echo "${LINE}"
        >&2 echo "ERROR: '${CURL_CMD}' returned a response with missing endpoint LastDiscoveryStatus fields"
        exit 1
    fi
done <<< "${PARSED_OUT}"

# verify that at least one endpoint was discovered successfully
PARSED_CHECK=$(echo "${PARSED_OUT}" | grep -E "ID.*LastDiscoveryStatus.*DiscoverOK")
if [[ -z "${PARSED_CHECK}" ]] ; then
    echo "${PARSED_OUT}"
    echo "FAIL: smd_discovery_status_test found no successfully discovered endpoints"
    exit 1
fi

# count the number of endpoints with unexpected discovery statuses
timestamp_print "Verifying endpoint discovery statuses..."
PARSED_FAILED=$(echo "${PARSED_OUT}" | grep -v "DiscoverOK")
NUM_FAILS=$(echo "${PARSED_FAILED}" | grep -E "ID.*LastDiscoveryStatus" | wc -l)
# one endpoint on the site network is expected to be unreachable and fail discovery with a status of 'HTTPSGetFailed'
if [[ ${NUM_FAILS} -gt 1 ]] ; then
    echo "${PARSED_FAILED}"
    echo "FAIL: smd_discovery_status_test found ${NUM_FAILS} endpoints that failed discovery, maximum allowable is 1"
    exit 1
else
    echo "PASS: smd_discovery_status_test passed!"
    exit 0
fi
