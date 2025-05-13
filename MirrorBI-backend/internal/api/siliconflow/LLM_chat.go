package siliconflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mrbi/config"
	"mrbi/internal/api/openai"
	"net/http"
)

// 制造一个请求结构体
func NewLLMChatReqeust(requirement string, csvData string) *openai.LLMRequest {
	sysPrompt := fmt.Sprintf("你是一个优秀的数据分析助手,以及前端开发师.\n" +
		"请根据用户提供的需求和 CSV 数据,生成精炼的分析结论和 EchartV5 option 代码.\n" +
		"请求格式如下:\n" +
		"分析需求:\n" +
		"{需求描述}\n" +
		"数据源:\n" +
		"{csv格式的原始数据源描述}\n" +
		"响应格式要求如下:\n" +
		"option = {代码}\n" +
		"数据分析结论:{中文结论分析}\n")
	return &openai.LLMRequest{
		Model:       "Pro/deepseek-ai/DeepSeek-V3",
		Temperature: 0.7,
		Messages: []openai.Message{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: fmt.Sprintf("分析需求:\n%s\n数据源:\n%s", requirement, csvData)},
		},
		Stream:    true,
		MaxTokens: 2000,
	}
}

// 没有上下文的对话请求,返回响应结果
func NewLLMChatReqeustNoContext(requirement string, csvData string) (*openai.LLMResponse, error) {
	//1.设置请求体参数
	req := NewLLMChatReqeust(requirement, csvData)
	//关闭流式响应
	req.Stream = false
	//调用API请求
	apiKey := config.LoadConfig().Siliconflow.APIkey
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal() failed: %v", err)
	}
	//2.发起HTTP请求
	url := "https://api.siliconflow.cn/v1/chat/completions"
	httpReq, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(payload))
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do() failed: %v", err)
	}
	defer resp.Body.Close()
	//状态检查
	if resp.StatusCode != http.StatusOK {
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		return nil, fmt.Errorf("API 返回 %d: %s", resp.StatusCode, body.String())
	}

	//3.进行流式响应，逐行读取SSE
	// reader := bufio.NewReader(resp.Body)
	// //var rawChunks []string
	// for {
	// 	//读取一行数据
	// 	line, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		//结束读取
	// 		break
	// 	}
	// 	line = strings.TrimSpace(line) //去除首尾空格

	// 	//调试
	// 	log.Println(line)

	// 	//剥离非数据行
	// 	if !strings.HasPrefix(line, "data:") {
	// 		continue
	// 	}
	// 	data := strings.TrimPrefix(line, "data: ")
	// 	if data == "[DONE]" {
	// 		break //结束传输
	// 	}
	// 	//TODO：打印增量内容，合并结果解析。
	// }

	//3.解析响应，返回结果
	var llmResponse openai.LLMResponse
	//获取body的byte数据
	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll() failed: %v", err)
	}
	//log.Println("响应数据:", string(byteData))
	//解析JSON数据
	err = json.Unmarshal(byteData, &llmResponse)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal() failed: %v", err)
	}
	//返回结果
	//log.Println(llmResponse)
	return &llmResponse, nil
}
