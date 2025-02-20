package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type Answer struct {
	Answers []int `json:"answers"`
}

type Result struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Dimensions []Dimension `json:"dimensions"`
	Report     string      `json:"report"`
}

type Dimension struct {
	Left       string  `json:"left"`
	Right      string  `json:"right"`
	LeftValue  float64 `json:"leftValue"`
	RightValue float64 `json:"rightValue"`
}

func GetQuestions(c *gin.Context) {
	testType := c.Query("type")
	if testType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "测试类型不能为空"})
		return
	}

	var types []string
	switch testType {
	case "simple":
		types = []string{"simple"}
	case "detailed":
		types = []string{"simple", "detailed"}
	case "full":
		types = []string{"simple", "detailed", "full"}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的测试类型"})
		return
	}

	query := "SELECT id, question, dimension, direction FROM questions WHERE type IN (?)"
	start := time.Now()
	log.Printf("开始执行SQL查询: %s, 参数: %s", query, testType)
	query, args, err := sqlx.In(query, types)
	if err != nil {
		log.Printf("构建IN查询失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取题目失败"})
		return
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Printf("SQL查询失败: %v, 耗时: %v", err, time.Since(start))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取题目失败"})
		return
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		var dimension string
		var direction int
		if err := rows.Scan(&q.ID, &q.Question, &dimension, &direction); err != nil {
			log.Printf("解析行数据失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析题目失败"})
			return
		}
		// 设置7级量表选项
		q.Options = []string{"非常不符合", "不符合", "有点不符合", "中立", "有点符合", "符合", "非常符合"}
		questions = append(questions, q)
	}

	log.Printf("SQL查询完成, 获取到%d条记录, 总耗时: %v", len(questions), time.Since(start))
	c.JSON(http.StatusOK, questions)
}

func SubmitAnswers(c *gin.Context) {
	start := time.Now()
	log.Printf("开始处理答案提交请求")

	var answer Answer
	if err := c.BindJSON(&answer); err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的答案数据"})
		return
	}

	// 获取测试类型
	testType := c.Query("type")
	if testType == "" {
		testType = "simple" // 默认使用simple类型
	}

	// 计算MBTI维度得分
	scores := CalculateScores(answer.Answers, testType)
	// 生成结果ID
	resultID := GenerateResultID()

	// 保存结果到数据库
	if err := SaveResult(resultID, scores); err != nil {
		log.Printf("保存结果失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存结果失败"})
		return
	}

	log.Printf("答案提交处理完成, 结果ID: %s, 总耗时: %v", resultID, time.Since(start))
	c.JSON(http.StatusOK, gin.H{"resultId": resultID})
}

