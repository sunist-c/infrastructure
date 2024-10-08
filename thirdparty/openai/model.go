package openai

import (
	"encoding/json"
	"io"
	"strconv"
)

// ModelObject 模型对象
// reference https://platform.openai.com/docs/api-reference/models/object
type ModelObject struct {
	ID      string `json:"id"`                 // 模型的唯一标识符
	Created int64  `json:"created"`            // 模型创建的时间
	Object  string `json:"object"`             // 模型的类型，一般为model
	OwnedBy string `json:"owned_by,omitempty"` // 模型的所有者，一般为openai
}

// ListModelRequest 列出模型请求
// reference https://platform.openai.com/docs/api-reference/models/list
type ListModelRequest struct{}

// ListModelResponseBody 列出模型响应
// reference https://platform.openai.com/docs/api-reference/models/list
type ListModelResponseBody struct {
	Object string        `json:"object"` // 一般为list
	Data   []ModelObject `json:"data"`   // 模型列表
}

// RetrieveModelRequest 获取模型请求
// reference https://platform.openai.com/docs/api-reference/models/retrieve
type RetrieveModelRequest struct {
	Model string // 模型的唯一标识符
}

// RetrieveModelResponseBody 获取模型响应
// reference https://platform.openai.com/docs/api-reference/models/retrieve
type RetrieveModelResponseBody struct {
	ModelObject // 模型对象
}

// ImageItem 生成的图片
// reference https://platform.openai.com/docs/api-reference/images/object
type ImageItem struct {
	Url           string `json:"url,omitempty"`            // openai生成的图片url，在response_format为url时返回，默认为url
	Base64        string `json:"b64_json,omitempty"`       // openai生成的图片base64编码，在response_format为b64_json时返回
	RevisedPrompt string `json:"revised_prompt,omitempty"` // openai生成的图片的提示，如果提示被修改了，会返回修改后的
}

// CreateImageRequestBody 生成图片请求体
// reference https://platform.openai.com/docs/api-reference/images/create
type CreateImageRequestBody struct {
	Prompt         string `json:"prompt"`                    // 生成图片的提示
	Model          string `json:"model,omitempty"`           // 生成图片的模型
	N              int    `json:"n"`                         // 生成图片的数量，1~10之间
	Size           string `json:"size"`                      // 生成图片的尺寸，dall-e-2模型只支持256x256,512x512,1024x1024，dall-e-3模型只支持1024x1024,1792x1024,1024x1792
	Quality        string `json:"quality,omitempty"`         // 质量，仅支持dall-e-3模型，标记为hd时生成的图片质量更高
	Style          string `json:"style,omitempty"`           // 风格，仅支持dall-e-3模型，可以选择natural和vivid
	ResponseFormat string `json:"response_format,omitempty"` // 返回格式，url或者b64_json
	User           string `json:"user,omitempty"`            // 用户的唯一标识符，用于openai跟踪
}

// CreateImageRequest 生成图片请求
// reference https://platform.openai.com/docs/api-reference/images/create
type CreateImageRequest struct {
	Body CreateImageRequestBody
}

// ImageResponseBody 生成图片响应
// reference https://platform.openai.com/docs/api-reference/images/create
type ImageResponseBody struct {
	Created int64       `json:"created"` // openai返回的消息回复时间
	Data    []ImageItem `json:"data"`    // openai生成的图片，一般情况下长度和请求的N一致
}

// ChatRoleEnum 聊天角色枚举
// reference https://platform.openai.com/docs/guides/text-generation/chat-completions-api
type ChatRoleEnum string

func (cre ChatRoleEnum) String() string { return string(cre) }

func getChatRoleEnum(enum ChatRoleEnum) string {
	if _, exist := supportedChatRoleEnum[enum.String()]; !exist {
		return ChatRoleEnumUser.String()
	}

	return enum.String()
}

// 聊天角色枚举值
const (
	ChatRoleEnumSystem    ChatRoleEnum = "system"    // 系统，用于给模型进行提示的角色，一般用于给出prompt
	ChatRoleEnumAssistant ChatRoleEnum = "assistant" // 助手，模型一般使用该角色
	ChatRoleEnumUser      ChatRoleEnum = "user"      // 用户，用户一般使用该角色
)

