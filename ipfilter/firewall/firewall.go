// Package firewall provides methods for managing firewall rules.
package firewall

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

type IPEntry struct {
	IP string
}

type Service struct {
	WrapperCmd string

	entries []*IPEntry
}

func (srv *Service) AddIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return errors.New("incorrect ip")
	}

	// add to registry
	entry := &IPEntry{
		IP: ip,
	}
	srv.entries = append(srv.entries, entry)

	// execute 'add to firewall' command
	cmdStr := fmt.Sprintf("ufw allow from %s to any proto tcp port 8080", ip)
	out, err := exec.CommandContext(context.Background(), "echo", strings.Split(cmdStr, " ")...).
		CombinedOutput()
	if err != nil {
		return fmt.Errorf("execute commmand: %w", err)
	}

	log.Println(string(out))

	return nil
}

func (srv *Service) DeleteIP(ip string) error {
	// remove from registry
	index, entry := srv.findByIP(ip)
	if entry == nil {
		return fmt.Errorf("ip %v: %w", ip, errors.New("not found"))
	}
	srv.deleteByIndex(index)

	// execute 'delete from firewall' command
	cmdStr := fmt.Sprintf("ufw delete allow from %s to any proto tcp port 8080", ip)
	out, err := exec.CommandContext(context.Background(), "echo", strings.Split(cmdStr, " ")...).
		CombinedOutput()
	if err != nil {
		return fmt.Errorf("execute commmand: %w", err)
	}

	log.Println(string(out))

	return nil
}

func (srv *Service) deleteByIndex(index int) {
	srv.entries = append(srv.entries[:index], srv.entries[index+1:]...)
}

func (srv *Service) findByIP(ip string) (int, *IPEntry) {
	for i, ee := range srv.entries {
		if ee.IP == ip {
			return i, ee
		}
	}
	return -1, nil
}

func (srv *Service) List() []IPEntry {
	entries := make([]IPEntry, len(srv.entries))
	for i, ee := range srv.entries {
		entries[i] = *ee
	}
	return entries
}
