<template>
  <view class="test-container">
    <view class="progress-bar">
      <view class="progress" :style="{ width: progress + '%' }"></view>
    </view>
    
    <view v-if="loading" class="loading">
      <text>加载中...</text>
    </view>

    <view class="question-card" v-else-if="currentQuestion">
      <text class="question-number">问题 {{ currentIndex + 1 }}/{{ totalQuestions }}</text>
      <text class="question-text">{{ currentQuestion.text }}</text>
      
      <view class="options">
        <view class="option-labels">
          <text class="extreme-label left">非常不符合</text>
          <text class="extreme-label right">非常符合</text>
        </view>
        <view class="option-circles">
          <view 
            class="option-circle" 
            v-for="value in 7" 
            :key="value"
            @click="selectOption(value)"
            :class="{ selected: answers[currentIndex] === value }"
            :style="{ transform: `scale(${0.8 + value * 0.05})` }"
          >
            <text class="option-value">{{ value }}</text>
          </view>
        </view>
        <view class="option-description">
          <text>{{ getOptionLabel(answers[currentIndex]) }}</text>
        </view>
      </view>
    </view>

    <view class="navigation-buttons">
      <button class="nav-button" @click="previousQuestion" :disabled="currentIndex === 0">上一题</button>
      <button 
        v-if="currentIndex === totalQuestions - 1"
        class="nav-button submit" 
        @click="submitTest" 
        :disabled="answers[currentIndex] === undefined"
      >提交测试</button>
      <button 
        v-else
        class="nav-button next" 
        @click="nextQuestion" 
        :disabled="answers[currentIndex] === undefined"
      >下一题</button>
    </view>
  </view>
</template>

<script>
import './test.css';

export default {
  data() {
    return {
      currentIndex: 0,
      answers: [],
      questions: [],
      loading: true,
      _testStartTime: Date.now()
    }
  },
  computed: {
    totalQuestions() {
      return this.questions.length
    },
    currentQuestion() {
      return this.questions[this.currentIndex]
    },
    progress() {
      return (this.currentIndex / this.totalQuestions) * 100
    }
  },
  onLoad() {
    this.fetchQuestions()
  },
  methods: {
    async fetchQuestions() {
      try {
        const response = await uni.request({
          url: 'http://localhost:8080/api/questions',
          method: 'GET'
        })

        if (response.statusCode === 200 && Array.isArray(response.data)) {
          this.questions = this.formatQuestions(response.data)
          this.initializeAnswers()
        } else {
          throw new Error(`获取题目失败: ${response.statusCode}`)
        }
      } catch (error) {
        this.handleError('获取题目失败', error)
      } finally {
        this.loading = false
      }
    },
    formatQuestions(data) {
      return data.map(question => ({
        text: question.question || '',
        dimension: question.dimension,
        direction: question.direction
      }))
    },
    getOptionLabel(value) {
      switch(value) {
        case 1: return '非常不符合'
        case 2: return '较为不符合'
        case 3: return '稍微不符合'
        case 4: return '中立'
        case 5: return '稍微符合'
        case 6: return '较为符合'
        case 7: return '非常符合'
        default: return ''
      }
    },
    initializeAnswers() {
      this.answers = new Array(this.totalQuestions).fill(undefined)
    },
    selectOption(value) {
      this.answers[this.currentIndex] = value
      // 自动跳转到下一题，但不在最后一题时跳转
      if (this.currentIndex < this.totalQuestions - 1) {
        setTimeout(() => {
          this.nextQuestion()
        }, 500)
      }
    },
    nextQuestion() {
      if (this.currentIndex < this.totalQuestions - 1) {
        this.currentIndex++
      }
    },
    previousQuestion() {
      if (this.currentIndex > 0) {
        this.currentIndex--
      }
    },
    async submitTest() {
      try {
        await this.showLoadingDialog()
        const dimensions = this.calculateDimensions()
        const response = await this.submitAnswers(dimensions)
        await this.handleSubmitResponse(response)
      } catch (error) {
        this.handleError('提交失败', error)
      } finally {
        uni.hideLoading()
      }
    },
    showLoadingDialog() {
      return uni.showLoading({
        title: 'AI正在深入分析您的答案\n请耐心等待（约10分钟）',
        mask: true
      })
    },
    calculateDimensions() {
      const dimensions = { EI: 0, SN: 0, TF: 0, JP: 0 }
      const dimensionCounts = { EI: 0, SN: 0, TF: 0, JP: 0 }
      
      // 计算每个维度的得分
      this.answers.forEach((answer, index) => {
        if (answer === undefined) return
        
        const question = this.questions[index]
        const dimensionKey = question.dimension
        const direction = question.direction
        
        let score = answer
        if (direction < 0) {
          score = 8 - score // 反向题目：7->1, 6->2, 5->3, 4->4, 3->5, 2->6, 1->7
        } else {
          score = score // 正向题目保持原始分数
        }
        score = ((score - 1) / 6.0) * 100 // 将1-7转换为0-100范围的百分比
        
        dimensions[dimensionKey] += score
        dimensionCounts[dimensionKey]++
      })

      // 计算每个维度的最终百分比
      for (const key in dimensions) {
        const totalQuestions = dimensionCounts[key]
        if (totalQuestions === 0) continue
        
        const averageScore = dimensions[key] / totalQuestions
        
        dimensions[key] = {
          leftValue: Math.round(averageScore),
          rightValue: Math.round(100 - averageScore)
        }
      }

      return dimensions
    },
    async submitAnswers(dimensions) {
      return await uni.request({
        url: 'http://localhost:8080/api/submit',
        method: 'POST',
        timeout: 600000,
        data: {
          answers: this.answers,
          dimensions,
          userInfo: {
            testDate: new Date().toISOString(),
            totalQuestions: this.totalQuestions,
            completionTime: Date.now() - this._testStartTime
          }
        },
        header: {
          'Content-Type': 'application/json'
        }
      })
    },
    async handleSubmitResponse(response) {
      if (response.statusCode === 200 && response.data?.resultId) {
        await this.navigateToResult(response.data.resultId)
      } else {
        throw new Error('提交失败，未获取到结果ID')
      }
    },
    async navigateToResult(resultId) {
      try {
        await uni.navigateTo({
          url: `/pages/result/result?id=${resultId}`
        })
      } catch (err) {
        this.handleError('跳转结果页面失败', err)
      }
    },
    handleError(message, error) {
      console.error(message, error)
      uni.showToast({
        title: message,
        icon: 'none',
        duration: 2000
      })
    }
  }
}
</script>