// 支持的聊天角色枚举
var supportedChatRoleEnum = map[string]ChatRoleEnum{
	ChatRoleEnumSystem.String():    ChatRoleEnumSystem,
	ChatRoleEnumAssistant.String(): ChatRoleEnumAssistant,
	ChatRoleEnumUser.String():      ChatRoleEnumUser,
}

// ChatMessageObject 聊天消息对象
// reference https://platform.openai.com/docs/api-reference/chat/object
type ChatMessageObject struct {
	Role    ChatRoleEnum    `json:"role"`
	Content json.RawMessage `json:"content"`
}

func (co *ChatMessageObject) GetStringContent() string {
	var contentBuffer any
	tryErr := json.Unmarshal(co.Content, &contentBuffer)
	if tryErr != nil {
		return string(co.Content)
	}

	type imageContent []struct {
		Text     string `json:"text,omitempty"`
		Type     string `json:"type"`
		ImageUrl struct {
			Detail string `json:"detail"`
			Url    string `json:"url"`
		} `json:"image_url,omitempty"`
	}

	if content, ok := contentBuffer.(string); ok {
		return content
	}

	obj := imageContent{}
	err := json.Unmarshal(co.Content, &obj)
	if err != nil {
		return string(co.Content)
	}

	for _, item := range obj {
		if item.Type == "text" {
			return item.Text
		}
	}

	return string(co.Content)
}

type StreamingReplyObject struct {
	Id                string                       `json:"id"`
	Object            string                       `json:"object"`
	Created           int                          `json:"created"`
	Model             string                       `json:"model"`
	Choices           []StreamingReplyChoiceObject `json:"choices"`
	SystemFingerprint string                       `json:"system_fingerprint"`
	Usage             *UsageObject                 `json:"usage,omitempty"`
}

type StreamingReplyChoiceObject struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Logprobs     ReplyLogprobs `json:"logprobs"`
	FinishReason string        `json:"finish_reason"`
}

type ReplyLogprobs struct {
	Refusal interface{} `json:"refusal"`
	Content []struct {
		TopLogprobs []interface{} `json:"top_logprobs"`
		Logprob     float64       `json:"logprob"`
		Bytes       []int         `json:"bytes"`
		Token       string        `json:"token"`
	} `json:"content"`
}

// ReplyChoiceObject 回复选择对象
// reference https://platform.openai.com/docs/api-reference/chat/object
type ReplyChoiceObject struct {
	Index        int               `json:"index"`         // 回复序号
	Message      ChatMessageObject `json:"message"`       // 回复信息
	FinishReason string            `json:"finish_reason"` // 停止回复原因
}

// UsageObject token使用量对象
// reference https://platform.openai.com/docs/api-reference/chat/object
type UsageObject struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提问token消耗
	CompletionTokens int `json:"completion_tokens"` // 回复token消耗
	TotalTokens      int `json:"total_tokens"`      // 全部token消耗
}

// ResponseFormat 聊天请求返回格式
// reference https://platform.openai.com/docs/api-reference/chat/create#chat-create-response_format
type ResponseFormat struct {
	Type string `json:"type,omitempty"` // 返回格式，支持 json_object 和 text
}

// 聊天请求返回格式
var (
	TextResponseFormat       = &ResponseFormat{Type: "text"}        // 返回文本格式
	JsonObjectResponseFormat = &ResponseFormat{Type: "json_object"} // 返回json格式
)

// Tools 聊天请求工具
// reference https://platform.openai.com/docs/api-reference/chat/create#chat-create-tools
type Tools struct {
	Type     string       `json:"type,omitempty"`     // 工具类型，目前只有 function
	Function ToolFunction `json:"function,omitempty"` // 工具函数
}

// ToolFunction 聊天请求工具函数
// reference https://platform.openai.com/docs/api-reference/chat/create#chat-create-tools-function
type ToolFunction struct {
	Name        string `json:"name,omitempty"`        // 工具函数名称
	Description string `json:"description,omitempty"` // 工具函数描述
	Parameters  any    `json:"parameters,omitempty"`  // 工具函数参数
}

