package can

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
)

func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	s, err := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		return nil, err
	}
	addr := &unix.SockaddrCAN{Ifindex: i.Index}
	if err := unix.Bind(s, addr); err != nil {
		return nil, err
	}

	// 打开本地回环
	err = syscall.SetsockoptInt(s, unix.SOL_CAN_RAW, unix.CAN_RAW_LOOPBACK, 1)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("本地回环错误%s", err.Error()))
	}

	err = syscall.SetsockoptByte(s, unix.SOL_CAN_RAW, unix.CAN_RAW_ERR_FILTER, unix.CAN_ERR_TX_TIMEOUT|unix.CAN_ERR_BUSOFF|unix.CAN_ERR_BUSERROR)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("filter错误%s", err.Error()))
	}

	f := os.NewFile(uintptr(s), fmt.Sprintf("fd %d", s))

	return &readWriteCloser{f}, nil
}
