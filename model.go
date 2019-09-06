package zlfile

import "os"

type Auth struct {
	User        string `json:"user"`          // 服务器登录用户名
	Password    string `json:"password"`      // 服务器登录密码
	Host        string `json:"host"`          // 服务器IP
	Port        *int   `json:"port"`          // 服务器端口号
	FileBufSize *int   `json:"file_buf_sise"` // 文件拷贝缓存大小
}

type ChmodFile struct {
	Path           string      `json:"path"`            // 文件路径
	FilePermission string      `json:"file_permission"` // 文件要修改后的权限,8进制数字字符串格式,0777
	Mode           os.FileMode `json:"mode"`            // 文件要修改后的权限,10进制格式,511
}