// CompleteChatRequestBody 聊天请求体
// reference https://platform.openai.com/docs/api-reference/chat/create
type CompleteChatRequestBody struct {
	Model            string              `json:"model"`                       // 生成图片的模型
	Messages         []ChatMessageObject `json:"messages"`                    // 提问信息
	Temperature      float64             `json:"temperature,omitempty"`       // 温度采样，0~2，越高越随机
	TopP             float64             `json:"top_p,omitempty"`             // 核采样，0~1，和temperature只有一个能生效
	N                int                 `json:"n,omitempty"`                 // 需要多少消息回复
	Stream           bool                `json:"stream,omitempty"`            // 流式传输开关
	MaxTokens        int                 `json:"max_tokens,omitempty"`        // 允许消耗的最大令牌数
	PresencePenalty  float64             `json:"presence_penalty,omitempty"`  // 创新惩罚，-2~+2，越高越可能出现新东西
	FrequencyPenalty float64             `json:"frequency_penalty,omitempty"` // 重复惩罚，-2~+2，越高越不可能重复
	User             string              `json:"user,omitempty"`              // 用户的唯一标识符，用于openai跟踪
	LogitBias        map[string]float64  `json:"logit_bias,omitempty"`        // 对特定token的偏好
	LogProbs         bool                `json:"logprobs,omitempty"`          // 返回token的log概率
	TopLogProbs      int                 `json:"top_logprobs,omitempty"`      // 返回token的top log概率
	ResponseFormat   *ResponseFormat     `json:"response_format,omitempty"`   // 返回格式，text 或者 json_object
	Seed             int                 `json:"seed,omitempty"`              // 随机种子
	ServiceTier      string              `json:"service_tier,omitempty"`      // 服务层级
	Tools            []Tools             `json:"tools,omitempty"`             // 模型可以调用的工具列表
	ToolChoice       any                 `json:"tool_choice,omitempty"`       // 控制模型调用哪个（如果有）工具
	StreamOptions    json.RawMessage     `json:"stream_options,omitempty"`    // 流式传输选项
}

// CompleteChatRequest 聊天请求
// reference https://platform.openai.com/docs/api-reference/chat/create
type CompleteChatRequest struct {
	Body CompleteChatRequestBody
}

// CompleteChatResponseBody 聊天响应
// reference https://platform.openai.com/docs/api-reference/chat/create
type CompleteChatResponseBody struct {
	ID                string              `json:"id"`                 // openai提供的回复id
	Object            string              `json:"object"`             // openai标记的返回对象，此处固定为chat.completion
	Created           int64               `json:"created"`            // openai返回的消息回复时间
	Choices           []ReplyChoiceObject `json:"choices"`            // openai的回复，一般情况下只有一个元素
	Usage             UsageObject         `json:"usage"`              // openai的token使用量
	Model             string              `json:"model"`              // openai的模型
	ServiceTier       string              `json:"service_tier"`       // openai的服务层级
	SystemFingerprint string              `json:"system_fingerprint"` // openai的系统指纹，与 seed 结合使用
}

type VoiceEnum string

func (ve VoiceEnum) String() string { return string(ve) }

func getVoiceEnum(enum VoiceEnum) string {
	if _, exist := supportedVoiceEnum[enum.String()]; !exist {
		return VoiceEnumAlloy.String()
	}

	return enum.String()
}

// 语音枚举值
const (
	VoiceEnumAlloy   VoiceEnum = "alloy"   // alloy音色，成年女性，中性声音
	VoiceEnumEcho    VoiceEnum = "echo"    // echo音色，成年男性，中性声音
	VoiceEnumFable   VoiceEnum = "fable"   // fable音色，成年女性，中性声音
	VoiceEnumOnyx    VoiceEnum = "onyx"    // onyx音色，成年男性，低沉声音
	VoiceEnumNova    VoiceEnum = "nova"    // nova音色，成年女性，年轻声音
	VoiceEnumShimmer VoiceEnum = "shimmer" // shimmer音色，成年女性，中性声音
)

