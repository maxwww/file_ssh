package main

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

func getClient(config *Config) (*ssh.Client, error) {
	dat, err := ioutil.ReadFile(config.KeyPath)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(dat)
	if err != nil {
		return nil, err
	}

	c := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return ssh.Dial("tcp", net.JoinHostPort(config.Host, config.Port), c)
}
