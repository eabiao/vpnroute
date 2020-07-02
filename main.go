package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if !isAdmin() {
		fmt.Println("请以管理员身份运行脚本")
		fmt.Println("按回车键退出程序...")
		fmt.Scanln()
		return
	}

	for {
		vpnIp := getVpnIPAddr()
		if vpnIp == "" {
			showMsg("vpn未连接，请检查vpn状态")
			goto wait
		}

		if isRouteExist() {
			showMsg("路由添加成功，VPN IP：" + vpnIp)
			goto wait
		}
		addRoute(vpnIp)

	wait:
		time.Sleep(1 * time.Second)
	}
}

// 缓存信息
var cacheMsg string

// 打印不重复的信息
func showMsg(msg string) {
	if msg != cacheMsg {
		fmt.Println(msg)
		cacheMsg = msg
	}
}

// 检查路由状态
func isRouteExist() bool {
	result := execute("route print")
	return strings.Contains(result, "192.168.138.0")
}

// 添加静态路由
func addRoute(vpnIp string) {
	result := execute("route add 192.168.138.0 mask 255.255.255.0 " + vpnIp)
	result = strings.TrimSpace(result)
	if strings.Contains(result, "操作完成") {
		showMsg("路由添加成功，VPN IP：" + vpnIp)
	} else {
		showMsg(result)
	}
}

// 获取本机vpn的ip地址
func getVpnIPAddr() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return ""
		}

		for _, addr := range addrs {
			var netIP net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				netIP = v.IP
			case *net.IPAddr:
				netIP = v.IP
			}

			if netIP != nil {
				ip := netIP.String()
				if strings.HasPrefix(ip, "10.20.20.") {
					return ip
				}
			}
		}
	}
	return ""
}

// 检查管理员身份
func isAdmin() bool {
	result := execute("net.exe session 1>NUL 2>NUL && echo admin")
	return strings.TrimRight(result, "\r\n") == "admin"
}

// 执行命令
func execute(cmd string) string {
	execCmd := exec.Command("cmd.exe", "/c", cmd)
	out, _ := execCmd.CombinedOutput()
	out, _ = GbkToUtf8(out)
	return string(out)
}

// 编码转换
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
