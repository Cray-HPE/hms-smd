#!/bin/bash -l
#
# Copyright 2020 Hewlett Packard Enterprise Development LP
#
###############################################################
#
#     CASM Test - Cray Inc.
#
#     TEST IDENTIFIER   : hsm_smoke_test
#
#     DESCRIPTION       : Automated test for verifying basic HSM/SMD API
#                         infrastructure and installation on Cray Shasta
#                         systems.
#                         
#     AUTHOR            : Mitch Schooler
#
#     DATE STARTED      : 09/22/2020
#
#     LAST MODIFIED     : 11/18/2020
#
#     SYNOPSIS
#       This is a smoke test for the HMS HSM/SMD API that makes basic HTTP
#       requests using curl to verify that the service's API endpoints
#       respond and function as expected after an installation.
#
#     INPUT SPECIFICATIONS
#       Usage: hsm_smoke_test
#       
#       Arguments: None
#
#     OUTPUT SPECIFICATIONS
#       Plaintext is printed to stdout and/or stderr. The script exits
#       with a status of '0' on success and '1' on failure.
#
#     DESIGN DESCRIPTION
#       This smoke test is based on the Shasta health check srv_check.sh
#       script in the CrayTest repository that verifies the basic health of
#       various microservices but instead focuses exclusively on the HSM/SMD
#       API. It was implemented to run from the ct-pipelines container off
#       of the NCN of the system under test within the DST group's Continuous
#       Testing (CT) framework as part of the remote-smoke test suite.
#
#     SPECIAL REQUIREMENTS
#       Must be executed from the ct-pipelines container on a remote host
#       (off of the NCNs of the test system) with the Continuous Test
#       infrastructure installed.
#
#     UPDATE HISTORY
#       user       date         description
#       -------------------------------------------------------
#       schooler   09/22/2020   initial implementation
#       schooler   11/18/2020   remove deprecated HSNInterfaces test
#
#     DEPENDENCIES
#       - hms_smoke_test_lib_ncn-resources_remote-resources.sh which is
#         expected to be packaged in
#         /opt/cray/tests/remote-resources/hms/hms-test in the ct-pipelines
#         container.
#
#     BUGS/LIMITATIONS
#       None
#
###############################################################

# HMS test metrics test cases: 31
# 1. Check cray-smd pod statuses
# 2. GET /service/ready API response code
# 3. GET /service/liveness API response code
# 4. GET /service/values API response code
# 5. GET /service/values/arch API response code
# 6. GET /service/values/class API response code
# 7. GET /service/values/flag API response code
# 8. GET /service/values/nettype API response code
# 9. GET /service/values/role API response code
# 10. GET /service/values/subrole API response code
# 11. GET /service/values/state API response code
# 12. GET /service/values/type API response code
# 13. GET /State/Components API response code
# 14. GET /Defaults/NodeMaps API response code
# 15. GET /Inventory/DiscoveryStatus API response code
# 16. GET /Inventory/Hardware API response code
# 17. GET /Inventory/HardwareByFRU API response code
# 18. GET /Inventory/Hardware/History API response code
# 19. GET /Inventory/HardwareByFRU/History API response code
# 20. GET /Inventory/RedfishEndpoints API response code
# 21. GET /Inventory/ComponentEndpoints API response code
# 22. GET /Inventory/ServiceEndpoints API response code
# 23. GET /Inventory/EthernetInterfaces API response code
# 24. GET /groups API response code
# 25. GET /groups/labels API response code
# 26. GET /partitions API response code
# 27. GET /partitions/names API response code
# 28. GET /memberships API response code
# 29. GET /Subscriptions/SCN API response code
# 30. GET /locks API response code
# 31. GET /sysinfo/powermaps API response code

# initialize test variables
TEST_RUN_TIMESTAMP=$(date +"%Y%m%dT%H%M%S")
TEST_RUN_SEED=${RANDOM}
OUTPUT_FILES_PATH="/tmp/hsm_smoke_test_out-${TEST_RUN_TIMESTAMP}.${TEST_RUN_SEED}"
SMOKE_TEST_LIB="/opt/cray/tests/remote-resources/hms/hms-test/hms_smoke_test_lib_ncn-resources_remote-resources.sh"
CURL_ARGS="-k -i -s -S"
MAIN_ERRORS=0
CURL_COUNT=0

# cleanup
function cleanup()
{
    echo "cleaning up..."
    rm -f ${OUTPUT_FILES_PATH}*
}

