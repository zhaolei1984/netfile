package zlfile

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshFtpInterface interface {
	GetSshSession() (*ssh.Session, error)
	SshRun(client *ssh.Session, cmd string) error
	GetSshClient() (*ssh.Client, error)
	GetFtpClient(sshClint *ssh.Client) (*sftp.Client, error)
	ScpCopyTo(client *sftp.Client, from, to string, chmodFileList []*ChmodFile) error
	ScpCopyFileTo(client *sftp.Client, from, to string, chmodFileList []*ChmodFile) error
	Chmod(client *sftp.Client, path string, modeFile *ChmodFile) error
	ScpCopyFrom(client *sftp.Client, from, to string, chmodFileList []*ChmodFile) error
	ScpCopyFileFrom(client *sftp.Client, from, to string, chmodFileList []*ChmodFile) error
}

type sshFptService struct {
	Para *Auth
}

func NewSshFtpInterface(para *Auth) SshFtpInterface {
	if para.Port == nil {
		port := 22
		para.Port = &port
	}
	if para.FileBufSize == nil {
		fileBufSize := 10240
		para.FileBufSize = &fileBufSize
	}
	return &sshFptService{para}
}
