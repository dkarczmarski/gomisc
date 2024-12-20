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
	"time"
)

var (
	ErrIncorrectIP = errors.New("incorrect ip")
	ErrIPNotFound  = errors.New("ip not found")
)

type IPEntry struct {
	IP        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type config struct {
	wrapperCmd string
	timeFunc   func() time.Time
}

func WithSudoWrapper() func(*config) {
	return func(c *config) {
		c.wrapperCmd = "sudo"
	}
}

func WithEchoWrapper() func(*config) {
	return func(c *config) {
		c.wrapperCmd = "echo"
	}
}

func WithTimeFunc(f func() time.Time) func(*config) {
	return func(c *config) {
		c.timeFunc = f
	}
}

type Service struct {
	wrapperCmd string
	entries    []*IPEntry
	timeFunc   func() time.Time
}

func NewService(opts ...func(*config)) *Service {
	var cnf config
	for _, ops := range opts {
		ops(&cnf)
	}

	return &Service{
		wrapperCmd: cnf.wrapperCmd,
		timeFunc:   cnf.timeFunc,
	}
}

// AddIP runs firewall command to add that ip.
// When ip has been already added by this method then the next call only update UpdatedAt field.
func (srv *Service) AddIP(ip string) error {
	return srv.AddIPCtx(context.Background(), ip)
}

func (srv *Service) AddIPCtx(ctx context.Context, ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%v: %w", ip, ErrIncorrectIP)
	}

	if _, entry := srv.findByIP(ip); entry != nil {
		entry.UpdatedAt = srv.timeFunc()
		return nil
	}

	// add to registry
	now := srv.timeFunc()
	entry := &IPEntry{
		IP:        ip,
		CreatedAt: now,
		UpdatedAt: now,
	}
	srv.entries = append(srv.entries, entry)

	// execute 'add to firewall' command
	cmdStr := fmt.Sprintf("ufw allow from %s to any proto tcp port 8080", ip)
	out, err := exec.CommandContext(ctx, "echo", strings.Split(cmdStr, " ")...).
		CombinedOutput()
	if err != nil {
		return fmt.Errorf("execute commmand: %w", err)
	}

	log.Println(string(out))

	return nil
}

func (srv *Service) DeleteIP(ip string) error {
	return srv.DeleteIPCtx(context.Background(), ip)
}

func (srv *Service) DeleteIPCtx(ctx context.Context, ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%v: %w", ip, ErrIncorrectIP)
	}

	// remove from registry
	index, entry := srv.findByIP(ip)
	if entry == nil {
		return fmt.Errorf("ip %v: %w", ip, ErrIPNotFound)
	}
	srv.deleteByIndex(index)

	// execute 'delete from firewall' command
	cmdStr := fmt.Sprintf("ufw delete allow from %s to any proto tcp port 8080", ip)
	out, err := exec.CommandContext(ctx, "echo", strings.Split(cmdStr, " ")...).
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

func (srv *Service) DeleteOutOfDate(duration time.Duration) ([]IPEntry, error) {
	return srv.DeleteOutOfDateCtx(context.Background(), duration)
}

func (srv *Service) DeleteOutOfDateCtx(ctx context.Context, duration time.Duration) ([]IPEntry, error) {
	before := srv.timeFunc().Add(-duration)

	entriesBefore := srv.findAllBefore(before)
	if entriesBefore == nil {
		return []IPEntry{}, nil
	}

	deletedEntries := make([]IPEntry, 0, len(entriesBefore))
	for _, entry := range entriesBefore {
		if err := srv.DeleteIPCtx(ctx, entry.IP); err != nil {
			return deletedEntries, err
		}

		deletedEntries = append(deletedEntries, *entry)
	}

	return deletedEntries, nil
}

func (srv *Service) findAllBefore(before time.Time) []*IPEntry {
	var counter int
	for _, ee := range srv.entries {
		if ee.UpdatedAt.Before(before) {
			counter++
		}
	}
	if counter == 0 {
		return nil
	}

	entries := make([]*IPEntry, 0, counter)
	for _, ee := range srv.entries {
		if ee.UpdatedAt.Before(before) {
			entries = append(entries, ee)
		}
	}
	return entries
}
