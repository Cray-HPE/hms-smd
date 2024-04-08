// MIT License
//
// (C) Copyright [2024] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package rf

import (
	"encoding/json"
	"fmt"
	"strings"
)

const FOXCONN_NODE_ETH_SUFFIX = "-node_eth"

///////////////////////////////////////////////////////////////////////////////
//
//
// localhost:~ # curl -ks -u root:${BMC_CREDS} https://10.5.1.106/redfish/v1/Systems/system | jq .
// {
//     "@odata.context": "",
//     "@odata.id": "/redfish/v1/Systems/system",
//     "@odata.type": "#ComputerSystem.v1_18_0.ComputerSystem",
//     ... <REDACTED> ...
//     "Oem": {
//         "InsydeNcsi": {
//             "Ncsi": {
//                 "@odata.id": "/redfish/v1/Systems/system/Oem/Insyde/Ncsi"
//             }
//         }
//     }
//     ... <REDACTED> ...
// }
//
///////////////////////////////////////////////////////////////////////////////

type ComputerSystemOemInsyde struct {
	Ncsi ResourceID `json:"Ncsi"`
}

///////////////////////////////////////////////////////////////////////////////
//
// localhost:~ # curl -ks -u root:${BMC_CREDS} https://10.5.1.106/redfish/v1/Systems/system/Oem/Insyde/Ncsi | jq .
// {
//   "@odata.id": "/redfish/v1/Systems/system/Oem/Insyde/Ncsi",
//   "@odata.type": "#InsydeNcsiCollection.InsydeNcsiCollection",
//   "Description": "The NetworkAdapterCollection schema describes a collection of network adapter instances.",
//   "Members": [
//     {
//       "@odata.id": "/redfish/v1/Systems/system/Oem/Insyde/Ncsi/1"
//     },
//     {
//       "@odata.id": "/redfish/v1/Systems/system/Oem/Insyde/Ncsi/2"
//     }
//   ],
//   "Members@odata.count": 2,
//   "Name": "Insyde Ncsi Collection"
// }
//
///////////////////////////////////////////////////////////////////////////////

type InsydeOemNcsiCollection GenericCollection

///////////////////////////////////////////////////////////////////////////////
//
// localhost:~ # curl -ks -u root:${BMC_CREDS} https://10.5.1.106/redfish/v1/Systems/system/Oem/Insyde/Ncsi/2 | jq .
// {
//   "@odata.id": "/redfish/v1/Systems/system/Oem/Insyde/Ncsi/2",
//   "@odata.type": "#InsydeNcsi.v1_0_0.InsydeNcsi",
//   "Description": "The InsydeNcsi schema contains properties related to NCSI device.",
//   "DeviceType": "NCSIOverRBT",
//   "Id": "2",
//   "Name": "2",
//   "Package": [
//     {
//       "@odata.id": "/redfish/v1/Systems/system/Oem/Insyde/Ncsi/2/Package/1"
//     }
//   ],
//   "VersionID": {
//     "FirmwareName": "X550 FW Ver ",
//     "FirmwareVersion": "00.00.03.60",
//     "ManufacturerID": "0x57010000",
//     "NcsiVersion": "1.0.1",
//     "PCIDID": "0x6315",
//     "PCISSID": "0x0000",
//     "PCIVID": "0x0000"
//   }
// }
//
///////////////////////////////////////////////////////////////////////////////

type InsydeOemNcsiMember struct {
	DeviceType              string                  `json: DeviceType`
	Id                      string                  `json: Id`
	Package                 []ResourceID  			`json:"Package"`
	VersionId               InsydeOemNcsiVersionId  `json:"VersionID"`
}

//type InsydeOemNcsiPackage struct {
//	Oid                     string          		`json:"@odata.id"`
//}

type InsydeOemNcsiVersionId struct {
	FirmwareName            string                  `json:"FirmwareName"`
}

///////////////////////////////////////////////////////////////////////////////
//
// The following is heavily redacted due to size.  Only relevant fields are shown.
//
// localhost:~ # curl -ks -u root:${BMC_CREDS} https://10.5.1.106/redfish/v1/Systems/system/Oem/Insyde/Ncsi/2/Package/1 | jq .
// {
//   "@odata.id": "/redfish/v1/Systems/system/Oem/Insyde/Ncsi/2/Package/1",
//   "@odata.type": "#InsydeNcsiPackage.v1_0_0.InsydeNcsiPackage",
//   "Description": "The InsydeNcsiPackage schema contains properties related to NcsiPackage.",
//   "Id": "1",
//   "Name": "1",
//   "PackageInfo": [
//     {
//       "ChannelIndex": 1,
//       "MACAddress": "04:D9:C8:5D:55:05",
//     },
//     {
//       "ChannelIndex": 2,
//       }
//     }
//   ]
// }
//
///////////////////////////////////////////////////////////////////////////////


type InsydeOemPackage struct {
	Id                      string                  `json: Id`
	PackageInfo             []InsydeOemPackageInfo  `json:"PackageInfo"`
}