// 支持的语音枚举
var supportedVoiceEnum = map[string]VoiceEnum{
	VoiceEnumAlloy.String():   VoiceEnumAlloy,
	VoiceEnumEcho.String():    VoiceEnumEcho,
	VoiceEnumFable.String():   VoiceEnumFable,
	VoiceEnumOnyx.String():    VoiceEnumOnyx,
	VoiceEnumNova.String():    VoiceEnumNova,
	VoiceEnumShimmer.String(): VoiceEnumShimmer,
}

// CreateSpeechRequestBody 生成语音请求体
// reference https://platform.openai.com/docs/api-reference/audio/createSpeech
type CreateSpeechRequestBody struct {
	Model          string    `json:"model"`                     // 生成语音的模型，目前支持tts-1和tts-1-hd
	Input          string    `json:"input"`                     // 需要生成语音的文本
	Voice          VoiceEnum `json:"voice"`                     // 生成语音的声音
	ResponseFormat string    `json:"response_format,omitempty"` // 返回格式，支持mp3,aac,flac,opus
	Speed          float64   `json:"speed,omitempty"`           // 语速，支持0.25~4.0，1.0为正常语速
}

// CreateSpeechRequest 生成语音请求
// reference https://platform.openai.com/docs/api-reference/audio/createSpeech
type CreateSpeechRequest struct {
	Body CreateSpeechRequestBody
}

// CreateSpeechResponseBody 生成语音响应
// reference https://platform.openai.com/docs/api-reference/audio/createSpeech
type CreateSpeechResponseBody []byte // openai生成的语音，根据response_format返回不同的格式

// CreateTranscriptionRequestBody 生成转录请求体
// reference https://platform.openai.com/docs/api-reference/audio/createTranscription
type CreateTranscriptionRequestBody struct {
	File           io.Reader // 需要转录的音频文件
	FileName       string    // 音频文件的名称
	Model          string    // 进行音频转写的模型
	Language       string    // 音频文件的语言，ISO-639-1标准
	Prompt         string    // 音频文件的提示
	ResponseFormat string    // 返回格式，支持json和txt
	Temperature    float64   // 温度采样，0~1，越高越随机
}

func (ctr CreateTranscriptionRequestBody) ToMultiPartBody() map[string]string {
	result := map[string]string{
		"model": ctr.Model,
	}
	if ctr.Language != "" {
		result["language"] = ctr.Language
	}
	if ctr.Prompt != "" {
		result["prompt"] = ctr.Prompt
	}
	if ctr.ResponseFormat != "" {
		result["response_format"] = ctr.ResponseFormat
	}
	if ctr.Temperature != 0 {
		result["temperature"] = strconv.FormatFloat(ctr.Temperature, 'f', -1, 64)
	}
	return result
}

// CreateTranscriptionRequest 生成转录请求
// reference https://platform.openai.com/docs/api-reference/audio/createTranscription
type CreateTranscriptionRequest struct {
	FormBody CreateTranscriptionRequestBody
}

// CreateTranscriptionResponseBody 生成转录响应
// reference https://platform.openai.com/docs/api-reference/audio/createTranscription
type CreateTranscriptionResponseBody struct {
	Text string `json:"text"` // openai生成的文本
}

// ModerationCategoryObject 内容分类对象
// reference: https://platform.openai.com/docs/api-reference/moderations/object
type ModerationCategoryObject struct {
	Sexual                bool `json:"sexual"`                 // 性内容
	Hate                  bool `json:"hate"`                   // 仇恨内容
	Harassment            bool `json:"harassment"`             // 骚扰内容
	SelfHarm              bool `json:"self-harm"`              // 自残内容
	SexualMinors          bool `json:"sexual/minors"`          // 未成年人性内容
	HateThreatening       bool `json:"hate/threatening"`       // 仇恨威胁内容
	ViolenceGraphic       bool `json:"violence/graphic"`       // 暴力内容
	SelfHarmIntent        bool `json:"self-harm/intent"`       // 自残意图内容
	SelfHarmInstr         bool `json:"self-harm/instructions"` // 自残指导内容
	HarassmentThreatening bool `json:"harassment/threatening"` // 骚扰威胁内容
	Violence              bool `json:"violence"`               // 暴力内容
}

