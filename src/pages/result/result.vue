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
    shareResult() {
      uni.share({
        provider: "weixin",
        scene: "WXSceneSession",
        type: 0,
        title: `我的MBTI性格类型是${this.result.type}`,
        summary: this.result.report,
        success: function (res) {
          console.log("success:" + JSON.stringify(res));
        },
        fail: function (err) {
          console.log("fail:" + JSON.stringify(err));
        }
      });
    }
  }
};
</script>

<style>
.result-container {
  padding: 30rpx;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.result-card {
  background: #ffffff;
  border-radius: 20rpx;
  padding: 40rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
  margin-bottom: 30rpx;
  animation: slideIn 0.5s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.result-header {
  text-align: center;
  margin-bottom: 50rpx;
}

.result-title {
  font-size: 32rpx;
  color: #666;
  margin-bottom: 20rpx;
  display: block;
}

.result-type {
  font-size: 64rpx;
  font-weight: bold;
  color: #333;
  display: block;
  letter-spacing: 4rpx;
  text-shadow: 2rpx 2rpx 4rpx rgba(0, 0, 0, 0.1);
}

.dimension-bar {
  margin-bottom: 30rpx;
}

.dimension-label {
  font-size: 28rpx;
  color: #666;
  margin: 0 20rpx;
}

.dimension-progress {
  height: 20rpx;
  background: #f0f0f0;
  border-radius: 10rpx;
  overflow: hidden;
  display: flex;
  margin: 15rpx 0;
  box-shadow: inset 0 2rpx 4rpx rgba(0, 0, 0, 0.1);
}

.dimension-value {
  height: 100%;
  transition: width 0.6s ease;
}

.report-section {
  margin-top: 50rpx;
  padding-top: 40rpx;
  border-top: 2rpx solid #eee;
}

.section-title {
  font-size: 36rpx;
  color: #333;
  font-weight: bold;
  margin-bottom: 30rpx;
  display: block;
}

.report-content {
  font-size: 30rpx;
  color: #666;
  line-height: 1.8;
  text-align: justify;
}

.share-button {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  padding: 25rpx 0;
  border-radius: 50rpx;
  font-size: 32rpx;
  border: none;
  width: 90%;
  margin: 40rpx auto;
  box-shadow: 0 4rpx 12rpx rgba(76, 175, 80, 0.3);
  transition: all 0.3s ease;
}

.share-button:active {
  transform: scale(0.98);
  box-shadow: 0 2rpx 6rpx rgba(76, 175, 80, 0.2);
}
</style>