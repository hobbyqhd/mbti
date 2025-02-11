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
        <button 
          class="option-button" 
          v-for="(option, index) in currentQuestion.options" 
          :key="index"
          @click="selectOption(option.value)"
          :class="{ selected: answers[currentIndex] === option.value }"
        >
          {{ option.text }}
        </button>
      </view>
    </view>

    <view class="navigation-buttons">
      <button class="nav-button" @click="previousQuestion" :disabled="currentIndex === 0">上一题</button>
      <button 
        class="nav-button submit" 
        @click="submitTest" 
        v-if="currentIndex === totalQuestions - 1"
        :disabled="answers[currentIndex] === undefined"
      >提交测试</button>
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
        options: this.parseOptions(question.options)
      }))
    },
    parseOptions(options) {
      try {
        const parsedOptions = Array.isArray(options) ? options : JSON.parse(options)
        return parsedOptions.map((optionText, index) => ({
          text: optionText,
          value: index
        }))
      } catch (e) {
        console.error('选项解析失败:', e)
        return []
      }
    },
    initializeAnswers() {
      this.answers = new Array(this.questions.length).fill(undefined)
    },
    selectOption(value) {
      this.answers[this.currentIndex] = value
      if (this.currentIndex < this.totalQuestions - 1) {
        setTimeout(this.nextQuestion, 300)
      }
    },
    previousQuestion() {
      if (this.currentIndex > 0) {
        this.currentIndex--
      }
    },
    nextQuestion() {
      if (this.currentIndex < this.totalQuestions - 1) {
        this.currentIndex++
      }
    },
    async submitTest() {
      if (!this.validateAnswers()) return

      try {
        await this.showLoadingIndicator()
        const dimensions = this.calculateDimensions()
        const response = await this.submitAnswers(dimensions)
        await this.handleSubmitResponse(response)
      } catch (error) {
        this.handleError('提交失败', error)
      }
    },
    validateAnswers() {
      if (this.answers.some(answer => answer === undefined)) {
        uni.showToast({
          title: '请回答所有问题',
          icon: 'none',
          duration: 2000
        })
        return false
      }
      return true
    },
    showLoadingIndicator() {
      return uni.showLoading({
        title: 'AI正在深入分析您的答案\n请耐心等待（约10分钟）',
        mask: true
      })
    },
    calculateDimensions() {
      const dimensions = { EI: 0, SN: 0, TF: 0, JP: 0 }
      this.answers.forEach((answer, index) => {
        const questionType = Math.floor(index / 23)
        const score = answer === 1 ? 1 : -1
        const dimensionKeys = ['EI', 'SN', 'TF', 'JP']
        dimensions[dimensionKeys[questionType]] += score
      })
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
        console.error('页面跳转失败:', err)
        await uni.redirectTo({
          url: `/pages/result/result?id=${resultId}`
        })
      } finally {
        uni.hideLoading()
      }
    },
    handleError(message, error) {
      console.error(message + ':', error)
      uni.showToast({
        title: `${message}: ${error.message}`,
        icon: 'none',
        duration: 3000
      })
    }
  }
}
</script>

<style>
.test-container {
  padding: 30rpx;
  min-height: 100vh;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.progress-bar {
  width: 100%;
  height: 8rpx;
  background-color: rgba(239, 239, 244, 0.9);
  border-radius: 8rpx;
  margin-bottom: 30rpx;
  overflow: hidden;
}

.progress {
  height: 100%;
  background-color: #007AFF;
  border-radius: 8rpx;
  transition: width 0.3s ease;
}

.loading {
  text-align: center;
  padding: 40px;
}

.loading text {
  font-size: 16px;
  color: #666;
}

.question-card {
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  margin-bottom: 30rpx;
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
}

.question-number {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
  display: block;
}

.question-text {
  font-size: 18px;
  color: #333;
  margin-bottom: 20px;
  display: block;
}

.options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.option-button {
  background-color: rgba(239, 239, 244, 0.9);
  border: none;
  padding: 24rpx 30rpx;
  border-radius: 12rpx;
  text-align: left;
  font-size: 28rpx;
  transition: all 0.2s ease;
  box-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.05);
  width: 100%;
  min-width: 500rpx;
  max-width: 600rpx;
  margin: 0 auto;
}

.option-button.selected {
  background-color: #2196F3;
  color: white;
}

.option-button:active {
  background-color: #e0e0e0;
}

.navigation-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
  position: fixed;
  bottom: 40rpx;
  left: 30rpx;
  right: 30rpx;
  z-index: 100;
}

.nav-button {
  padding: 20rpx 40rpx;
  border-radius: 12rpx;
  background-color: #007AFF;
  color: white;
  border: none;
  font-size: 28rpx;
  box-shadow: 0 2rpx 6rpx rgba(0, 122, 255, 0.2);
  transition: all 0.2s ease;
}

.nav-button:disabled {
  background-color: #ccc;
}

.submit {
  background-color: #4CAF50;
}
</style>