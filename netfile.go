package netfile

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Auth struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type SshFptInterface interface {
	GetSshSession() (*ssh.Session, error)
	SshRun(client *ssh.Session, cmd string) error
	GetSshClient() (*ssh.Client, error)
	GetFtpClient(sshClint *ssh.Client) (*sftp.Client, error)
	ScpCopy(client *sftp.Client, from, to string) error
}

type sshFptService struct {
	Para *Auth
}

func NewSshFtpInterface(para *Auth) SshFptInterface {
	return &sshFptService{para}
}

func (s *sshFptService) GetSshSession() (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(s.Para.Password))

	clientConfig = &ssh.ClientConfig{
		User:    s.Para.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", s.Para.Host, s.Para.Port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sshFptService) SshRun(client *ssh.Session, cmd string) error {
	if client == nil {
		var err error
		client, err = s.GetSshSession()
		if err != nil {
			return err
		}
	}
	return client.Run(cmd)
}

func (s *sshFptService) GetSshClient() (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(s.Para.Password))
	clientConfig = &ssh.ClientConfig{
		User:    s.Para.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		//这个是问你要不要验证远程主机，以保证安全性。这里不验证
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", s.Para.Host, s.Para.Port)
	return ssh.Dial("tcp", addr, clientConfig)
}

func (s *sshFptService) GetFtpClient(client *ssh.Client) (*sftp.Client, error) {
	if client == nil {
		var err error
		client, err = s.GetSshClient()
		if err != nil {
			return nil, err
		}
	}
	return sftp.NewClient(client)
}

func (s *sshFptService) ScpCopy(client *sftp.Client, from, to string) error {
	if client == nil {
		var err error
		client, err = s.GetFtpClient(nil)
		if err != nil {
			return err
		}
	}

	fileList, err := GetDirFileList(from)
	if err != nil {
		return err
	}
	for _, fileName := range fileList {
		newFrom := path.Join(from, fileName)
		newTo := path.Join(to, fileName)
		client.Remove(newTo)

		isDir, err := IsDir(newFrom)
		if err != nil {
			return err
		}
		if isDir {
			if err := client.MkdirAll(newTo); err != nil {
				return err
			}
			if err := s.ScpCopy(client, newFrom, newTo); err != nil {
				return err
			}
		} else {
			if err = s.ScpCopyFile(client, newFrom, newTo); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *sshFptService) ScpCopyFile(client *sftp.Client, from, to string) error {
	if client == nil {
		var err error
		client, err = s.GetFtpClient(nil)
		if err != nil {
			return err
		}
	}

	srcFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(to)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 10240)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	return nil
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

func GetDirAndSubDirFileList(pathname string) ([]string, error) {
	var fileList []string
	err := filepath.Walk(pathname,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				fileList = append(fileList, path)
				return nil
			}

			return nil
		})
	return fileList, err
}

func GetDirFileList(pathname string) ([]string, error) {
	var fileList []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return fileList, err
	}
	for _, fi := range rd {
		fileList = append(fileList, fi.Name())
	}
	return fileList, nil
}
