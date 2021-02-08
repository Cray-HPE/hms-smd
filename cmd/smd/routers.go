/*
 * MIT License
 *
 * (C) Copyright [2018-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (s *SmD) NewRouter(routes []Route) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if s.lgLvl >= LOG_DEBUG ||
			(route.Name != "doReadyGet" && route.Name != "doLivenessGet") {
			handler = s.Logger(handler, route.Name)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.MethodNotAllowedHandler = http.HandlerFunc(s.doMethodNotAllowedHandler)
	s.router = router

	return router
}

func (s *SmD) getAllMethodsForRequest(req *http.Request) []string {
	var allMethods []string
	smdRoutes := s.generateRoutes()
	for _, smdRoute := range smdRoutes {
		route := s.router.Get(smdRoute.Name)
		if route != nil {
			var match mux.RouteMatch
			if route.Match(req, &match) || match.MatchErr == mux.ErrMethodMismatch {
				methods, err := route.GetMethods()
				if err == nil {
					allMethods = append(allMethods, methods...)
				}
			}
		}
	}
	return allMethods
}

func (s *SmD) doMethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	allowString := ""
	//need to dynamically build the list of HTTP methods supported for the failed request
	allMethods := s.getAllMethodsForRequest(r)
	if len(allMethods) > 0 {
		allowString = strings.Join(allMethods[:], ",")
	}
	if allowString != "" {
		w.Header().Set("allow", allowString)
	}
	sendJsonError(w, http.StatusMethodNotAllowed, "allow "+allowString)
}

func (s *SmD) generateRoutes() Routes {
	return Routes{
		// HSM Service State
		Route{
			"doReadyGet",
			strings.ToUpper("Get"),
			s.serviceBase + "/ready",
			s.doReadyGet,
		},
		Route{
			"doLivenessGet",
			strings.ToUpper("Get"),
			s.serviceBase + "/liveness",
			s.doLivenessGet,
		},
		Route{
			"doValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase,
			s.doValuesGet,
		},
		Route{
			"doArchValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/arch",
			s.doArchValuesGet,
		},
		Route{
			"doClassValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/class",
			s.doClassValuesGet,
		},
		Route{
			"doFlagValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/flag",
			s.doFlagValuesGet,
		},
		Route{
			"doNetTypeValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/nettype",
			s.doNetTypeValuesGet,
		},
		Route{
			"doRoleValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/role",
			s.doRoleValuesGet,
		},
		Route{
			"doSubRoleValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/subrole",
			s.doSubRoleValuesGet,
		},
		Route{
			"doStateValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/state",
			s.doStateValuesGet,
		},
		Route{
			"doTypeValuesGet",
			strings.ToUpper("Get"),
			s.valuesBase + "/type",
			s.doTypeValuesGet,
		},
		// Components
		Route{
			"doComponentGet",
			strings.ToUpper("Get"),
			s.componentsBase + "/{xname}",
			s.doComponentGet,
		},
		Route{
			"doComponentPut",
			strings.ToUpper("Put"),
			s.componentsBase + "/{xname}",
			s.doComponentPut,
		},
		Route{
			"doComponentDelete",
			strings.ToUpper("Delete"),
			s.componentsBase + "/{xname}",
			s.doComponentDelete,
		},
		Route{
			"doComponentsGet",
			strings.ToUpper("Get"),
			s.componentsBase,
			s.doComponentsGet,
		},
		Route{
			"doComponentsPost",
			strings.ToUpper("Post"),
			s.componentsBase,
			s.doComponentsPost,
		},
		Route{
			"doComponentsDeleteAll",
			strings.ToUpper("Delete"),
			s.componentsBase,
			s.doComponentsDeleteAll,
		},
		Route{
			"doCompBulkStateDataPatch",
			"PATCH",
			s.componentsBase + "/BulkStateData",
			s.doCompBulkStateDataPatch,
		},
		Route{
			"doCompStateDataPatch",
			"PATCH",
			s.componentsBase + "/{xname}/StateData",
			s.doCompStateDataPatch,
		},
		Route{
			"doCompBulkFlagOnlyPatch",
			"PATCH",
			s.componentsBase + "/BulkFlagOnly",
			s.doCompBulkFlagOnlyPatch,
		},
		Route{
			"doCompFlagOnlyPatch",
			"PATCH",
			s.componentsBase + "/{xname}/FlagOnly",
			s.doCompFlagOnlyPatch,
		},
		Route{
			"doCompBulkEnabledPatch",
			"PATCH",
			s.componentsBase + "/BulkEnabled",
			s.doCompBulkEnabledPatch,
		},
		Route{
			"doCompEnabled",
			"PATCH",
			s.componentsBase + "/{xname}/Enabled",
			s.doCompEnabledPatch,
		},
		Route{
			"doCompBulkSwStatusPatch",
			"PATCH",
			s.componentsBase + "/BulkSoftwareStatus",
			s.doCompBulkSwStatusPatch,
		},
		Route{
			"doCompSwStatus",
			"PATCH",
			s.componentsBase + "/{xname}/SoftwareStatus",
			s.doCompSwStatusPatch,
		},
		Route{
			"doCompBulkRolePatch",
			"PATCH",
			s.componentsBase + "/BulkRole",
			s.doCompBulkRolePatch,
		},
		Route{
			"doCompRole",
			"PATCH",
			s.componentsBase + "/{xname}/Role",
			s.doCompRolePatch,
		},
		Route{
			"doCompBulkNIDPatch",
			"PATCH",
			s.componentsBase + "/BulkNID",
			s.doCompBulkNIDPatch,
		},
		Route{
			"doCompNIDPatch",
			"PATCH",
			s.componentsBase + "/{xname}/NID",
			s.doCompNIDPatch,
		},
		Route{
			"doComponentByNIDGet",
			strings.ToUpper("Get"),
			s.componentsBase + "/ByNID/{nid}",
			s.doComponentByNIDGet,
		},
		Route{
			"doComponentByNIDQueryPost",
			strings.ToUpper("Post"),
			s.componentsBase + "/ByNID/Query",
			s.doComponentByNIDQueryPost,
		},
		Route{
			"doComponentsQueryPost",
			strings.ToUpper("Post"),
			s.componentsBase + "/Query",
			s.doComponentsQueryPost,
		},
		Route{
			"doComponentsQueryGet",
			strings.ToUpper("Get"),
			s.componentsBase + "/Query/{xname}",
			s.doComponentsQueryGet,
		},
		//
		//		Route{
		//			"DoComponentResourceURIByNIDGet",
		//			strings.ToUpper("Get"),
		//			"/stateManager/v1/State/Components/ByNID",
		//			DoComponentResourceURIByNIDGet,
		//		},
		//
		//		Route{
		//			"DoComponentResourceURIGet",
		//			strings.ToUpper("Get"),
		//			"/stateManager/v1/State/Components/",
		//			DoComponentResourceURIGet,
		//		},
		//
		//
		//		Route{
		//			"DoComponentEndpointResourceURIGet",
		//			strings.ToUpper("Get"),
		//			"/stateManager/v1/Inventory/RedfishEndpoints/ByComponent",
		//			DoComponentEndpointResourceURIGet,
		//		},
		//

		// ComponentEndpoints
		Route{
			"doComponentEndpointGet", // Individual entry
			strings.ToUpper("Get"),
			s.compEPBase + "/{xname}",
			s.doComponentEndpointGet,
		},
		Route{
			"doComponentEndpointDelete", // Individual entry
			strings.ToUpper("Delete"),
			s.compEPBase + "/{xname}",
			s.doComponentEndpointDelete,
		},
		Route{
			"doComponentEndpointsGet", // Whole collection
			strings.ToUpper("Get"),
			s.compEPBase,
			s.doComponentEndpointsGet,
		},
		Route{
			"doComponentEndpointsDeleteAll", // Whole collection
			strings.ToUpper("Delete"),
			s.compEPBase,
			s.doComponentEndpointsDeleteAll,
		},
		Route{
			"doComponentEndpointQueryGet",
			strings.ToUpper("Get"),
			s.compEPBase + "/{xname}",
			s.doComponentEndpointQueryGet,
		},
		//
		//		Route{
		//			"DoComponentEndpointQueryPost",
		//			strings.ToUpper("Post"),
		//			"/stateManager/v1/Inventory/RedfishEndpoints/ByComponent/Query",
		//			DoComponentEndpointQueryPost,
		//		},
		//

		// ServiceEndpoints
		Route{
			"doServiceEndpointGet", // Individual entry
			strings.ToUpper("Get"),
			s.serviceEPBase + "/{service}/RedfishEndpoints/{xname}",
			s.doServiceEndpointGet,
		},
		Route{
			"doServiceEndpointDelete", // Individual entry
			strings.ToUpper("Delete"),
			s.serviceEPBase + "/{service}/RedfishEndpoints/{xname}",
			s.doServiceEndpointDelete,
		},
		Route{
			"doServiceEndpointsGet", // Collection by service
			strings.ToUpper("Get"),
			s.serviceEPBase + "/{service}",
			s.doServiceEndpointsGet,
		},
		Route{
			"doServiceEndpointsGetAll", // Whole collection
			strings.ToUpper("Get"),
			s.serviceEPBase,
			s.doServiceEndpointsGetAll,
		},
		Route{
			"doServiceEndpointsDeleteAll", // Whole collection
			strings.ToUpper("Delete"),
			s.serviceEPBase,
			s.doServiceEndpointsDeleteAll,
		},

		// Component Ethernet Interfaces - V1
		Route{
			"doCompEthInterfacesGet",
			strings.ToUpper("Get"),
			s.compEthIntBase,
			s.doCompEthInterfacesGet,
		},
		Route{
			"doCompEthInterfacePost",
			strings.ToUpper("Post"),
			s.compEthIntBase,
			s.doCompEthInterfacePost,
		},
		Route{
			"doCompEthInterfaceDeleteAll",
			strings.ToUpper("Delete"),
			s.compEthIntBase,
			s.doCompEthInterfaceDeleteAll,
		},
		Route{
			"doCompEthInterfaceGet",
			strings.ToUpper("Get"),
			s.compEthIntBase + "/{id}",
			s.doCompEthInterfaceGet,
		},
		Route{
			"doCompEthInterfaceDelete",
			strings.ToUpper("Delete"),
			s.compEthIntBase + "/{id}",
			s.doCompEthInterfaceDelete,
		},
		Route{
			"doCompEthInterfacePatch",
			strings.ToUpper("Patch"),
			s.compEthIntBase + "/{id}",
			s.doCompEthInterfacePatch,
		},

		// NodeMaps
		Route{
			"doNodeMapGet",
			strings.ToUpper("Get"),
			s.nodeMapBase + "/{xname}",
			s.doNodeMapGet,
		},
		Route{
			"doNodeMapsGet",
			strings.ToUpper("Get"),
			s.nodeMapBase,
			s.doNodeMapsGet,
		},
		Route{
			"doNodeMapDelete",
			strings.ToUpper("Delete"),
			s.nodeMapBase + "/{xname}",
			s.doNodeMapDelete,
		},
		Route{
			"doNodeMapPut",
			strings.ToUpper("Put"),
			s.nodeMapBase + "/{xname}",
			s.doNodeMapPut,
		},
		Route{
			"doNodeMapsPost",
			strings.ToUpper("Post"),
			s.nodeMapBase,
			s.doNodeMapsPost,
		},
		Route{
			"doNodeMapsDeleteAll",
			strings.ToUpper("Delete"),
			s.nodeMapBase,
			s.doNodeMapsDeleteAll,
		},

		// Hardware Inventory History
		Route{
			"doHWInvHistByLocationGet",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "/History/{xname}",
			s.doHWInvHistByLocationGet,
		},
		Route{
			"doHWInvHistByLocationGetAll",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "/History",
			s.doHWInvHistByLocationGetAll,
		},
		Route{
			"doHWInvHistByFRUGet",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "ByFRU/History/{fruid}",
			s.doHWInvHistByFRUGet,
		},
		Route{
			"doHWInvHistByFRUGetAll",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "ByFRU/History",
			s.doHWInvHistByFRUGetAll,
		},
		Route{
			"doHWInvHistByLocationDelete",
			strings.ToUpper("Delete"),
			s.hwinvByLocBase + "/History/{xname}",
			s.doHWInvHistByLocationDelete,
		},
		Route{
			"doHWInvHistDeleteAll",
			strings.ToUpper("Delete"),
			s.hwinvByLocBase + "/History",
			s.doHWInvHistDeleteAll,
		},
		Route{
			"doHWInvHistByFRUDelete",
			strings.ToUpper("Delete"),
			s.hwinvByLocBase + "ByFRU/History/{fruid}",
			s.doHWInvHistByFRUDelete,
		},

		// Hardware Inventory
		Route{
			"doHWInvByLocationQueryGet",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "/Query/{xname}",
			s.doHWInvByLocationQueryGet,
		},
		Route{
			"doHWInvByFRUGet",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "ByFRU/{fruid}",
			s.doHWInvByFRUGet,
		},
		Route{
			"doHWInvByFRUGetAll",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "ByFRU",
			s.doHWInvByFRUGetAll,
		},
		Route{
			"doHWInvByLocationGet",
			strings.ToUpper("Get"),
			s.hwinvByLocBase + "/{xname}",
			s.doHWInvByLocationGet,
		},
		Route{
			"doHWInvByLocationGetAll",
			strings.ToUpper("Get"),
			s.hwinvByLocBase,
			s.doHWInvByLocationGetAll,
		},
		Route{
			"doHWInvByLocationPost",
			strings.ToUpper("Post"),
			s.hwinvByLocBase,
			s.doHWInvByLocationPost,
		},
		Route{
			"doHWInvByFRUDelete",
			strings.ToUpper("Delete"),
			s.hwinvByLocBase + "ByFRU/{fruid}",
			s.doHWInvByFRUDelete,
		},
		Route{
			"doHWInvByFRUDeleteAll",
			strings.ToUpper("Delete"),
			s.hwinvByLocBase + "ByFRU",
			s.doHWInvByFRUDeleteAll,
		},
		Route{
			"doHWInvByLocationDelete",
			strings.ToUpper("Delete"),
			s.hwinvByLocBase + "/{xname}",
			s.doHWInvByLocationDelete,
		},
		Route{
			"doHWInvByLocationDeleteAll",
			strings.ToUpper("Delete"),
			s.hwinvByLocBase,
			s.doHWInvByLocationDeleteAll,
		},

		//		Route{
		//			"DoRedfishEndpointResourceURIGet",
		//			strings.ToUpper("Get"),
		//			"/stateManager/v1/Inventory/RedfishEndpoints",
		//			DoRedfishEndpointResourceURIGet,
		//		},
		//

		// RefishEndpoints
		Route{
			"doRedfishEndpointGet",
			strings.ToUpper("Get"),
			s.redfishEPBase + "/{xname}",
			s.doRedfishEndpointGet,
		},
		Route{
			"doRedfishEndpointsGet",
			strings.ToUpper("Get"),
			s.redfishEPBase,
			s.doRedfishEndpointsGet,
		},
		Route{
			"doRedfishEndpointDelete",
			strings.ToUpper("Delete"),
			s.redfishEPBase + "/{xname}",
			s.doRedfishEndpointDelete,
		},
		Route{
			"doRedfishEndpointPut",
			strings.ToUpper("Put"),
			s.redfishEPBase + "/{xname}",
			s.doRedfishEndpointPut,
		},
		Route{
			"doRedfishEndpointPatch",
			strings.ToUpper("Patch"),
			s.redfishEPBase + "/{xname}",
			s.doRedfishEndpointPatch,
		},
		Route{
			"doRedfishEndpointsPost",
			strings.ToUpper("Post"),
			s.redfishEPBase,
			s.doRedfishEndpointsPost,
		},
		Route{
			"doRedfishEndpointsDeleteAll",
			strings.ToUpper("Delete"),
			s.redfishEPBase,
			s.doRedfishEndpointsDeleteAll,
		},
		Route{
			"doRedfishEndpointQueryGet",
			strings.ToUpper("Get"),
			s.redfishEPBase + "/Query/{xname}",
			s.doRedfishEndpointQueryGet,
		},
		//
		//		Route{
		//			"DoRedfishEndpointQueryPost",
		//			strings.ToUpper("Post"),
		//			"/stateManager/v1/Inventory/RedfishEndpoints/Query",
		//			DoRedfishEndpointQueryPost,
		//		},
		//
		Route{
			"doInventoryDiscoverPost",
			strings.ToUpper("Post"),
			s.invDiscoverBase,
			s.doInventoryDiscoverPost,
		},
		Route{
			"doDiscoveryStatusGetAll",
			strings.ToUpper("Get"),
			s.invDiscStatusBase,
			s.doDiscoveryStatusGetAll,
		},
		Route{
			"doDiscoveryStatusGet",
			strings.ToUpper("Get"),
			s.invDiscStatusBase + "/{id}",
			s.doDiscoveryStatusGet,
		},
		//		Route{
		//			"doInventoryDiscoverPost",
		//			strings.ToUpper("Post"),
		//			s.invDiscoverBase + "/{xname}",
		//			s.doInventoryDiscoverPost,
		//		},
		//		Route{
		//			"DoSIBQueryGet",
		//			strings.ToUpper("Get"),
		//			"/stateManager/v1/State/SIB/Query/{xname}",
		//			DoSIBQueryGet,
		//		},
		//
		//		Route{
		//			"DoSIBQueryPost",
		//			strings.ToUpper("Post"),
		//			"/stateManager/v1/State/SIB/Query",
		//			DoSIBQueryPost,
		//		},
		Route{
			"doGetSCNSubscription",
			strings.ToUpper("Get"),
			s.subscriptionBase + "/SCN",
			s.doGetSCNSubscriptionsAll,
		},
		Route{
			"doPostSCNSubscription",
			strings.ToUpper("Post"),
			s.subscriptionBase + "/SCN",
			s.doPostSCNSubscription,
		},
		Route{
			"doDeleteSCNSubscriptions",
			strings.ToUpper("Delete"),
			s.subscriptionBase + "/SCN",
			s.doDeleteSCNSubscriptionsAll,
		},
		Route{
			"doGetSCNSubscription",
			strings.ToUpper("Get"),
			s.subscriptionBase + "/SCN/{id}",
			s.doGetSCNSubscription,
		},
		Route{
			"doPutSCNSubscription",
			strings.ToUpper("Put"),
			s.subscriptionBase + "/SCN/{id}",
			s.doPutSCNSubscription,
		},
		Route{
			"doPatchSCNSubscription",
			strings.ToUpper("Patch"),
			s.subscriptionBase + "/SCN/{id}",
			s.doPatchSCNSubscription,
		},
		Route{
			"doDeleteSCNSubscription",
			strings.ToUpper("Delete"),
			s.subscriptionBase + "/SCN/{id}",
			s.doDeleteSCNSubscription,
		},

		// Groups
		Route{
			"doGroupsGet",
			strings.ToUpper("Get"),
			s.groupsBase,
			s.doGroupsGet,
		},
		Route{
			"doGroupsPost",
			strings.ToUpper("Post"),
			s.groupsBase,
			s.doGroupsPost,
		},
		Route{
			"doGroupLabelsGet",
			strings.ToUpper("Get"),
			s.groupsBase + "/labels",
			s.doGroupLabelsGet,
		},
		Route{
			"doGroupGet",
			strings.ToUpper("Get"),
			s.groupsBase + "/{group_label}",
			s.doGroupGet,
		},
		Route{
			"doGroupDelete",
			strings.ToUpper("Delete"),
			s.groupsBase + "/{group_label}",
			s.doGroupDelete,
		},
		Route{
			"doGroupPatch",
			strings.ToUpper("Patch"),
			s.groupsBase + "/{group_label}",
			s.doGroupPatch,
		},
		Route{
			"doGroupMembersGet",
			strings.ToUpper("Get"),
			s.groupsBase + "/{group_label}/members",
			s.doGroupMembersGet,
		},
		Route{
			"doGroupMembersPost",
			strings.ToUpper("Post"),
			s.groupsBase + "/{group_label}/members",
			s.doGroupMembersPost,
		},
		Route{
			"doGroupMemberDelete",
			strings.ToUpper("Delete"),
			s.groupsBase + "/{group_label}/members/{xname_id}",
			s.doGroupMemberDelete,
		},

		// Partitions
		Route{
			"doPartitionsGet",
			strings.ToUpper("Get"),
			s.partitionsBase,
			s.doPartitionsGet,
		},
		Route{
			"doPartitionsPost",
			strings.ToUpper("Post"),
			s.partitionsBase,
			s.doPartitionsPost,
		},
		Route{
			"doPartitionNamesGet",
			strings.ToUpper("Get"),
			s.partitionsBase + "/names",
			s.doPartitionNamesGet,
		},
		Route{
			"doPartitionGet",
			strings.ToUpper("Get"),
			s.partitionsBase + "/{partition_name}",
			s.doPartitionGet,
		},
		Route{
			"doPartitionDelete",
			strings.ToUpper("Delete"),
			s.partitionsBase + "/{partition_name}",
			s.doPartitionDelete,
		},
		Route{
			"doPartitionPatch",
			strings.ToUpper("Patch"),
			s.partitionsBase + "/{partition_name}",
			s.doPartitionPatch,
		},
		Route{
			"doPartitionMembersGet",
			strings.ToUpper("Get"),
			s.partitionsBase + "/{partition_name}/members",
			s.doPartitionMembersGet,
		},
		Route{
			"doPartitionMembersPost",
			strings.ToUpper("Post"),
			s.partitionsBase + "/{partition_name}/members",
			s.doPartitionMembersPost,
		},
		Route{
			"doPartitionMemberDelete",
			strings.ToUpper("Delete"),
			s.partitionsBase + "/{partition_name}/members/{xname_id}",
			s.doPartitionMemberDelete,
		},

		// Memberships
		Route{
			"doMembershipsGet",
			strings.ToUpper("Get"),
			s.membershipsBase,
			s.doMembershipsGet,
		},
		Route{
			"doMembershipGet",
			strings.ToUpper("Get"),
			s.membershipsBase + "/{xname}",
			s.doMembershipGet,
		},

		// Component Locks
		Route{
			"doCompLocksGet",
			strings.ToUpper("Get"),
			s.compLockBase,
			s.doCompLocksGet,
		},
		Route{
			"doCompLocksPost",
			strings.ToUpper("Post"),
			s.compLockBase,
			s.doCompLocksPost,
		},
		Route{
			"doCompLockGet",
			strings.ToUpper("Get"),
			s.compLockBase + "/{lock_id}",
			s.doCompLockGet,
		},
		Route{
			"doCompLockDelete",
			strings.ToUpper("Delete"),
			s.compLockBase + "/{lock_id}",
			s.doCompLockDelete,
		},
		Route{
			"doCompLockPatch",
			strings.ToUpper("Patch"),
			s.compLockBase + "/{lock_id}",
			s.doCompLockPatch,
		},

		// PowerMaps
		Route{
			"doPowerMapGet",
			strings.ToUpper("Get"),
			s.powerMapBase + "/{xname}",
			s.doPowerMapGet,
		},
		Route{
			"doPowerMapsGet",
			strings.ToUpper("Get"),
			s.powerMapBase,
			s.doPowerMapsGet,
		},
		Route{
			"doPowerMapDelete",
			strings.ToUpper("Delete"),
			s.powerMapBase + "/{xname}",
			s.doPowerMapDelete,
		},
		Route{
			"doPowerMapPut",
			strings.ToUpper("Put"),
			s.powerMapBase + "/{xname}",
			s.doPowerMapPut,
		},
		Route{
			"doPowerMapsPost",
			strings.ToUpper("Post"),
			s.powerMapBase,
			s.doPowerMapsPost,
		},
		Route{
			"doPowerMapsDeleteAll",
			strings.ToUpper("Delete"),
			s.powerMapBase,
			s.doPowerMapsDeleteAll,
		},

		///////////////////////////////////////////////////////////////////////
		// v2 API routes
		///////////////////////////////////////////////////////////////////////

		// HSM Service State
		Route{
			"doReadyGetV2",
			strings.ToUpper("Get"),
			s.serviceBaseV2 + "/ready",
			s.doReadyGet,
		},
		Route{
			"doLivenessGetV2",
			strings.ToUpper("Get"),
			s.serviceBaseV2 + "/liveness",
			s.doLivenessGet,
		},
		Route{
			"doValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2,
			s.doValuesGet,
		},
		Route{
			"doArchValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/arch",
			s.doArchValuesGet,
		},
		Route{
			"doClassValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/class",
			s.doClassValuesGet,
		},
		Route{
			"doFlagValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/flag",
			s.doFlagValuesGet,
		},
		Route{
			"doNetTypeValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/nettype",
			s.doNetTypeValuesGet,
		},
		Route{
			"doRoleValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/role",
			s.doRoleValuesGet,
		},
		Route{
			"doSubRoleValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/subrole",
			s.doSubRoleValuesGet,
		},
		Route{
			"doStateValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/state",
			s.doStateValuesGet,
		},
		Route{
			"doTypeValuesGetV2",
			strings.ToUpper("Get"),
			s.valuesBaseV2 + "/type",
			s.doTypeValuesGet,
		},
		// Components
		Route{
			"doComponentGetV2",
			strings.ToUpper("Get"),
			s.componentsBaseV2 + "/{xname}",
			s.doComponentGet,
		},
		Route{
			"doComponentPutV2",
			strings.ToUpper("Put"),
			s.componentsBaseV2 + "/{xname}",
			s.doComponentPut,
		},
		Route{
			"doComponentDeleteV2",
			strings.ToUpper("Delete"),
			s.componentsBaseV2 + "/{xname}",
			s.doComponentDelete,
		},
		Route{
			"doComponentsGetV2",
			strings.ToUpper("Get"),
			s.componentsBaseV2,
			s.doComponentsGet,
		},
		Route{
			"doComponentsPostV2",
			strings.ToUpper("Post"),
			s.componentsBaseV2,
			s.doComponentsPost,
		},
		Route{
			"doComponentsDeleteAllV2",
			strings.ToUpper("Delete"),
			s.componentsBaseV2,
			s.doComponentsDeleteAll,
		},
		Route{
			"doCompBulkStateDataPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/BulkStateData",
			s.doCompBulkStateDataPatch,
		},
		Route{
			"doCompStateDataPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/{xname}/StateData",
			s.doCompStateDataPatch,
		},
		Route{
			"doCompBulkFlagOnlyPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/BulkFlagOnly",
			s.doCompBulkFlagOnlyPatch,
		},
		Route{
			"doCompFlagOnlyPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/{xname}/FlagOnly",
			s.doCompFlagOnlyPatch,
		},
		Route{
			"doCompBulkEnabledPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/BulkEnabled",
			s.doCompBulkEnabledPatch,
		},
		Route{
			"doCompEnabledV2",
			"PATCH",
			s.componentsBaseV2 + "/{xname}/Enabled",
			s.doCompEnabledPatch,
		},
		Route{
			"doCompBulkSwStatusPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/BulkSoftwareStatus",
			s.doCompBulkSwStatusPatch,
		},
		Route{
			"doCompSwStatusV2",
			"PATCH",
			s.componentsBaseV2 + "/{xname}/SoftwareStatus",
			s.doCompSwStatusPatch,
		},
		Route{
			"doCompBulkRolePatchV2",
			"PATCH",
			s.componentsBaseV2 + "/BulkRole",
			s.doCompBulkRolePatch,
		},
		Route{
			"doCompRoleV2",
			"PATCH",
			s.componentsBaseV2 + "/{xname}/Role",
			s.doCompRolePatch,
		},
		Route{
			"doCompBulkNIDPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/BulkNID",
			s.doCompBulkNIDPatch,
		},
		Route{
			"doCompNIDPatchV2",
			"PATCH",
			s.componentsBaseV2 + "/{xname}/NID",
			s.doCompNIDPatch,
		},
		Route{
			"doComponentByNIDGetV2",
			strings.ToUpper("Get"),
			s.componentsBaseV2 + "/ByNID/{nid}",
			s.doComponentByNIDGet,
		},
		Route{
			"doComponentByNIDQueryPostV2",
			strings.ToUpper("Post"),
			s.componentsBaseV2 + "/ByNID/Query",
			s.doComponentByNIDQueryPost,
		},
		Route{
			"doComponentsQueryPostV2",
			strings.ToUpper("Post"),
			s.componentsBaseV2 + "/Query",
			s.doComponentsQueryPost,
		},
		Route{
			"doComponentsQueryGetV2",
			strings.ToUpper("Get"),
			s.componentsBaseV2 + "/Query/{xname}",
			s.doComponentsQueryGet,
		},

		// ComponentEndpoints
		Route{
			"doComponentEndpointGetV2", // Individual entry
			strings.ToUpper("Get"),
			s.compEPBaseV2 + "/{xname}",
			s.doComponentEndpointGet,
		},
		Route{
			"doComponentEndpointDeleteV2", // Individual entry
			strings.ToUpper("Delete"),
			s.compEPBaseV2 + "/{xname}",
			s.doComponentEndpointDelete,
		},
		Route{
			"doComponentEndpointsGetV2", // Whole collection
			strings.ToUpper("Get"),
			s.compEPBaseV2,
			s.doComponentEndpointsGet,
		},
		Route{
			"doComponentEndpointsDeleteAllV2", // Whole collection
			strings.ToUpper("Delete"),
			s.compEPBaseV2,
			s.doComponentEndpointsDeleteAll,
		},
		Route{
			"doComponentEndpointQueryGetV2",
			strings.ToUpper("Get"),
			s.compEPBaseV2 + "/{xname}",
			s.doComponentEndpointQueryGet,
		},

		// ServiceEndpoints
		Route{
			"doServiceEndpointGetV2", // Individual entry
			strings.ToUpper("Get"),
			s.serviceEPBaseV2 + "/{service}/RedfishEndpoints/{xname}",
			s.doServiceEndpointGet,
		},
		Route{
			"doServiceEndpointDeleteV2", // Individual entry
			strings.ToUpper("Delete"),
			s.serviceEPBaseV2 + "/{service}/RedfishEndpoints/{xname}",
			s.doServiceEndpointDelete,
		},
		Route{
			"doServiceEndpointsGetV2", // Collection by service
			strings.ToUpper("Get"),
			s.serviceEPBaseV2 + "/{service}",
			s.doServiceEndpointsGet,
		},
		Route{
			"doServiceEndpointsGetAllV2", // Whole collection
			strings.ToUpper("Get"),
			s.serviceEPBaseV2,
			s.doServiceEndpointsGetAll,
		},
		Route{
			"doServiceEndpointsDeleteAllV2", // Whole collection
			strings.ToUpper("Delete"),
			s.serviceEPBaseV2,
			s.doServiceEndpointsDeleteAll,
		},

		// Component Ethernet Interfaces - V2
		Route{
			"doCompEthInterfacesGetV2",
			strings.ToUpper("Get"),
			s.compEthIntBaseV2,
			s.doCompEthInterfacesGetV2,
		},
		Route{
			"doCompEthInterfacePostV2",
			strings.ToUpper("Post"),
			s.compEthIntBaseV2,
			s.doCompEthInterfacePostV2,
		},
		Route{
			"doCompEthInterfaceDeleteAllV2",
			strings.ToUpper("Delete"),
			s.compEthIntBaseV2,
			s.doCompEthInterfaceDeleteAll, // V1/V2 APIS share the same handler
		},
		Route{
			"doCompEthInterfaceGetV2",
			strings.ToUpper("Get"),
			s.compEthIntBaseV2 + "/{id}",
			s.doCompEthInterfaceGetV2,
		},
		Route{
			"doCompEthInterfaceDeleteV2",
			strings.ToUpper("Delete"),
			s.compEthIntBaseV2 + "/{id}",
			s.doCompEthInterfaceDelete, // V1/V2 APIS share the same handler
		},
		Route{
			"doCompEthInterfacePatchV2",
			strings.ToUpper("Patch"),
			s.compEthIntBaseV2 + "/{id}",
			s.doCompEthInterfacePatchV2,
		},
		Route{
			"doCompEthInterfaceIPAddressesGetV2",
			strings.ToUpper("Get"),
			s.compEthIntBaseV2 + "/{id}/IPAddresses",
			s.doCompEthInterfaceIPAddressesGetV2,
		},
		Route{
			"doCompEthInterfaceIPAddressPostV2",
			strings.ToUpper("Post"),
			s.compEthIntBaseV2 + "/{id}/IPAddresses",
			s.doCompEthInterfaceIPAddressPostV2,
		},
		Route{
			"doCompEthInterfaceMemberPatchV2",
			strings.ToUpper("Patch"),
			s.compEthIntBaseV2 + "/{id}/IPAddresses/{ipaddr}",
			s.doCompEthInterfaceIPAddressPatchV2,
		},
		Route{
			"doCompEthInterfaceMemberDeleteV2",
			strings.ToUpper("Delete"),
			s.compEthIntBaseV2 + "/{id}/IPAddresses/{ipaddr}",
			s.doCompEthInterfaceIPAddressDeleteV2,
		},

		// NodeMaps
		Route{
			"doNodeMapGetV2",
			strings.ToUpper("Get"),
			s.nodeMapBaseV2 + "/{xname}",
			s.doNodeMapGet,
		},
		Route{
			"doNodeMapsGetV2",
			strings.ToUpper("Get"),
			s.nodeMapBaseV2,
			s.doNodeMapsGet,
		},
		Route{
			"doNodeMapDeleteV2",
			strings.ToUpper("Delete"),
			s.nodeMapBaseV2 + "/{xname}",
			s.doNodeMapDelete,
		},
		Route{
			"doNodeMapPutV2",
			strings.ToUpper("Put"),
			s.nodeMapBaseV2 + "/{xname}",
			s.doNodeMapPut,
		},
		Route{
			"doNodeMapsPostV2",
			strings.ToUpper("Post"),
			s.nodeMapBaseV2,
			s.doNodeMapsPost,
		},
		Route{
			"doNodeMapsDeleteAllV2",
			strings.ToUpper("Delete"),
			s.nodeMapBaseV2,
			s.doNodeMapsDeleteAll,
		},

		// Hardware Inventory History
		Route{
			"doHWInvHistByLocationGetV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "/History/{xname}",
			s.doHWInvHistByLocationGet,
		},
		Route{
			"doHWInvHistByLocationGetAllV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "/History",
			s.doHWInvHistByLocationGetAll,
		},
		Route{
			"doHWInvHistByFRUGetV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "ByFRU/History/{fruid}",
			s.doHWInvHistByFRUGet,
		},
		Route{
			"doHWInvHistByFRUGetAllV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "ByFRU/History",
			s.doHWInvHistByFRUGetAll,
		},
		Route{
			"doHWInvHistByLocationDeleteV2",
			strings.ToUpper("Delete"),
			s.hwinvByLocBaseV2 + "/History/{xname}",
			s.doHWInvHistByLocationDelete,
		},
		Route{
			"doHWInvHistDeleteAllV2",
			strings.ToUpper("Delete"),
			s.hwinvByLocBaseV2 + "/History",
			s.doHWInvHistDeleteAll,
		},
		Route{
			"doHWInvHistByFRUDeleteV2",
			strings.ToUpper("Delete"),
			s.hwinvByLocBaseV2 + "ByFRU/History/{fruid}",
			s.doHWInvHistByFRUDelete,
		},

		// Hardware Inventory
		Route{
			"doHWInvByLocationQueryGetV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "/Query/{xname}",
			s.doHWInvByLocationQueryGet,
		},
		Route{
			"doHWInvByFRUGetV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "ByFRU/{fruid}",
			s.doHWInvByFRUGet,
		},
		Route{
			"doHWInvByFRUGetAllV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "ByFRU",
			s.doHWInvByFRUGetAll,
		},
		Route{
			"doHWInvByLocationGetV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2 + "/{xname}",
			s.doHWInvByLocationGet,
		},

		Route{
			"doHWInvByLocationPostV2",
			strings.ToUpper("Post"),
			s.hwinvByLocBaseV2,
			s.doHWInvByLocationPost,
		},

		Route{
			"doHWInvByLocationGetAllV2",
			strings.ToUpper("Get"),
			s.hwinvByLocBaseV2,
			s.doHWInvByLocationGetAll,
		},
		Route{
			"doHWInvByFRUDeleteV2",
			strings.ToUpper("Delete"),
			s.hwinvByLocBaseV2 + "ByFRU/{fruid}",
			s.doHWInvByFRUDelete,
		},
		Route{
			"doHWInvByFRUDeleteAllV2",
			strings.ToUpper("Delete"),
			s.hwinvByLocBaseV2 + "ByFRU",
			s.doHWInvByFRUDeleteAll,
		},
		Route{
			"doHWInvByLocationDeleteV2",
			strings.ToUpper("Delete"),
			s.hwinvByLocBaseV2 + "/{xname}",
			s.doHWInvByLocationDelete,
		},
		Route{
			"doHWInvByLocationDeleteAllV2",
			strings.ToUpper("Delete"),
			s.hwinvByLocBaseV2,
			s.doHWInvByLocationDeleteAll,
		},

		// RefishEndpoints
		Route{
			"doRedfishEndpointGetV2",
			strings.ToUpper("Get"),
			s.redfishEPBaseV2 + "/{xname}",
			s.doRedfishEndpointGet,
		},
		Route{
			"doRedfishEndpointsGetV2",
			strings.ToUpper("Get"),
			s.redfishEPBaseV2,
			s.doRedfishEndpointsGet,
		},
		Route{
			"doRedfishEndpointDeleteV2",
			strings.ToUpper("Delete"),
			s.redfishEPBaseV2 + "/{xname}",
			s.doRedfishEndpointDelete,
		},
		Route{
			"doRedfishEndpointPutV2",
			strings.ToUpper("Put"),
			s.redfishEPBaseV2 + "/{xname}",
			s.doRedfishEndpointPut,
		},
		Route{
			"doRedfishEndpointPatchV2",
			strings.ToUpper("Patch"),
			s.redfishEPBaseV2 + "/{xname}",
			s.doRedfishEndpointPatch,
		},
		Route{
			"doRedfishEndpointsPostV2",
			strings.ToUpper("Post"),
			s.redfishEPBaseV2,
			s.doRedfishEndpointsPost,
		},
		Route{
			"doRedfishEndpointsDeleteAllV2",
			strings.ToUpper("Delete"),
			s.redfishEPBaseV2,
			s.doRedfishEndpointsDeleteAll,
		},
		Route{
			"doRedfishEndpointQueryGetV2",
			strings.ToUpper("Get"),
			s.redfishEPBaseV2 + "/Query/{xname}",
			s.doRedfishEndpointQueryGet,
		},
		Route{
			"doInventoryDiscoverPostV2",
			strings.ToUpper("Post"),
			s.invDiscoverBaseV2,
			s.doInventoryDiscoverPost,
		},
		Route{
			"doDiscoveryStatusGetAllV2",
			strings.ToUpper("Get"),
			s.invDiscStatusBaseV2,
			s.doDiscoveryStatusGetAll,
		},
		Route{
			"doDiscoveryStatusGetV2",
			strings.ToUpper("Get"),
			s.invDiscStatusBaseV2 + "/{id}",
			s.doDiscoveryStatusGet,
		},

		Route{
			"doGetSCNSubscriptionV2",
			strings.ToUpper("Get"),
			s.subscriptionBaseV2 + "/SCN",
			s.doGetSCNSubscriptionsAll,
		},
		Route{
			"doPostSCNSubscriptionV2",
			strings.ToUpper("Post"),
			s.subscriptionBaseV2 + "/SCN",
			s.doPostSCNSubscription,
		},
		Route{
			"doDeleteSCNSubscriptionsV2",
			strings.ToUpper("Delete"),
			s.subscriptionBaseV2 + "/SCN",
			s.doDeleteSCNSubscriptionsAll,
		},
		Route{
			"doGetSCNSubscriptionV2",
			strings.ToUpper("Get"),
			s.subscriptionBaseV2 + "/SCN/{id}",
			s.doGetSCNSubscription,
		},
		Route{
			"doPutSCNSubscriptionV2",
			strings.ToUpper("Put"),
			s.subscriptionBaseV2 + "/SCN/{id}",
			s.doPutSCNSubscription,
		},
		Route{
			"doPatchSCNSubscriptionV2",
			strings.ToUpper("Patch"),
			s.subscriptionBaseV2 + "/SCN/{id}",
			s.doPatchSCNSubscription,
		},
		Route{
			"doDeleteSCNSubscriptionV2",
			strings.ToUpper("Delete"),
			s.subscriptionBaseV2 + "/SCN/{id}",
			s.doDeleteSCNSubscription,
		},

		// Groups
		Route{
			"doGroupsGetV2",
			strings.ToUpper("Get"),
			s.groupsBaseV2,
			s.doGroupsGet,
		},
		Route{
			"doGroupsPostV2",
			strings.ToUpper("Post"),
			s.groupsBaseV2,
			s.doGroupsPost,
		},
		Route{
			"doGroupLabelsGetV2",
			strings.ToUpper("Get"),
			s.groupsBaseV2 + "/labels",
			s.doGroupLabelsGet,
		},
		Route{
			"doGroupGetV2",
			strings.ToUpper("Get"),
			s.groupsBaseV2 + "/{group_label}",
			s.doGroupGet,
		},
		Route{
			"doGroupDeleteV2",
			strings.ToUpper("Delete"),
			s.groupsBaseV2 + "/{group_label}",
			s.doGroupDelete,
		},
		Route{
			"doGroupPatchV2",
			strings.ToUpper("Patch"),
			s.groupsBaseV2 + "/{group_label}",
			s.doGroupPatch,
		},
		Route{
			"doGroupMembersGetV2",
			strings.ToUpper("Get"),
			s.groupsBaseV2 + "/{group_label}/members",
			s.doGroupMembersGet,
		},
		Route{
			"doGroupMembersPostV2",
			strings.ToUpper("Post"),
			s.groupsBaseV2 + "/{group_label}/members",
			s.doGroupMembersPost,
		},
		Route{
			"doGroupMemberDeleteV2",
			strings.ToUpper("Delete"),
			s.groupsBaseV2 + "/{group_label}/members/{xname_id}",
			s.doGroupMemberDelete,
		},

		// Partitions
		Route{
			"doPartitionsGetV2",
			strings.ToUpper("Get"),
			s.partitionsBaseV2,
			s.doPartitionsGet,
		},
		Route{
			"doPartitionsPostV2",
			strings.ToUpper("Post"),
			s.partitionsBaseV2,
			s.doPartitionsPost,
		},
		Route{
			"doPartitionNamesGetV2",
			strings.ToUpper("Get"),
			s.partitionsBaseV2 + "/names",
			s.doPartitionNamesGet,
		},
		Route{
			"doPartitionGetV2",
			strings.ToUpper("Get"),
			s.partitionsBaseV2 + "/{partition_name}",
			s.doPartitionGet,
		},
		Route{
			"doPartitionDeleteV2",
			strings.ToUpper("Delete"),
			s.partitionsBaseV2 + "/{partition_name}",
			s.doPartitionDelete,
		},
		Route{
			"doPartitionPatchV2",
			strings.ToUpper("Patch"),
			s.partitionsBaseV2 + "/{partition_name}",
			s.doPartitionPatch,
		},
		Route{
			"doPartitionMembersGetV2",
			strings.ToUpper("Get"),
			s.partitionsBaseV2 + "/{partition_name}/members",
			s.doPartitionMembersGet,
		},
		Route{
			"doPartitionMembersPostV2",
			strings.ToUpper("Post"),
			s.partitionsBaseV2 + "/{partition_name}/members",
			s.doPartitionMembersPost,
		},
		Route{
			"doPartitionMemberDeleteV2",
			strings.ToUpper("Delete"),
			s.partitionsBaseV2 + "/{partition_name}/members/{xname_id}",
			s.doPartitionMemberDelete,
		},

		// Memberships
		Route{
			"doMembershipsGetV2",
			strings.ToUpper("Get"),
			s.membershipsBaseV2,
			s.doMembershipsGet,
		},
		Route{
			"doMembershipGetV2",
			strings.ToUpper("Get"),
			s.membershipsBaseV2 + "/{xname}",
			s.doMembershipGet,
		},

		// V2 Component Locks

		//Admin reservedMap
		Route{
			"doCompLocksReservationRemoveV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/reservations/remove",
			s.doCompLocksReservationRemove,
		},
		Route{
			"doCompLocksReservationReleaseV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/reservations/release",
			s.doCompLocksReservationRelease,
		},
		Route{
			"doCompLocksReservationCreateV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/reservations",
			s.doCompLocksReservationCreate,
		},

		//Service reservedMap
		Route{
			"doCompLocksServiceReservationRenewV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/service/reservations/renew",
			s.doCompLocksServiceReservationRenew,
		},
		Route{
			"doCompLocksServiceReservationReleaseV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/service/reservations/release",
			s.doCompLocksServiceReservationRelease,
		},
		Route{
			"doCompLocksServiceReservationCreateV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/service/reservations",
			s.doCompLocksServiceReservationCreate,
		},
		Route{
			"doCompLocksServiceReservationCheckV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/service/reservations/check",
			s.doCompLocksServiceReservationCheck,
		},

		//Admin Locks
		Route{
			"doCompLocksStatusV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/status",
			s.doCompLocksStatus,
		},
		Route{
			"doCompLocksLockV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/lock",
			s.doCompLocksLock,
		},
		Route{
			"doCompLocksUnlockV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/unlock",
			s.doCompLocksUnlock,
		},
		Route{
			"doCompLocksRepairV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/repair",
			s.doCompLocksRepair,
		},
		Route{
			"doCompLocksDisableV2",
			strings.ToUpper("Post"),
			s.compLockBaseV2 + "/disable",
			s.doCompLocksDisable,
		},

		// PowerMaps
		Route{
			"doPowerMapGetV2",
			strings.ToUpper("Get"),
			s.powerMapBaseV2 + "/{xname}",
			s.doPowerMapGet,
		},
		Route{
			"doPowerMapsGetV2",
			strings.ToUpper("Get"),
			s.powerMapBaseV2,
			s.doPowerMapsGet,
		},
		Route{
			"doPowerMapDeleteV2",
			strings.ToUpper("Delete"),
			s.powerMapBaseV2 + "/{xname}",
			s.doPowerMapDelete,
		},
		Route{
			"doPowerMapPutV2",
			strings.ToUpper("Put"),
			s.powerMapBaseV2 + "/{xname}",
			s.doPowerMapPut,
		},
		Route{
			"doPowerMapsPostV2",
			strings.ToUpper("Post"),
			s.powerMapBaseV2,
			s.doPowerMapsPost,
		},
		Route{
			"doPowerMapsDeleteAllV2",
			strings.ToUpper("Delete"),
			s.powerMapBaseV2,
			s.doPowerMapsDeleteAll,
		},
	}
}