// ModerationCategoryScoreObject 内容分类得分对象
// reference: https://platform.openai.com/docs/api-reference/moderations/object
type ModerationCategoryScoreObject struct {
	Sexual                float64 `json:"sexual"`                 // 性内容
	Hate                  float64 `json:"hate"`                   // 仇恨内容
	Harassment            float64 `json:"harassment"`             // 骚扰内容
	SelfHarm              float64 `json:"self-harm"`              // 自残内容
	SexualMinors          float64 `json:"sexual/minors"`          // 未成年人性内容
	HateThreatening       float64 `json:"hate/threatening"`       // 仇恨威胁内容
	ViolenceGraphic       float64 `json:"violence/graphic"`       // 暴力内容
	SelfHarmIntent        float64 `json:"self-harm/intent"`       // 自残意图内容
	SelfHarmInstr         float64 `json:"self-harm/instructions"` // 自残指导内容
	HarassmentThreatening float64 `json:"harassment/threatening"` // 骚扰威胁内容
	Violence              float64 `json:"violence"`               // 暴力内容
}

// ModerationResultObject 内容审核结果对象
// reference: https://platform.openai.com/docs/api-reference/moderations/object
type ModerationResultObject struct {
	Flagged        bool                          `json:"flagged"`         // 是否被标记
	Categories     ModerationCategoryObject      `json:"categories"`      // 内容分类
	CategoryScores ModerationCategoryScoreObject `json:"category_scores"` // 内容分类得分
}

// CompleteModerationRequestBody 内容审核请求体
// reference https://platform.openai.com/docs/api-reference/moderations/create
type CompleteModerationRequestBody struct {
	Model string `json:"model"` // 进行内容审核的模型
	Input string `json:"input"` // 需要审核的文本
}

// CompleteModerationRequest 内容审核请求
// reference https://platform.openai.com/docs/api-reference/moderations/create
type CompleteModerationRequest struct {
	Body CompleteModerationRequestBody
}

// EmbeddingRequest 嵌入请求
type EmbeddingRequest struct {
	Body EmbeddingRequestBody
}

// EmbeddingRequestBody 嵌入请求体
type EmbeddingRequestBody struct {
	Input string `json:"input"`           // 需要生成嵌入的文本
	Model string `json:"model,omitempty"` // 进行嵌入的模型
}

// EmbeddingDataItem 嵌入数据项
type EmbeddingDataItem struct {
	Object    string    `json:"object"`    // 数据类型，一般为embedding
	Embedding []float64 `json:"embedding"` // 嵌入数据
	Index     int       `json:"index"`     // 嵌入数据的索引
}

// EmbeddingResponseBody 嵌入响应
type EmbeddingResponseBody struct {
	Object string              `json:"object"` // 数据类型，一般为 list
	Data   []EmbeddingDataItem `json:"data"`   // 嵌入数据
	Model  string              `json:"model"`  // 进行嵌入的模型
	Usage  UsageObject         `json:"usage"`  // openai的token使用量
}

// CompleteModerationResponseBody 内容审核响应
// reference https://platform.openai.com/docs/api-reference/moderations/create
type CompleteModerationResponseBody struct {
	ID      string                   `json:"id"`      // openai提供的回复id
	Model   string                   `json:"model"`   // 进行内容审核的模型
	Results []ModerationResultObject `json:"results"` // 内容审核结果
}

// AvailableFineTuningModelEnum 可用的微调模型枚举值
// reference https://platform.openai.com/docs/guides/fine-tuning
type AvailableFineTuningModelEnum string

func (afme AvailableFineTuningModelEnum) String() string { return string(afme) }

const (
	AvailableFineTuningModelEnumGPT35   AvailableFineTuningModelEnum = "gpt-3.5-turbo-1106" // gpt3.5模型(default)
	AvailableFineTuningModelEnumDavinci AvailableFineTuningModelEnum = "davinci-002"        // davinci模型
)

// 支持的微调模型枚举
var supportedAvailableFineTuningModelEnum = map[string]AvailableFineTuningModelEnum{
	AvailableFineTuningModelEnumGPT35.String():   AvailableFineTuningModelEnumGPT35,
	AvailableFineTuningModelEnumDavinci.String(): AvailableFineTuningModelEnumDavinci,
}