# main
function main()
{
    AUTH_ARG="-H \"Authorization: Bearer $TOKEN\""

    # GET tests
    for URL_ARGS in \
        "apis/smd/hsm/v1/service/ready" \
        "apis/smd/hsm/v1/service/liveness" \
        "apis/smd/hsm/v1/service/values" \
        "apis/smd/hsm/v1/service/values/arch" \
        "apis/smd/hsm/v1/service/values/class" \
        "apis/smd/hsm/v1/service/values/flag" \
        "apis/smd/hsm/v1/service/values/nettype" \
        "apis/smd/hsm/v1/service/values/role" \
        "apis/smd/hsm/v1/service/values/subrole" \
        "apis/smd/hsm/v1/service/values/state" \
        "apis/smd/hsm/v1/service/values/type" \
        "apis/smd/hsm/v1/State/Components" \
        "apis/smd/hsm/v1/Defaults/NodeMaps" \
        "apis/smd/hsm/v1/Inventory/DiscoveryStatus" \
        "apis/smd/hsm/v1/Inventory/Hardware" \
        "apis/smd/hsm/v1/Inventory/HardwareByFRU" \
        "apis/smd/hsm/v1/Inventory/Hardware/History" \
        "apis/smd/hsm/v1/Inventory/HardwareByFRU/History" \
        "apis/smd/hsm/v1/Inventory/RedfishEndpoints" \
        "apis/smd/hsm/v1/Inventory/ComponentEndpoints" \
        "apis/smd/hsm/v1/Inventory/ServiceEndpoints" \
        "apis/smd/hsm/v1/Inventory/EthernetInterfaces" \
        "apis/smd/hsm/v1/groups" \
        "apis/smd/hsm/v1/groups/labels" \
        "apis/smd/hsm/v1/partitions" \
        "apis/smd/hsm/v1/partitions/names" \
        "apis/smd/hsm/v1/memberships" \
        "apis/smd/hsm/v1/Subscriptions/SCN" \
        "apis/smd/hsm/v1/locks" \
        "apis/smd/hsm/v1/sysinfo/powermaps"
    do
        URL=$(url "${URL_ARGS}")
        URL_RET=$?
        if [[ ${URL_RET} -ne 0 ]] ; then
            cleanup
            exit 1
        fi
        run_curl "GET ${AUTH_ARG} ${URL}"
        if [[ $? -ne 0 ]] ; then
            ((MAIN_ERRORS++))
        fi
    done

    echo "MAIN_ERRORS=${MAIN_ERRORS}"
    return ${MAIN_ERRORS}
}

# check_pod_status
function check_pod_status()
{
    run_check_pod_status "cray-smd"
    return $?
}

# TARGET_SYSTEM is expected to be set in the ct-pipelines container
if [[ -z ${TARGET_SYSTEM} ]] ; then
    >&2 echo "ERROR: TARGET_SYSTEM environment variable is not set"
    cleanup
    exit 1
else
    echo "TARGET_SYSTEM=${TARGET_SYSTEM}"
    TARGET="auth.${TARGET_SYSTEM}"
    echo "TARGET=${TARGET}"
fi

# TOKEN is expected to be set in the ct-pipelines container
if [[ -z ${TOKEN} ]] ; then
    >&2 echo "ERROR: TOKEN environment variable is not set"
    cleanup
    exit 1
else
    echo "TOKEN=${TOKEN}"
fi

trap ">&2 echo \"recieved kill signal, exiting with status of '1'...\" ; \
    cleanup ; \
    exit 1" SIGHUP SIGINT SIGTERM

# source HMS smoke test library file
if [[ -r ${SMOKE_TEST_LIB} ]] ; then
    . ${SMOKE_TEST_LIB}
else
    >&2 echo "ERROR: failed to source HMS smoke test library: ${SMOKE_TEST_LIB}"
    exit 1
fi

# make sure filesystem is writable for output files
touch ${OUTPUT_FILES_PATH}
if [[ $? -ne 0 ]] ; then
    >&2 echo "ERROR: output file location not writable: ${OUTPUT_FILES_PATH}"
    cleanup
    exit 1
fi

echo "Running hsm_smoke_test..."

# run initial pod status test
check_pod_status
if [[ $? -ne 0 ]] ; then
    echo "FAIL: hsm_smoke_test ran with failures"
    cleanup
    exit 1
fi

# run main API tests
main
if [[ $? -ne 0 ]] ; then
    echo "FAIL: hsm_smoke_test ran with failures"
    cleanup
    exit 1
else
    echo "PASS: hsm_smoke_test passed!"
    cleanup
    exit 0
fi
