package tool

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ReadByte 字节流读取文件
func ReadByte(filePath string) ([]byte, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(fi)
	var data []byte
	buf := make([]byte, 2)
	for {
		n, err := r.Read(buf[:])
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		data = append(data, buf[:n]...)
	}
	return data, nil
}

// FileExist 文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// WriteFile 写文件
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// ByteToMb 字节转Mb
func ByteToMb(bakSpaceSize uint64) string {
	var resMb string
	res := fmt.Sprintf("%.2f", float64(bakSpaceSize)/float64(1024*1024))
	if res == "" {
		resMb = "--"
	} else {
		resMb = res + "M"
	}
	return resMb
}

// PathExists 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreatePathIfNotExist 不存在路径创建
//
//	@Param	path	query	string	false	"路径"
func CreatePathIfNotExist(ctx context.Context, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 路径不存在，创建路径
		err := os.MkdirAll(path, 0755) // 使用0755权限创建路径
		if err != nil {
			return err
		}
		slog.Debug("<-----输出日志----->", "Path created", path)
	}
	return nil
}

// BytesToInt 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x) //转换有两种不同的方式，也就是大端和小端。大端就是内存中低地址对应着整数的高位。

	return int(x)
}

// AddByteHeader 添加tcp 请求头
func AddByteHeader(wr []byte, header []byte) []byte {
	v := uint32(len(header))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	buf = append(buf, header...)
	buf = append(buf, wr...)
	return buf
}

// AddTcpHeader 添加tcp 请求头
func AddTcpHeader(wr []byte) []byte {
	v := uint32(len(wr))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	buf = append(buf, wr...)
	return buf
}

func AddFileContent(filename string, fileString []string) error {
	fileByte, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := os.Remove(filename); err != nil {
		return err
	}
	newFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	for _, i2 := range fileString {
		newFile.WriteString(i2 + "\n")
	}
	newFile.Write(fileByte)
	defer newFile.Close()
	return nil
}

// ReadLine 按行读取文件
//
//	@Param	rs		body	string	true	"文件内容"
//	@Param	count	body	int		true	"开始读取第几行"
func ReadLine(rs string, count int) ([]string, []byte) {
	r := strings.NewReader(rs)
	br := bufio.NewReader(r)
	var i int
	var resContent []string //获取前几个数据
	var overByte []byte     //剩余的数据
	for {
		l, e := br.ReadBytes('\n')
		if e != nil && len(l) == 0 {
			break
		}
		if i < count {
			resContent = append(resContent, strings.Trim(string(l), "\n")) //特殊凭接参数
		} else {
			overByte = append(overByte, l...)
		}

		i++
	}
	return resContent, overByte
}

// Encode 将消息编码
func Encode(message string) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4) // 读取前4个字节的数据
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// Buffered返回缓冲中现有的可读取的字节数。
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// 读取真正的消息数据
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}

// FileReplaceGtid 替换gtid
func FileReplaceGtid(file []byte) []byte {
	var replaceWord = []string{"SET @MYSQLDUMP_TEMP_LOG_BIN", "SET @@SESSION.SQL_LOG_BIN", "SET @@GLOBAL.GTID_PURGED"}
	for _, i2 := range replaceWord {
		file = bytes.Replace(file, []byte(i2), []byte("-- "+i2), -1)
	}
	//logs.Debug(string(file), "调试")
	//fileName := "./newrepl.sql"
	//ioutil.WriteFile(fileName, file, 0644)

	return file
}

// AllDirByPath 获取所有子目录
func AllDirByPath(path string) ([]string, error) {
	var directories []string
	//defer wg.Done()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})

	return directories, err
}

// AllFileByPath 目录下的所有文件
func AllFileByPath(path string) ([]string, error) {
	var directories []string
	//defer wg.Done()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})

	return directories, err
}

// encodeURIComponent 编码方式，会对特殊符号编码
func encodeURIComponent(s string, excluded ...[]byte) string {
	var b bytes.Buffer
	written := 0

	for i, n := 0, len(s); i < n; i++ {
		c := s[i]

		switch c {
		case '-', '_', '.', '!', '~', '*', '\'', '(', ')':
			continue
		default:
			// Unreserved according to RFC 3986 sec 2.3
			if 'a' <= c && c <= 'z' {
				continue
			}
			if 'A' <= c && c <= 'Z' {
				continue
			}
			if '0' <= c && c <= '9' {
				continue
			}
			if len(excluded) > 0 {
				conti := false
				for _, ch := range excluded[0] {
					if ch == c {
						conti = true
						break
					}
				}
				if conti {
					continue
				}
			}
		}

		b.WriteString(s[written:i])
		fmt.Fprintf(&b, "%%%02X", c)
		written = i + 1
	}

	if written == 0 {
		return s
	}
	b.WriteString(s[written:])
	return b.String()
}

func decodeURIComponent(s string) (string, error) {
	decodeStr, err := url.QueryUnescape(s)
	if err != nil {
		return s, err
	}
	return decodeStr, err
}

// EncodeByte  gob 编码的数据只能被 Go 语言程序解码。
func EncodeByte(data any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecodeByte  gob 编码的数据只能被 Go 语言程序解码。
func DecodeByte(data []byte, target any) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(target)
}
