package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	query := "SELECT id, question, options FROM questions"
	start := time.Now()
	log.Printf("开始执行SQL查询: %s", query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("SQL查询失败: %v, 耗时: %v", err, time.Since(start))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取题目失败"})
		return
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		var optionsJSON string
		if err := rows.Scan(&q.ID, &q.Question, &optionsJSON); err != nil {
			log.Printf("解析行数据失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析题目失败"})
			return
		}
		if err := json.Unmarshal([]byte(optionsJSON), &q.Options); err != nil {
			log.Printf("解析选项JSON失败: %v, 数据: %s", err, optionsJSON)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析选项失败"})
			return
		}
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

	// 计算MBTI维度得分
	scores := CalculateScores(answer.Answers)
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

func CalculateScores(answers []int) map[string]float64 {
	// 实现MBTI得分计算逻辑
	scores := make(map[string]float64)
	// E/I, S/N, T/F, J/P 维度的题目索引（从0开始）
	dimensions := map[string][]int{
		"EI": {0, 7, 14, 21, 28, 35, 42, 49, 56, 63},
		"SN": {1, 8, 15, 22, 29, 36, 43, 50, 57, 64},
		"TF": {2, 9, 16, 23, 30, 37, 44, 51, 58, 65},
		"JP": {3, 10, 17, 24, 31, 38, 45, 52, 59, 66},
	}

	for dim, indices := range dimensions {
		var score float64
		for _, idx := range indices {
			// 添加边界检查
			if idx >= len(answers) {
				continue
			}
			if answers[idx] == 0 {
				score++
			}
		}
		scores[dim] = (score / float64(len(indices))) * 100
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
		{Left: "E", Right: "I", LeftValue: 100 - scores["EI"], RightValue: scores["EI"]},
		{Left: "S", Right: "N", LeftValue: 100 - scores["SN"], RightValue: scores["SN"]},
		{Left: "T", Right: "F", LeftValue: 100 - scores["TF"], RightValue: scores["TF"]},
		{Left: "J", Right: "P", LeftValue: 100 - scores["JP"], RightValue: scores["JP"]},
	}

	// 生成个性化报告
	report := GenerateReport(mbtiType)

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
	if scores["EI"] < 50 {
		mbtiType += "E"
	} else {
		mbtiType += "I"
	}
	if scores["SN"] < 50 {
		mbtiType += "S"
	} else {
		mbtiType += "N"
	}
	if scores["TF"] < 50 {
		mbtiType += "T"
	} else {
		mbtiType += "F"
	}
	if scores["JP"] < 50 {
		mbtiType += "J"
	} else {
		mbtiType += "P"
	}
	return mbtiType
}

func GenerateReport(mbtiType string) string {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Printf("Warning: DEEPSEEK_API_KEY not set, using default report template")
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	// 构建提示词
	prompt := fmt.Sprintf("请用专业的语气，详细描述MBTI中%s类型的性格特点，包括：\n1. 总体特征\n2. 人际关系\n3. 工作风格\n4. 个人发展建议", mbtiType)
	log.Printf("开始生成%s类型的个性化报告", mbtiType)

	// 调用DeepSeek API
	client := &http.Client{}
	reqBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{{
			"role":    "user",
			"content": prompt,
		}},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	log.Printf("正在调用DeepSeek API生成报告...")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling DeepSeek API: %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}
	defer resp.Body.Close()

	log.Printf("DeepSeek API响应状态码: %d", resp.StatusCode)

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding response: %v", err)
		return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
	}

	// 提取生成的报告内容
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				log.Printf("成功生成%s类型的个性化报告", mbtiType)
				return content
			}
		}
	}

	log.Printf("无法从API响应中提取报告内容，使用默认模板")
	return fmt.Sprintf("%s类型的性格特点是...", mbtiType)
}