func getAvailableFineTuningModelEnum(enum AvailableFineTuningModelEnum) string {
	if _, exist := supportedAvailableFineTuningModelEnum[enum.String()]; !exist {
		return AvailableFineTuningModelEnumDavinci.String()
	} else {
		return enum.String()
	}
}

// PurposeEnum 可用的用途枚举值
// reference https://platform.openai.com/docs/api-reference/files/create
type PurposeEnum string

func (pe PurposeEnum) String() string { return string(pe) }

const (
	PurposeEnumFineTune   PurposeEnum = "fine-tune"  // 微调
	PurposeEnumAssistants PurposeEnum = "assistants" // 助理
)

// HyperparametersObject 超参数对象
type HyperparametersObject struct {
	NEpochs                interface{} `json:"n_epochs"`                           // 训练的Epoch数, 一般为int或者"auto"
	BatchSize              interface{} `json:"batch_size,omitempty"`               // 每个batch的大小
	LearningRateMultiplier interface{} `json:"learning_rate_multiplier,omitempty"` // 学习率乘数
}

// FineTuningJobObject 微调任务对象
// reference https://platform.openai.com/docs/api-reference/fine-tuning/object
type FineTuningJobObject struct {
	Object          string                `json:"object"`                    // 一般为fine-tuning.job(通过API创建的微调任务)
	ID              string                `json:"id"`                        // 微调任务的唯一标识符
	Model           string                `json:"model"`                     // 进行微调的模型
	CreatedAt       int64                 `json:"created_at"`                // 微调任务创建的时间
	FinishedAt      int64                 `json:"finished_at,omitempty"`     // 微调任务结束的时间 (Optional)
	FineTunedModel  string                `json:"fine_tuned_model"`          // 微调后的模型
	OrganizationID  string                `json:"organization_id"`           // 微调任务所属的组织
	ResultFiles     []string              `json:"result_files"`              // 微调任务的结果文件
	Status          string                `json:"status"`                    // 微调任务的状态
	ValidationFile  interface{}           `json:"validation_file,omitempty"` // 微调任务的验证文件 (Optional)
	TrainingFile    string                `json:"training_file"`             // 微调任务的训练文件
	Hyperparameters HyperparametersObject `json:"hyperparameters"`           // 微调任务的超参数
	TrainedTokens   int                   `json:"trained_tokens"`            // 微调任务的训练token数
}

// CreateFineTuningJobRequestBody 创建微调任务请求体
// reference https://platform.openai.com/docs/api-reference/fine-tuning/create
type CreateFineTuningJobRequestBody struct {
	Model           string                `json:"model"`                     // 进行微调的模型
	TrainingFile    string                `json:"training_file"`             // 微调任务的训练文件
	ValidationFile  string                `json:"validation_file,omitempty"` // 微调任务的验证文件
	Hyperparameters HyperparametersObject `json:"hyperparameters,omitempty"` // 微调任务的超参数
	Suffix          string                `json:"suffix,omitempty"`          // 微调任务的后缀 (Optional)
}

// CreateFineTuningJobRequest 创建微调任务请求
type CreateFineTuningJobRequest struct {
	Body CreateFineTuningJobRequestBody
}

// CreateFineTuningJobResponseBody 创建微调任务响应
type CreateFineTuningJobResponseBody struct {
	FineTuningJobObject
}

// RetrieveFineTuningJobRequestBody 获取微调任务请求
// reference https://platform.openai.com/docs/api-reference/fine-tuning/retrieve
type RetrieveFineTuningJobRequestBody struct {
	ID string // 微调任务的唯一标识符
}

// RetrieveFineTuningJobRequest 获取微调任务请求
type RetrieveFineTuningJobRequest struct {
	Body RetrieveFineTuningJobRequestBody
}

// RetrieveFineTuningJobResponseBody 获取微调任务响应
type RetrieveFineTuningJobResponseBody struct {
	FineTuningJobObject
}

