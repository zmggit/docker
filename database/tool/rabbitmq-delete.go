package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

//./rabbitmq-delete --prefix=openapi~ --host=rabbit.example.com --user=admin --pass=secret

// RabbitMQ 配置
type Config struct {
	Host     string // RabbitMQ 主机地址
	Port     int    // HTTP API 端口
	Username string // 用户名
	Password string // 密码
	Vhost    string // 虚拟主机
}

// Exchange 表示 RabbitMQ 交换机
type Exchange struct {
	Name       string `json:"name"`
	Vhost      string `json:"vhost"`
	Type       string `json:"type"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
}

func main() {
	// 解析命令行参数
	prefix := flag.String("prefix", "openapi~", "交换机名称前缀")
	host := flag.String("host", "172.16.19.5", "RabbitMQ 主机地址")
	port := flag.Int("port", 15672, "RabbitMQ HTTP API 端口")
	username := flag.String("user", "won-cloud-admin", "用户名")
	password := flag.String("pass", "sk,j*@b2i8_7hGAoiDho3", "密码")
	vhost := flag.String("vhost", "/", "虚拟主机")
	workers := flag.Int("workers", 5, "并发删除的工作协程数")
	flag.Parse()

	if *prefix == "" {
		fmt.Println("错误：必须指定交换机前缀")
		fmt.Println("用法示例：")
		fmt.Println("  ./rabbitmq-delete --prefix=temp_ --host=rabbit.example.com --user=admin --pass=secret")
		os.Exit(1)
	}

	config := Config{
		Host:     *host,
		Port:     *port,
		Username: *username,
		Password: *password,
		Vhost:    *vhost,
	}

	fmt.Printf("正在删除前缀为 '%s' 的交换机...\n", *prefix)

	// 获取交换机列表
	exchanges, err := getExchanges(config)
	if err != nil {
		log.Fatalf("获取交换机列表失败: %v", err)
	}
	// 过滤匹配前缀的交换机
	matchingExchanges := filterExchanges(exchanges, *prefix)
	if len(matchingExchanges) == 0 {
		fmt.Println("没有找到匹配的交换机")
		return
	}

	fmt.Printf("找到 %d 个匹配的交换机:\n", len(matchingExchanges))

	for _, ex := range matchingExchanges {
		fmt.Printf("  - %s (类型: %s, 持久化: %t)\n", ex.Name, ex.Type, ex.Durable)
	}
	// 确认操作
	fmt.Print("\n确定要删除这些交换机吗? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)
	if strings.ToLower(confirm) != "y" {
		fmt.Println("操作已取消")
		return
	}

	// 并发删除交换机
	deleteExchanges(config, matchingExchanges, *workers)
}

// 获取所有交换机列表
func getExchanges(config Config) ([]Exchange, error) {
	url := fmt.Sprintf("http://%s:%d/api/exchanges/%s",
		config.Host, config.Port, urlEncodeVhost(config.Vhost))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", basicAuth(config.Username, config.Password))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 请求失败: %s, 响应: %s", resp.Status, string(body))
	}

	var exchanges []Exchange
	if err := json.NewDecoder(resp.Body).Decode(&exchanges); err != nil {
		return nil, err
	}

	return exchanges, nil
}

// 过滤匹配指定前缀的交换机
func filterExchanges(exchanges []Exchange, prefix string) []Exchange {
	var result []Exchange
	for _, ex := range exchanges {
		// 排除默认交换机 (amq.*)
		if strings.HasPrefix(ex.Name, "amq.") {
			continue
		}

		if strings.HasPrefix(ex.Name, prefix) {
			result = append(result, ex)
		}
	}
	return result
}

// 并发删除交换机
func deleteExchanges(config Config, exchanges []Exchange, workers int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workers) // 并发控制信号量

	fmt.Printf("\n开始删除 %d 个交换机...\n", len(exchanges))

	successCount := 0
	failCount := 0
	var mu sync.Mutex // 保护计数器的互斥锁

	for _, ex := range exchanges {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(ex Exchange) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			err := deleteExchange(config, ex)
			mu.Lock()
			if err != nil {
				failCount++
				fmt.Printf("删除失败: %s - %v\n", ex.Name, err)
			} else {
				successCount++
				fmt.Printf("已删除: %s\n", ex.Name)
			}
			mu.Unlock()
		}(ex)
	}

	wg.Wait()
	fmt.Printf("\n操作完成! 成功: %d, 失败: %d\n", successCount, failCount)
}

// 删除单个交换机
func deleteExchange(config Config, ex Exchange) error {
	url := fmt.Sprintf("http://%s:%d/api/exchanges/%s/%s",
		config.Host, config.Port, urlEncodeVhost(ex.Vhost), urlEncodeName(ex.Name))

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", basicAuth(config.Username, config.Password))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API 错误: %s, 响应: %s", resp.Status, string(body))
	}

	return nil
}

// 基础认证头生成
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// 虚拟主机URL编码
func urlEncodeVhost(vhost string) string {
	if vhost == "/" {
		return "%2F"
	}
	return vhost
}

// 交换机名称URL编码
func urlEncodeName(name string) string {
	var buf bytes.Buffer
	for _, r := range name {
		switch {
		case 'a' <= r && r <= 'z':
			buf.WriteRune(r)
		case 'A' <= r && r <= 'Z':
			buf.WriteRune(r)
		case '0' <= r && r <= '9':
			buf.WriteRune(r)
		case r == '-':
			buf.WriteRune(r)
		case r == '_':
			buf.WriteRune(r)
		case r == '.':
			buf.WriteRune(r)
		default:
			buf.WriteString(fmt.Sprintf("%%%02X", r))
		}
	}
	return buf.String()
}
