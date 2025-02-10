<template>
  <view class="result-container">
    <view v-if="loading" class="loading">
      <text>加载中...</text>
    </view>

    <view v-else class="result-content">
      <view class="result-header">
        <text class="result-title">您的MBTI类型是</text>
        <text class="mbti-type">{{ result.mbtiType }}</text>
      </view>

      <view class="report-section">
        <rich-text :nodes="formattedReport"></rich-text>
      </view>

      <button class="share-button" @click="shareResult">分享结果</button>
      <button class="restart-button" @click="restartTest">重新测试</button>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      loading: true,
      result: {
        mbtiType: '',
        report: ''
      }
    }
  },
  computed: {
    formattedReport() {
      return this.result.report.replace(/\n/g, '<br>')
    }
  },
  onLoad(options) {
    if (options.id) {
      this.fetchResult(options.id)
    }
  },
  methods: {
    async fetchResult(resultId) {
      try {
        const response = await uni.request({
          url: `http://localhost:8080/api/result/${resultId}`,
          method: 'GET'
        })

        if (response.statusCode === 200) {
          this.result = response.data
        } else {
          throw new Error('获取结果失败')
        }
      } catch (error) {
        uni.showToast({
          title: '获取结果失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
      }
    },
    shareResult() {
      // 实现分享功能
      uni.showToast({
        title: '分享功能开发中',
        icon: 'none'
      })
    },
    restartTest() {
      uni.redirectTo({
        url: '/pages/test/test'
      })
    }
  }
}
</script>

<style>
.result-container {
  padding: 20px;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.result-header {
  text-align: center;
  margin-bottom: 30px;
}

.result-title {
  font-size: 20px;
  color: #333;
  margin-bottom: 10px;
  display: block;
}

.mbti-type {
  font-size: 36px;
  font-weight: bold;
  color: #2196F3;
  display: block;
}

.report-section {
  background: #fff;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  line-height: 1.6;
}

.share-button,
.restart-button {
  width: 100%;
  margin-bottom: 15px;
  padding: 12px;
  border-radius: 5px;
  font-size: 16px;
}

.share-button {
  background-color: #2196F3;
  color: white;
}

.restart-button {
  background-color: #f5f5f5;
  color: #333;
}
</style>