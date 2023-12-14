package mikrotik

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/routeros.v2"
)

type mikrotikRepository struct {
	client *routeros.Client
}

type Mikrotik interface {
	GetIdentity() (string, error)
	GetSecrets(bts string, host string) ([]Secret, error)
	GetActiveConnections() ([]map[string]string, error)
	EnableSNMP()
	SetMacFromAC()
	SetRemoteAddress()
	GetAddressList(listName string) []AddressList
	AddSecretToAddressList(addressList AddressList) AddressList
	RemoveSecretFromAddressList(addressListData []AddressList, remoteAddress string) AddressList
	SetStaticRoute(dstAddr, rType, prefSource, bgpComm string) error
}

func NewMikrotikRepository(addr, user, password string) (Mikrotik, error) {
	dial, err := routeros.Dial(addr+":8728", user, password)
	if err != nil {
		return nil, err
	}

	return &mikrotikRepository{client: dial}, nil
}

func (mr *mikrotikRepository) GetIdentity() (string, error) {
	identity := []map[string]string{}
	mkt, err := mr.client.Run("/system/identity/print")
	if err != nil {
		return "", err
	}

	for _, r := range mkt.Re {
		identity = append(identity, r.Map)
	}
	return identity[0]["name"], nil
}

func (mr *mikrotikRepository) GetSecrets(bts string, host string) ([]Secret, error) {
	var secret []Secret

	mkt, err := mr.client.Run("/ppp/secret/print")

	if err != nil {
		return []Secret{}, err
	}

	for _, d := range mkt.Re {
		row := Secret{
			Name:          d.Map["name"],
			CallerID:      d.Map["caller-id"],
			Profile:       d.Map["profile"],
			Comment:       d.Map["comment"],
			RemoteAddress: d.Map["remote-address"],
			Bts:           bts,
			Host:          host,
		}
		secret = append(secret, row)
	}
	return secret, nil
}

func (mr *mikrotikRepository) GetActiveConnections() ([]map[string]string, error) {
	var activeUsers []map[string]string
	reply, err := mr.client.Run("/ppp/active/print")
	if err != nil {
		return nil, err
	}

	for _, r := range reply.Re {
		activeUsers = append(activeUsers, r.Map)
	}

	return activeUsers, nil
}

func (mr *mikrotikRepository) EnableSNMP() {
	reply, err := mr.client.Run("/snmp/set", "=enabled=yes", "=trap-version=2")
	if err != nil {
		log.Println(err)
	}
	log.Println(reply.Done.Word)
}

func (mr *mikrotikRepository) SetMacFromAC() {
	reply, err := mr.client.Run("/ppp/secret/print")
	if err != nil {
		log.Println("error to get secrets users", err)
	}
	var secret []map[string]string

	for _, s := range reply.Re {
		fmt.Println(s.Map)
		secret = append(secret, s.Map)
	}

	activeUsers, err := mr.GetActiveConnections()
	if err != nil {
		log.Println("error to get active connections", err)
	}

	for _, su := range secret {
		for _, au := range activeUsers {
			if su["name"] == au["name"] {
				_, err := mr.client.Run(
					"/ppp/secret/set",
					fmt.Sprintf("=numbers=%s", su["name"]),
					fmt.Sprintf("=caller-id=%s", au["caller-id"]),
				)
				if err != nil {
					log.Println("error to set mac", err)
				}

			}
		}
	}
}

func (mr *mikrotikRepository) SetRemoteAddress() {
	secrets, err := mr.GetSecrets("", "")
	if err != nil {
		log.Println(err)
	}
	activeConnections, err := mr.GetActiveConnections()
	if err != nil {
		log.Println("error to get active connections")
	}

	for _, secret := range secrets {
		for _, active := range activeConnections {
			if secret.Name == active["name"] {
				_, err := mr.client.Run(
					"/ppp/secret/set",
					fmt.Sprintf("=numbers=%s", secret.Name),
					fmt.Sprintf("=remote-address=%s", active["address"]),
				)
				if err != nil {
					log.Println("error to set remote address")
				}
			}
		}
	}
}

func (mr *mikrotikRepository) AddSecretToAddressList(data AddressList) AddressList {

	reply, _ := mr.client.Run(
		"/ip/firewall/address-list/add",
		fmt.Sprintf("=list=%s", data.List),
		fmt.Sprintf("=address=%s", data.Address),
		fmt.Sprintf("=comment=%s", data.Comment),
	)
	if reply != nil {
		return AddressList{
			Address:      data.Address,
			Comment:      data.Comment,
			CreationTime: time.Now().Format("Jan/02/2006 15:04:05"),
			List:         data.List,
			Status:       "CORTADO",
		}
	}

	return AddressList{}
}

func (mr *mikrotikRepository) RemoveSecretFromAddressList(addressListData []AddressList, remoteAddress string) AddressList {

	for _, list := range addressListData {
		if list.Address == remoteAddress {
			_, err := mr.client.Run(
				"/ip/firewall/address-list/remove",
				fmt.Sprintf("=numbers=%s", list.ID),
			)

			if err == nil {
				return AddressList{
					Address:      list.Address,
					Comment:      list.Comment,
					CreationTime: list.CreationTime,
					List:         list.List,
					Status:       "ACTIVO",
				}
			}
		}
	}
	return AddressList{}
}

func (mr *mikrotikRepository) GetAddressList(listName string) []AddressList {
	var results []AddressList
	reply, err := mr.client.Run(
		"/ip/firewall/address-list/print",
	)

	if err != nil {
		log.Println("Error al imprimir los address list", err)
	}

	for _, alist := range reply.Re {
		if alist.Map["list"] == listName {
			L := AddressList{
				ID:           alist.Map[".id"],
				Address:      alist.Map["address"],
				Comment:      alist.Map["comment"],
				CreationTime: alist.Map["creation-time"],
				List:         alist.Map["list"],
			}
			results = append(results, L)
		}
	}
	return results
}

// ip route add dst-address=1.1.1.0/24 type=blackhole pref-src=1.1.1.1 bgp-communities=18013:13
func (mr *mikrotikRepository) SetStaticRoute(dstAddr, rType, prefSrc, bgpComm string) error {
	_, err := mr.client.Run(
		"/ip/route/add",
		fmt.Sprintf("=dst-address=%s", dstAddr),
		fmt.Sprintf("=type=%s", rType),
		fmt.Sprintf("=pref-src=%s", prefSrc),
		fmt.Sprintf("=bgp-communities=%s", bgpComm),
	)
	return err
}