func GetResult(c *gin.Context) {
	start := time.Now()
	resultID := c.Param("id")
	log.Printf("开始获取结果, ID: %s", resultID)

	// 从数据库获取结果
	var result Result
	var dimensionsJSON string
	query := "SELECT id, type, dimensions, report FROM results WHERE id = ?"
	log.Printf("执行SQL: %s, 参数: [%s]", query, resultID)

	err := DB.QueryRow(query, resultID).Scan(
		&result.ID, &result.Type, &dimensionsJSON, &result.Report)
	if err != nil {
		log.Printf("获取结果失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "结果不存在"})
		return
	}

	// 解析维度数据
	if err := json.Unmarshal([]byte(dimensionsJSON), &result.Dimensions); err != nil {
		log.Printf("解析维度数据失败: %v, 数据: %s", err, dimensionsJSON)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析结果失败"})
		return
	}

	log.Printf("获取结果完成, ID: %s, 类型: %s, 总耗时: %v", resultID, result.Type, time.Since(start))
	c.JSON(http.StatusOK, result)
}

func CalculateScores(answers []int, testType string) map[string]float64 {
	// 实现MBTI得分计算逻辑
	scores := make(map[string]float64)

	// 根据测试类型确定题目范围
	var types []string
	switch testType {
	case "simple":
		types = []string{"simple"}
	case "detailed":
		types = []string{"simple", "detailed"}
	case "full":
		types = []string{"simple", "detailed", "full"}
	default:
		log.Printf("无效的测试类型: %s", testType)
		return scores
	}

	// 构建SQL查询
	query := "SELECT dimension, direction FROM questions WHERE type IN (?) ORDER BY id"
	query, args, err := sqlx.In(query, types)
	if err != nil {
		log.Printf("构建IN查询失败: %v", err)
		return scores
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Printf("获取题目信息失败: %v", err)
		return scores
	}
	defer rows.Close()

	// 存储每个维度的题目信息
	dimensionScores := make(map[string]struct {
		total float64
		count int
	})

	// 处理每个题目的答案
	for i := 0; i < len(answers) && rows.Next(); i++ {
		var dimension string
		var direction int
		if err := rows.Scan(&dimension, &direction); err != nil {
			log.Printf("解析题目信息失败: %v", err)
			continue
		}

		// 确保答案在1-7的范围内
		answer := answers[i]
		if answer < 1 {
			answer = 1
		} else if answer > 7 {
			answer = 7
		}

		// 计算得分（考虑题目方向）
		// 对于正向题目(direction=1)：1->0%, 7->100%
		// 对于反向题目(direction=-1)：1->100%, 7->0%
		var score float64
		if direction == 1 {
			score = float64(answer-1) * (100.0 / 6.0) // 正向题目：将1-7线性映射到0-100
		} else {
			score = float64(7-answer) * (100.0 / 6.0) // 反向题目：将7-1线性映射到0-100
		}

		// 确保得分在0-100范围内
		if score < 0 {
			score = 0
		} else if score > 100 {
			score = 100
		}

		// 累加分数
		scoreInfo := dimensionScores[dimension]
		scoreInfo.total += score
		scoreInfo.count++
		dimensionScores[dimension] = scoreInfo
	}

	// 计算每个维度的最终得分
	for dim, info := range dimensionScores {
		if info.count > 0 {
			// 计算平均分并确保在0-100范围内
			avgScore := info.total / float64(info.count)
			if avgScore < 0 {
				avgScore = 0
			} else if avgScore > 100 {
				avgScore = 100
			}
			scores[dim] = math.Round(avgScore)
		}
	}

	return scores
}

func GenerateResultID() string {
	return fmt.Sprintf("result_%d", time.Now().UnixNano())
}

func SaveResult(resultID string, scores map[string]float64) error {
	start := time.Now()
	log.Printf("开始保存结果, ID: %s", resultID)

	// 根据得分确定MBTI类型
	mbtiType := DetermineMBTIType(scores)
	log.Printf("计算得到MBTI类型: %s, 得分: %v", mbtiType, scores)

	// 生成维度数据
	dimensions := []Dimension{
		{Left: "E", Right: "I", LeftValue: math.Round(100 - scores["EI"]), RightValue: math.Round(scores["EI"])},
		{Left: "S", Right: "N", LeftValue: math.Round(100 - scores["SN"]), RightValue: math.Round(scores["SN"])},
		{Left: "T", Right: "F", LeftValue: math.Round(100 - scores["TF"]), RightValue: math.Round(scores["TF"])},
		{Left: "J", Right: "P", LeftValue: math.Round(100 - scores["JP"]), RightValue: math.Round(scores["JP"])},
	}

	// 生成个性化报告
	report := GenerateReport(mbtiType, dimensions)

	// 将维度数据转换为JSON
	dimensionsJSON, err := json.Marshal(dimensions)
	if err != nil {
		log.Printf("维度数据JSON序列化失败: %v", err)
		return err
	}

	// 保存到数据库
	query := "INSERT INTO results (id, type, dimensions, report) VALUES (?, ?, ?, ?)"
	log.Printf("执行SQL: %s, 参数: [%s, %s, %s, %s]", query, resultID, mbtiType, string(dimensionsJSON), report)

	result, err := DB.Exec(query, resultID, mbtiType, string(dimensionsJSON), report)
	if err != nil {
		log.Printf("数据库插入失败: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("结果保存完成, 影响行数: %d, 总耗时: %v", rowsAffected, time.Since(start))
	return nil
}

func DetermineMBTIType(scores map[string]float64) string {
	var mbtiType string
	// 当分数等于50时，默认选择第二个类型
	if scores["EI"] >= 50 {
		mbtiType += "E"
	} else {
		mbtiType += "I"
	}
	if scores["SN"] >= 50 {
		mbtiType += "S"
	} else {
		mbtiType += "N"
	}
	if scores["TF"] >= 50 {
		mbtiType += "T"
	} else {
		mbtiType += "F"
	}
	if scores["JP"] >= 50 {
		mbtiType += "J"
	} else {
		mbtiType += "P"
	}
	return mbtiType
}

func GenerateReport(mbtiType string, dimensions []Dimension) string {

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Printf("[MBTI报告生成] 错误: DEEPSEEK_API_KEY未设置，将使用默认模板")
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	// 构建维度倾向描述
	dimensionDescriptions := make([]string, len(dimensions))
	for i, dim := range dimensions {
		var tendency string
		if dim.LeftValue > dim.RightValue {
			tendency = fmt.Sprintf("偏向%s (%.0f%%)", dim.Right, dim.LeftValue)
		} else {
			tendency = fmt.Sprintf("偏向%s (%.0f%%)", dim.Left, dim.RightValue)
		}
		dimensionDescriptions[i] = fmt.Sprintf("%s vs %s: %s", dim.Left, dim.Right, tendency)
	}

	// 构建提示词
	prompt := fmt.Sprintf(`请以专业的心理学视角，对MBTI中的%s类型进行全面分析，并使用HTML标签组织内容。该用户在各维度上的具体倾向如下：

%s

请按以下HTML结构组织内容：

<div class="mbti-report">
  <section class="core-traits">
    <h2>核心特质概述</h2>
    <ul>
      <li><strong>主要认知功能：</strong>[内容]</li>
      <li><strong>思维方式特点：</strong>[内容]</li>
      <li><strong>行为模式倾向：</strong>[内容]</li>
    </ul>
  </section>

  <section class="relationships">
    <h2>人际关系分析</h2>
    <ul>
      <li><strong>沟通风格：</strong>[内容]</li>
      <li><strong>与他人互动方式：</strong>[内容]</li>
      <li><strong>在团队中的角色：</strong>[内容]</li>
      <li><strong>理想的社交环境：</strong>[内容]</li>
    </ul>
  </section>

  <section class="career">
    <h2>职业发展洞察</h2>
    <ul>
      <li><strong>最适合的工作环境：</strong>[内容]</li>
      <li><strong>职业优势：</strong>[内容]</li>
      <li><strong>潜在的职业挑战：</strong>[内容]</li>
      <li><strong>理想的职业方向：</strong>[内容]</li>
    </ul>
  </section>

  <section class="growth">
    <h2>个人成长建议</h2>
    <ul>
      <li><strong>需要培养的能力：</strong>[内容]</li>
      <li><strong>潜在的发展盲点：</strong>[内容]</li>
      <li><strong>压力管理方式：</strong>[内容]</li>
      <li><strong>自我提升方向：</strong>[内容]</li>
    </ul>
  </section>

  <section class="strengths-weaknesses">
    <h2>独特优势和劣势</h2>
    <ul>
      <li><strong>核心优势：</strong>[内容]</li>
      <li><strong>潜在劣势：</strong>[内容]</li>
      <li><strong>如何发挥优势：</strong>[内容]</li>
      <li><strong>如何克服劣势：</strong>[内容]</li>
    </ul>
  </section>

  <section class="love-relationships">
    <h2>爱情关系分析</h2>
    <ul>
      <li><strong>恋爱中的表现：</strong>[内容]</li>
      <li><strong>理想伴侣特质：</strong>[内容]</li>
      <li><strong>感情相处方式：</strong>[内容]</li>
      <li><strong>维系感情建议：</strong>[内容]</li>
    </ul>
  </section>

  <section class="friendships">
    <h2>友情互动特点</h2>
    <ul>
      <li><strong>交友方式：</strong>[内容]</li>
      <li><strong>理想友谊类型：</strong>[内容]</li>
      <li><strong>友情维护特点：</strong>[内容]</li>
      <li><strong>社交圈建议：</strong>[内容]</li>
    </ul>
  </section>

  <section class="parenting">
    <h2>育儿方式建议</h2>
    <ul>
      <li><strong>教育理念：</strong>[内容]</li>
      <li><strong>亲子互动风格：</strong>[内容]</li>
      <li><strong>可能的教育挑战：</strong>[内容]</li>
      <li><strong>育儿建议：</strong>[内容]</li>
    </ul>
  </section>

  <section class="work-habits">
    <h2>工作习惯分析</h2>
    <ul>
      <li><strong>工作节奏偏好：</strong>[内容]</li>
      <li><strong>任务管理方式：</strong>[内容]</li>
      <li><strong>工作环境需求：</strong>[内容]</li>
      <li><strong>效率提升建议：</strong>[内容]</li>
    </ul>
  </section>
</div>

请根据用户在各维度上的具体倾向程度，用专业但易懂的语言填充上述HTML结构中的[内容]部分。特别注意根据维度分数的强弱来调整建议的针对性。请确保生成的HTML格式正确，可以直接在网页中显示。`, mbtiType, strings.Join(dimensionDescriptions, "\n"))

	// 调用DeepSeek API
	client := &http.Client{
		Timeout: 5 * time.Minute, // 增加超时时间到5分钟
	}
	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-chat-1.0"
	}
	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{{
			"role":    "user",
			"content": prompt,
		}},
		"temperature":       0.5,
		"max_tokens":        2500,
		"stream":            false,
		"presence_penalty":  0.1,
		"frequency_penalty": 0.1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[MBTI报告生成] 错误: 请求序列化失败 - %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	apiURL := os.Getenv("DEEPSEEK_API_URL")
	if apiURL == "" {
		apiURL = "https://api.deepseek.com/v1/chat/completions"
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[MBTI报告生成] 错误: 创建HTTP请求失败 - %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[MBTI报告生成] 错误: 调用DeepSeek API失败 - %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}
	defer resp.Body.Close()

	// 读取响应体内容用于记录
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[MBTI报告生成] 错误: 读取响应内容失败 - %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("[MBTI报告生成] 错误: API返回非200状态码 - %d", resp.StatusCode)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	// 重新创建一个新的Reader用于JSON解码
	var result map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&result); err != nil {
		log.Printf("[MBTI报告生成] 错误: 解析API响应失败 - %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	// 提取生成的报告内容
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				// 清理内容，只保留HTML部分
				content = strings.TrimSpace(content)
				// 如果内容以```html开头且以```结尾，移除这些标记
				if strings.HasPrefix(content, "```html") {
					content = strings.TrimPrefix(content, "```html")
				}
				if strings.HasSuffix(content, "```") {
					content = strings.TrimRight(content, "```")
				}
				content = strings.TrimSpace(content)
				return content
			}
		}
	}

	log.Printf("[MBTI报告生成] 错误: 无法从API响应中提取报告内容，将使用默认模板")
	return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
}