// ListFineTuningJobsRequestBody 列出微调任务请求体
// reference https://platform.openai.com/docs/api-reference/fine-tuning/list
type ListFineTuningJobsRequestBody struct {
	After string // 上一页最后一个微调任务的标识符，用于分页
	Limit int    // 每页的微调任务数量, 默认20
}

// ListFineTuningJobsRequest 列出微调任务请求
type ListFineTuningJobsRequest struct {
	Body ListFineTuningJobsRequestBody
}

// ListFineTuningJobsResponseBody 列出微调任务响应
type ListFineTuningJobsResponseBody struct {
	Object  string                `json:"object"`   // 一般为list
	Data    []FineTuningJobObject `json:"data"`     // 微调任务列表
	HasMore bool                  `json:"has_more"` // 是否还有更多的微调任务
}

// CancelFineTuningJobRequestBody 取消微调任务请求体
// reference https://platform.openai.com/docs/api-reference/fine-tuning/cancel
type CancelFineTuningJobRequestBody struct {
	ID string // 微调任务的唯一标识符
}

// CancelFineTuningJobRequest 取消微调任务请求
type CancelFineTuningJobRequest struct {
	Body CancelFineTuningJobRequestBody
}

// CancelFineTuningJobResponseBody 取消微调任务响应
type CancelFineTuningJobResponseBody struct {
	FineTuningJobObject
}

// FileObject 文件对象
// reference https://platform.openai.com/docs/api-reference/files/object
type FileObject struct {
	ID       string `json:"id"`       // 文件的唯一标识符
	Bytes    int64  `json:"bytes"`    // 文件的大小
	Created  int64  `json:"created"`  // 文件创建的时间
	FileName string `json:"filename"` // 文件的名称
	Object   string `json:"object"`   // 文件的类型，一定为file
	Purpose  string `json:"purpose"`  // 文件的用途, 可能的值见reference
	// Status        string `json:"status"`         // 已过时。文件的状态，uploaded/processed/error
	// StatusDetails string `json:"status_details"` // 已过时。文件的状态详情
}

// UploadFileRequestBody 文件上传请求体
// reference https://platform.openai.com/docs/api-reference/files/create
type UploadFileRequestBody struct {
	File     io.Reader `json:"file"`     // 需要上传的文件
	FileName string    `json:"filename"` // 文件的名称, openai不需要此字段
	Purpose  string    `json:"purpose"`  // 文件的用途
}

func (ufr UploadFileRequestBody) ToMultiPartBody() map[string]string {
	return map[string]string{
		"purpose": ufr.Purpose,
	}
}

// UploadFileRequest 文件上传请求
type UploadFileRequest struct {
	FormBody UploadFileRequestBody
}

// UploadFileResponseBody 文件上传响应
// reference https://platform.openai.com/docs/api-reference/files/create
type UploadFileResponseBody struct {
	FileObject
}

// ListFilesRequestBody 列出文件请求体
// reference https://platform.openai.com/docs/api-reference/files/list
type ListFilesRequestBody struct {
	Purpose string `json:"purpose,omitempty"` // 文件的用途
}

// ListFilesRequest 列出文件请求
type ListFilesRequest struct {
	Body ListFilesRequestBody
}

// ListFilesResponseBody 列出文件响应
type ListFilesResponseBody struct {
	Object string       `json:"object"` // 一般为list
	Data   []FileObject `json:"data"`   // 文件列表
}

// DeleteFileRequestBody 删除文件请求体
// reference https://platform.openai.com/docs/api-reference/files/delete
type DeleteFileRequestBody struct {
	ID string // 文件的唯一标识符
}

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	Body DeleteFileRequestBody
}

// DeleteFileResponseBody 删除文件响应
type DeleteFileResponseBody struct {
	FileObject
}

// RetrieveFileRequestBody 获取文件请求体
// reference https://platform.openai.com/docs/api-reference/files/retrieve
type RetrieveFileRequestBody struct {
	ID string // 文件的唯一标识符
}

// RetrieveFileRequest 获取文件请求
type RetrieveFileRequest struct {
	Body RetrieveFileRequestBody
}

// RetrieveFileResponseBody 获取文件响应
type RetrieveFileResponseBody struct {
	FileObject
}
