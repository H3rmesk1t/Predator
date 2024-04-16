package smb

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"Predator/pkg/xhttp"
	"bytes"
	"fmt"
	"net"
	"time"
)

const (
	pkt = "\x00" +
		"\x00\x00\xc0" +
		"\xfeSMB@\x00" +

		//[MS-SMB2]: SMB2 NEGOTIATE Request
		//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/e14db7ff-763a-4263-8b10-0c3944f52fc5

		"\x00\x00" +
		"\x00\x00" +
		"\x00\x00" +
		"\x00\x00" +
		"\x1f\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +

		// [MS-SMB2]: SMB2 NEGOTIATE_CONTEXT
		// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/15332256-522e-4a53-8cd7-0bd17678a2f7

		"$\x00" +
		"\x08\x00" +
		"\x01\x00" +
		"\x00\x00" +
		"\x7f\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"x\x00" +
		"\x00\x00" +
		"\x02\x00" +
		"\x00\x00" +
		"\x02\x02" +
		"\x10\x02" +
		"\x22\x02" +
		"$\x02" +
		"\x00\x03" +
		"\x02\x03" +
		"\x10\x03" +
		"\x11\x03" +
		"\x00\x00\x00\x00" +

		// [MS-SMB2]: SMB2_PREAUTH_INTEGRITY_CAPABILITIES
		// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/5a07bd66-4734-4af8-abcf-5a44ff7ee0e5

		"\x01\x00" +
		"&\x00" +
		"\x00\x00\x00\x00" +
		"\x01\x00" +
		"\x20\x00" +
		"\x01\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00\x00\x00" +
		"\x00\x00" +

		// [MS-SMB2]: SMB2_COMPRESSION_CAPABILITIES
		// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/78e0c942-ab41-472b-b117-4a95ebe88271

		"\x03\x00" +
		"\x0e\x00" +
		"\x00\x00\x00\x00" +
		"\x01\x00" + //CompressionAlgorithmCount
		"\x00\x00" +
		"\x01\x00\x00\x00" +
		"\x01\x00" + //LZNT1
		"\x00\x00" +
		"\x00\x00\x00\x00"
)

func SmbGhost(info *config.HostInfo) error {
	if config.IsBrutePass {
		return nil
	}
	err := SmbGhostScan(info)
	return err
}

func SmbGhostScan(info *config.HostInfo) error {
	ip, port, timeout := info.Host, 445, time.Duration(config.Timeout)*time.Second
	addr := fmt.Sprintf("%s:%v", info.Host, port)
	conn, err := xhttp.WrapperTcpWithTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	_, err = conn.Write([]byte(pkt))
	if err != nil {
		return err
	}
	buff := make([]byte, 1024)
	err = conn.SetReadDeadline(time.Now().Add(timeout))
	n, err := conn.Read(buff)
	if err != nil || n == 0 {
		return err
	}
	if bytes.Contains(buff[:n], []byte("Public")) == true && len(buff[:n]) >= 76 && bytes.Equal(buff[72:74], []byte{0x11, 0x03}) && bytes.Equal(buff[74:76], []byte{0x02, 0x00}) {
		result := fmt.Sprintf("[!] %v CVE-2020-0796 SmbGhost Vulnerable!", ip)
		utils.LogSuccess(result)
	}
	return err
}