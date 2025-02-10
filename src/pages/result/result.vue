<template>
  <view class="result-container">
    <view class="result-card">
      <view class="result-header">
        <text class="result-title">您的MBTI性格类型是</text>
        <text class="result-type">{{ result.type }}</text>
      </view>

      <view class="result-details">
        <view class="dimension-bar" v-for="(dimension, index) in result.dimensions" :key="index">
          <text class="dimension-label">{{ dimension.left }}</text>
          <view class="dimension-progress">
            <view 
              class="dimension-value" 
              :style="{ width: dimension.leftValue + '%', backgroundColor: '#4CAF50' }"
            ></view>
            <view 
              class="dimension-value" 
              :style="{ width: dimension.rightValue + '%', backgroundColor: '#2196F3' }"
            ></view>
          </view>
          <text class="dimension-label">{{ dimension.right }}</text>
        </view>
      </view>

      <view class="report-section">
        <text class="section-title">个性化分析报告</text>
        <text class="report-content">{{ result.report }}</text>
      </view>
    </view>

    <button class="share-button" @click="shareResult">分享结果</button>
  </view>
</template>

<script>
export default {
  data() {
    return {
      result: {
        type: '',
        dimensions: [],
        report: ''
      }
    };
  },
  onLoad(options) {
    if (options.id) {
      this.fetchResult(options.id);
    }
  },
  methods: {
    async fetchResult(id) {
      try {
        const response = await uni.request({
          url: `http://localhost:8080/api/results/${id}`,
          method: 'GET'
        });
        if (response.statusCode === 200) {
          this.result = response.data;
        } else {
          uni.showToast({
            title: '获取结果失败',
            icon: 'none'
          });
        }
      } catch (error) {
        uni.showToast({
          title: '网络错误',
          icon: 'none'
        });
      }
    },
    shareResult() {
      uni.share({
        provider: 'weixin',
        scene: 'WXSceneSession',
        type: 0,
        title: 'MBTI性格测试结果',
        summary: `我的MBTI类型是${this.result.type}，快来测测你的类型吧！`,
        imageUrl: '/static/share-image.png',
        success: function () {
          uni.showToast({
            title: '分享成功',
            icon: 'success'
          });
        },
        fail: function () {
          uni.showToast({
            title: '分享失败',
            icon: 'none'
          });
        }
      });
    }
  }
};
</script>

<style>
.result-container {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: 100vh;
}

.result-card {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.result-header {
  text-align: center;
  margin-bottom: 30px;
}

.result-title {
  font-size: 18px;
  color: #666;
  margin-bottom: 10px;
}

.result-type {
  font-size: 36px;
  font-weight: bold;
  color: #333;
}

.result-details {
  margin-bottom: 30px;
}

.dimension-bar {
  margin-bottom: 20px;
}

.dimension-label {
  font-size: 14px;
  color: #666;
  width: 30px;
  text-align: center;
}

.dimension-progress {
  flex: 1;
  height: 20px;
  background-color: #f0f0f0;
  border-radius: 10px;
  overflow: hidden;
  margin: 0 10px;
  display: flex;
}

.dimension-value {
  height: 100%;
  transition: width 0.3s ease;
}

.report-section {
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 15px;
}

.report-content {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}

.share-button {
  margin-top: 30px;
  background-color: #4CAF50;
  color: #fff;
  border: none;
  border-radius: 25px;
  padding: 12px 0;
  font-size: 16px;
  width: 100%;
}

.share-button:active {
  background-color: #45a049;
}
</style>