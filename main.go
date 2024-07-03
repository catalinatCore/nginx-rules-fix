package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	// 定义一个命令行参数
	rootDir := flag.String("root", ".", "the root directory to process")
	flag.Parse()

	err := filepath.Walk(*rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是 .conf 文件
		if !info.IsDir() && filepath.Ext(path) == ".conf" {
			processFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
	}
}

func processFile(path string) {

	// 获取路径中的特定部分
	dir, _ := filepath.Split(path)
	parts := strings.Split(dir, "/")
	part := parts[len(parts)-2]

	// 如果 part 不是数字，说明不是服务目录
	if _, err := strconv.Atoi(part); err != nil {
		return
	}

	// fmt.Println("处理配置文件:", path)

	lines, err := readLines(path)
	if err != nil {
		fmt.Printf("error reading the file: %v\n", err)
		return
	}

	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "ssl_") {
			continue
		}

		if strings.Contains(line, "reuseport ssl") {
			line = strings.Replace(line, "reuseport ssl", "reuseport", -1)
		}

		newLines = append(newLines, line)
	}

	err = writeLines(path, newLines)
	if err != nil {
		fmt.Printf("error writing the file: %v\n", err)
	}

	fmt.Println("重载服务:", "xiandan"+part+"xiandan")

	// 重载特定的服务
	// cmd := exec.Command("sudo", "systemctl", "reload", "xiandan"+part+"xiandan")
	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Printf("error reloading service: %v\n", err)
	// }
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
