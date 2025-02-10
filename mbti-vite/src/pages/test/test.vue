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
      <button class="nav-button" @click="nextQuestion" v-if="currentIndex < totalQuestions - 1">下一题</button>
      <button class="nav-button submit" @click="submitTest" v-else>提交测试</button>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      currentIndex: 0,
      answers: [],
      questions: [],
      loading: true
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
          this.questions = response.data.map(question => {
            let options = [];
            try {
              options = Array.isArray(question.options) ? question.options : JSON.parse(question.options);
            } catch (e) {
              console.error('选项解析失败:', e);
              options = [];
            }
            return {
              text: question.question || '',
              options: options.map((optionText, index) => ({
                text: optionText,
                value: index
              }))
            }
          })
        } else {
          throw new Error(`获取题目失败: ${response.statusCode}`)
        }
      } catch (error) {
        console.error('请求失败:', error);
        uni.showToast({
          title: `获取题目失败: ${error.message}`,
          icon: 'none',
          duration: 3000
        })
      } finally {
        this.loading = false
      }
    },
    selectOption(value) {
      this.answers[this.currentIndex] = value
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
      if (this.answers.length < this.totalQuestions) {
        uni.showToast({
          title: '请回答所有问题',
          icon: 'none'
        })
        return
      }

      try {
        const response = await uni.request({
          url: 'http://localhost:8080/api/submit',
          method: 'POST',
          data: {
            answers: this.answers
          }
        })

        if (response.statusCode === 200 && response.data.resultId) {
          uni.navigateTo({
            url: `/pages/result/result?id=${response.data.resultId}`
          })
        } else {
          throw new Error('提交失败')
        }
      } catch (error) {
        uni.showToast({
          title: '提交失败，请重试',
          icon: 'none'
        })
      }
    }
  }
}
</script>

<style>
.test-container {
  padding: 20px;
}

.progress-bar {
  width: 100%;
  height: 10px;
  background-color: #eee;
  border-radius: 5px;
  margin-bottom: 20px;
}

.progress {
  height: 100%;
  background-color: #4CAF50;
  border-radius: 5px;
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
  background-color: #fff;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
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
  background-color: #f5f5f5;
  border: none;
  padding: 15px;
  border-radius: 5px;
  text-align: left;
  font-size: 16px;
  transition: all 0.3s ease;
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
}

.nav-button {
  padding: 10px 20px;
  border-radius: 5px;
  background-color: #2196F3;
  color: white;
  border: none;
}

.nav-button:disabled {
  background-color: #ccc;
}

.submit {
  background-color: #4CAF50;
}
</style>