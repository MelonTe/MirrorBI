package openai

// LLMRequest 表示发送给LLM（大型语言模型）的请求参数
// 包含模型选择、消息内容、响应控制和生成参数等配置选项
type LLMRequest struct {
	Model            string    `json:"model"`                       // 使用的LLM模型名称
	Messages         []Message `json:"messages"`                    // 到目前为止的对话的消息列表
	Stream           bool      `json:"stream"`                      // 是否使用流式响应
	MaxTokens        int       `json:"max_tokens"`                  // 生成的最大令牌数
	EnableThinking   bool      `json:"enable_thinking,omitempty"`   // 是否启用思考模式
	ThinkingBudget   int       `json:"thinking_budget,omitempty"`   // 思考预算（令牌数）
	MinP             float64   `json:"min_p,omitempty"`             // 最小概率参数
	Stop             []string  `json:"stop,omitempty"`              // 遇到这些字符串时停止生成
	Temperature      float64   `json:"temperature,omitempty"`       // 温度参数，控制随机性（越高越随机）
	TopP             float64   `json:"top_p,omitempty"`             // 核采样参数
	TopK             int       `json:"top_k,omitempty"`             // K值采样参数
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"` // 频率惩罚系数，避免重复
	N                int       `json:"n,omitempty"`                 // 生成多少个补全结果
	ResponseFormat   *Format   `json:"response_format,omitempty"`   // 响应格式设置
	Tools            []Tool    `json:"tools,omitempty"`             // 模型可以使用的工具列表
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Format defines the response format type
type Format struct {
	Type string `json:"type"`
}

// Tool represents a tool that can be used by the model
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents a function definition for tool
type Function struct {
	Description string                 `json:"description"`
	Name        string                 `json:"name"`
	Parameters  map[string]interface{} `json:"parameters"`
	Strict      bool                   `json:"strict,omitempty"`
}
