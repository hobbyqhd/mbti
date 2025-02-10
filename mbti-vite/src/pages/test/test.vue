<template>
  <view class="test-container">
    <view class="progress-bar">
      <view class="progress" :style="{ width: progress + '%' }"></view>
    </view>
    
    <view class="question-card" v-if="currentQuestion">
      <text class="question-number">问题 {{ currentIndex + 1 }}/{{ totalQuestions }}</text>
      <text class="question-text">{{ currentQuestion.text }}</text>
      
      <view class="options">
        <button 
          class="option-button" 
          v-for="(option, index) in currentQuestion.options" 
          :key="index"
          @click="selectOption(option.value)"
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
      questions: [
        {
          text: '在社交场合中，你通常会：',
          options: [
            { text: '与很多人交谈，获得能量', value: 'E' },
            { text: '倾向于与少数人深入交谈', value: 'I' }
          ]
        },
        // 添加更多MBTI测试题目
      ]
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
  methods: {
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
    submitTest() {
      // 计算MBTI结果并跳转到结果页面
      const result = this.calculateResult()
      uni.navigateTo({
        url: `/pages/result/result?type=${result}`
      })
    },
    calculateResult() {
      // 简单的结果计算逻辑
      return 'ENFP' // 示例返回值
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
  transition: background-color 0.3s;
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