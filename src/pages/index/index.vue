<template>
  <view class="test-container">
    <view class="progress-bar">
      <progress :percent="progressPercent" active-color="#007AFF" />
      <text class="progress-text">{{ currentQuestion + 1 }}/93</text>
    </view>

    <view class="question-card">
      <text class="question-text">{{ currentQuestionData.question }}</text>
      
      <view class="options-container">
        <button
          v-for="(option, index) in currentQuestionData.options"
          :key="index"
          class="option-button"
          :class="{ selected: selectedOption === index }"
          @click="selectOption(index)"
        >
          {{ option }}
        </button>
      </view>
    </view>

    <view class="navigation-buttons">
      <button
        class="nav-button"
        :disabled="currentQuestion === 0"
        @click="previousQuestion"
      >
        上一题
      </button>
      <button
        class="nav-button primary"
        :disabled="selectedOption === null"
        @click="nextQuestion"
      >
        {{ isLastQuestion ? '提交测试' : '下一题' }}
      </button>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      currentQuestion: 0,
      questions: [],
      answers: new Array(93).fill(null),
      selectedOption: null
    };
  },
  computed: {
    progressPercent() {
      return ((this.currentQuestion + 1) / 93) * 100;
    },
    currentQuestionData() {
      return this.questions[this.currentQuestion] || {
        question: '加载中...',
        options: []
      };
    },
    isLastQuestion() {
      return this.currentQuestion === 92;
    }
  },
  onLoad() {
    this.fetchQuestions();
  },
  methods: {
    async fetchQuestions() {
      try {
        const response = await uni.request({
          url: 'http://localhost:8080/api/questions',
          method: 'GET'
        });
        this.questions = response.data;
      } catch (error) {
        uni.showToast({
          title: '获取题目失败',
          icon: 'none'
        });
      }
    },
    selectOption(index) {
      this.selectedOption = index;
      this.answers[this.currentQuestion] = index;
    },
    previousQuestion() {
      if (this.currentQuestion > 0) {
        this.currentQuestion--;
        this.selectedOption = this.answers[this.currentQuestion];
      }
    },
    async nextQuestion() {
      if (this.isLastQuestion) {
        await this.submitTest();
      } else {
        this.currentQuestion++;
        this.selectedOption = this.answers[this.currentQuestion];
      }
    },
    async submitTest() {
      try {
        const response = await uni.request({
          url: 'http://localhost:8080/api/submit',
          method: 'POST',
          data: {
            answers: this.answers
          }
        });
        
        uni.navigateTo({
          url: `/pages/result/result?id=${response.data.resultId}`
        });
      } catch (error) {
        uni.showToast({
          title: '提交失败',
          icon: 'none'
        });
      }
    }
  }
};
</script>

<style>
.test-container {
  padding: 20px;
}

.progress-bar {
  margin-bottom: 30px;
}

.progress-text {
  text-align: center;
  display: block;
  margin-top: 10px;
  color: #666;
}

.question-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
}

.question-text {
  font-size: 18px;
  color: #333;
  margin-bottom: 20px;
  display: block;
  line-height: 1.5;
}

.options-container {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.option-button {
  background: #f5f5f5;
  border: none;
  padding: 15px 20px;
  border-radius: 8px;
  text-align: left;
  font-size: 16px;
  color: #333;
}

.option-button.selected {
  background: #007AFF;
  color: #fff;
}

.navigation-buttons {
  display: flex;
  justify-content: space-between;
  gap: 15px;
}

.nav-button {
  flex: 1;
  padding: 12px;
  border-radius: 8px;
  font-size: 16px;
}

.nav-button.primary {
  background: #007AFF;
  color: #fff;
}

.nav-button[disabled] {
  opacity: 0.5;
}
</style>