type InsydeOemPackageInfo struct {
	ChannelIndex            int                     `json:"ChannelIndex"`
	MACAddress              string                  `json:"MACAddress"`
}

// Parses redfish to find ethernet interfaces for Foxconn Paradise
func discoverFoxconnENetInterfaces(s *EpSystem) {
	//////////////////////////////////////////////////////
 	// Parse /redfish/v1/Systems/system/Oem/Insyde/Ncsi

	path := s.SystemRF.OEM.InsydeNcsi.Ncsi.Oid

	url := s.epRF.FQDN + path
	jsonData, err := s.epRF.GETRelative(path)
	if err != nil || jsonData == nil {
		s.LastStatus = HTTPsGetFailed
		return
	}
	s.LastStatus = HTTPsGetOk

	var n InsydeOemNcsiCollection
	if err := json.Unmarshal(jsonData, &n); err != nil {
		errlog.Printf("Failed to decode %s: %s\n", url, err)
		s.LastStatus = EPResponseFailedDecode
	}

	if n.MembersOCount < 1 {
		errlog.Printf("No Ncsi members detected")
		return
	}

	//////////////////////////////////////////////////////
 	// Parse each /redfish/v1/Systems/system/Oem/Insyde/Ncsi/#

	for _, ncsiMember := range n.Members {
		path := ncsiMember.Oid

		url := s.epRF.FQDN + path
		jsonData, err = s.epRF.GETRelative(path)
		if err != nil || jsonData == nil {
			s.LastStatus = HTTPsGetFailed
			return
		}
		s.LastStatus = HTTPsGetOk

		var nm InsydeOemNcsiMember
		if err := json.Unmarshal(jsonData, &nm); err != nil {
			errlog.Printf("Failed to decode %s: %s\n", url, err)
			s.LastStatus = EPResponseFailedDecode
		}

		errlog.Printf("<========== JW_DEBUG ==========> discoverFoxconnENetInterfaces: FirmwareName=%s\n", nm.VersionId.FirmwareName)

		//////////////////////////////////////////////////////
 		// Parse /redfish/v1/Systems/system/Oem/Insyde/Ncsi/#/Package/#

		path = nm.Package[0].Oid

		url = s.epRF.FQDN + path
		jsonData, err = s.epRF.GETRelative(path)
		if err != nil || jsonData == nil {
			s.LastStatus = HTTPsGetFailed
			return
		}
		s.LastStatus = HTTPsGetOk

		var p InsydeOemPackage
		if err := json.Unmarshal(jsonData, &p); err != nil {
			errlog.Printf("Failed to decode %s: %s\n", url, err)
			s.LastStatus = EPResponseFailedDecode
		}

		//////////////////////////////////////////////////////
 		// Parse /redfish/v1/Systems/system/Oem/Insyde/Ncsi/#/Package/#.PackageInfo[]
		//
		// Some controllers have multiple MACs but the host ethernet controller will have only one
		// so stop parsing after the first MAC address is found.

		for j, pi := range p.PackageInfo {
			errlog.Printf("<========== JW_DEBUG ==========> discoverFoxconnENetInterfaces: channel=%d\n", pi.ChannelIndex)
			if pi.MACAddress != "" {
				s.ENetInterfaces.Num++

				ei := NewEpEthInterface(s.epRF, ncsiMember.Oid, s.RedfishSubtype, nm.Package[0], s.ENetInterfaces.Num)

				ei.EtherIfaceRF.Oid = path
				ei.EtherIfaceRF.MACAddress = pi.MACAddress
				ei.EtherIfaceRF.Description = "Auto-detected Foxconn NCSI Ethernet Interface"

				// ID = "foxconn-ncsi-" + ncsi number + "-" + package number + "-" + channel index
				ei.EtherIfaceRF.Id = "foxconn-ncsi-" + nm.Id + "-" + p.Id + "-" + fmt.Sprint(j)

				// This is the only (hopefully) unique identifier for the onboard host ethernet
				if strings.TrimSpace(nm.VersionId.FirmwareName) == "X550 FW Ver" {
					// We append a "-node_eth" string to the end of the Description so that we can
					// identify it later.
					ei.EtherIfaceRF.Id += FOXCONN_NODE_ETH_SUFFIX
					errlog.Printf("<========== JW_DEBUG ==========> discoverFoxconnENetInterfaces: setting system MACAddr=%s\n", ei.MACAddr)
				}

				ei.BaseOdataID = ei.EtherIfaceRF.Id
				ei.etherIfaceRaw = &jsonData
				ei.LastStatus = VerifyingData

				s.ENetInterfaces.OIDs[ei.EtherIfaceRF.Id] = ei

				errlog.Printf("<========== JW_DEBUG ==========> discoverFoxconnENetInterfaces: s.ENetInterfaces.OIDs[%s]=%+v\n", ei.EtherIfaceRF.Id, s.ENetInterfaces.OIDs[ei.EtherIfaceRF.Id])
			}
		}
	}
}