package sync

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"regexp"
)

type gsync struct {
	addr     string
	user     string
	password string
	ignoreList []*regexp.Regexp
}

func NewGsync(addr string, user string, password string, ignoreList []string) *gsync {
	g := &gsync{
		addr:     addr,
		user:     user,
		password: password,
	}

	if len(ignoreList)>0 {
		g.ignoreList = make([]*regexp.Regexp,0,len(ignoreList))
		for _,ignore := range ignoreList {
			if len(ignore)>0 {
				g.ignoreList = append(g.ignoreList, regexp.MustCompile(ignore))
			}
		}
	}

	return g
}

func (self gsync) connectSSH() *ssh.Client {
	client, err := connect(self.addr, self.user, self.password)
	ce(err, "connect ssh")

	return client
}

func connect(addr string, user string, password string) (*ssh.Client, error) {
	client, err := ssh.Dial("tcp4", addr, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 5,
	})
	return client, err
}

//上传文件
func putfile(localPath string, remotePath string, sftp *sftp.Client) {
	log.Println("上传文件: " + localPath + " ======> " + remotePath)

	file, err := sftp.OpenFile(remotePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	ce(err, "ftp openfile: "+remotePath)
	defer file.Close()

	lfile, err := os.Open(localPath)
	ce(err, "open local file: "+localPath)
	defer lfile.Close()

	sendN, err := io.Copy(file, lfile)
	if err != io.EOF {
		ce(err, "read file")
	}
	log.Printf("sended: %d", sendN)

}

func getRemoteFileMD5(filePath string, client *ssh.Client) (string, error) {
	s, err := runcmd("md5sum "+filePath, client)
	if err != nil {
		return "", nil
	}
	ss := strings.Split(s, " ")
	if len(ss) > 0 {
		return ss[0], nil
	}
	return "", errors.New("can't get md5")
}

func runcmd(cmd string, client *ssh.Client) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output := &bytes.Buffer{}
	session.Stdout = output

	if err = session.Run(cmd); err != nil {
		return "", err
	}

	return output.String(), nil
}

func (self gsync)isIgnorePath(path string) bool {
	for _, e := range self.ignoreList {
		if e.MatchString(path) {
			return true
		}
	}
	return false
}

func (self *gsync) SyncDir(localPath string, remoteDir string) error {
	log.Printf("开始同步文件夹: %v ====> %v", localPath, remoteDir)

	client := self.connectSSH()

	startTime := time.Now()
	totalSyncNum := 0
	realSyncNum := 0

	//md5不一致，传
	sftp, err := sftp.NewClient(client)
	ce(err, "sftp client")
	defer sftp.Close()

	err = filepath.Walk(localPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if self.isIgnorePath(path[len(localPath):]) {
			log.Printf("忽略文件：%v", path)
			return nil;
		}

		remotePath := filepath.Join(remoteDir, path[len(localPath):])
		if f.IsDir() {
			sftp.Mkdir(remotePath)
			return nil
		}

		totalSyncNum++

		//文件md5值对比，如果一样则不传
		md5, err := GetFileMd5(path)
		if err != nil {
			log.Printf("md5计算失败: %v", err)
		}

		remoteMd5, err := getRemoteFileMD5(remotePath, client)
		if md5 != remoteMd5 {
			putfile(path, remotePath, sftp)
			realSyncNum++
		} else {
			log.Printf("文件md5(%v)相同，忽略上传: %v === %v", md5, path, remotePath)
		}

		return nil
	})

	t := time.Now()
	elapsed := t.Sub(startTime) / time.Second
	log.Printf("同步结果：用时[%d]秒，总文件数[%d]个, 实际上传[%d]个",
		elapsed, totalSyncNum, realSyncNum)

	return err
}