<style>
.test-container {
  padding: 20px;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.progress-bar {
  height: 4px;
  background-color: #e0e0e0;
  border-radius: 2px;
  margin-bottom: 20px;
}

.progress {
  height: 100%;
  background-color: #4CAF50;
  border-radius: 2px;
  transition: width 0.3s ease;
}

.question-card {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.question-number {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.question-text {
  font-size: 18px;
  color: #333;
  margin-bottom: 20px;
  line-height: 1.5;
}

.options {
  display: flex;
  flex-direction: column;
  margin: 20px 0;
  padding: 0 20px;
}

.option-labels {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
}

.extreme-label {
  font-size: 14px;
  color: #666;
}

.option-circles {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e3e6e8 100%);
  border-radius: 12px;
  margin: 10px 0;
  gap: 10px;
}

.option-circle {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
  color: white;
}

.option-circle:nth-child(1) { background: #FFD700; }
.option-circle:nth-child(2) { background: #FFB347; }
.option-circle:nth-child(3) { background: #FF8C00; }
.option-circle:nth-child(4) { background: #9370DB; }
.option-circle:nth-child(5) { background: #8A2BE2; }
.option-circle:nth-child(6) { background: #800080; }
.option-circle:nth-child(7) { background: #4B0082; }

.option-circle:hover {
  transform: scale(1.15);
  box-shadow: 0 6px 12px rgba(0,0,0,0.15);
}

.option-circle.selected {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  filter: brightness(1.2);
}

.option-value {
  font-size: 14px;
  font-weight: 500;
}

.option-description {
  text-align: center;
  margin-top: 10px;
  min-height: 20px;
}

.option-description text {
  color: #666;
  font-size: 14px;
}

.navigation-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: 40px;
  position: relative;
  padding: 0 20px;
  margin-bottom: 30px;
  gap: 20px;
}

.nav-button {
  flex: 1;
  padding: 12px 20px;
  border-radius: 4px;
  border: none;
  background-color: #4CAF50;
  color: white;
  font-size: 16px;
  cursor: pointer;
  min-width: 100px;
}

.nav-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.nav-button.submit {
  background-color: #2196F3;
}

.loading {
  text-align: center;
  padding: 20px;
}
</style>