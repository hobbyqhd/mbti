<template>
  <view class="content">
    <view class="header">
      <text class="title">MBTI性格测试</text>
      <text class="subtitle">了解你的性格类型</text>
    </view>
    
    <view class="description">
      <text>MBTI（迈尔斯-布里格斯类型指标）是一个帮助你了解自己性格特征的测试工具。请选择适合您的测试版本：</text>
    </view>

    <view class="test-options">
      <view class="test-card" @click="startTest('simple')">
        <text class="card-title">简单版</text>
        <text class="card-subtitle">40题 · 约10分钟</text>
        <text class="card-description">适合快速了解自己的性格类型</text>
        <button class="start-button primary">开始测试</button>
      </view>

      <view class="test-card" @click="startTest('detailed')">
        <text class="card-title">详细版</text>
        <text class="card-subtitle">80题 · 约20分钟</text>
        <text class="card-description">提供更准确的性格分析结果</text>
        <button class="start-button">开始测试</button>
      </view>

      <view class="test-card" @click="startTest('full')">
        <text class="card-title">超详细版</text>
        <text class="card-subtitle">120题 · 约30分钟</text>
        <text class="card-description">最全面的性格测试体验</text>
        <button class="start-button">开始测试</button>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      title: 'MBTI性格测试'
    }
  },
  methods: {
    startTest(type) {
      uni.navigateTo({
        url: `/pages/test/test?type=${type}`
      })
    }
  }
}
</script>

<style>
.content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 40rpx;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.header {
  text-align: center;
  margin-bottom: 40rpx;
}

.title {
  font-size: 48rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.subtitle {
  font-size: 32rpx;
  color: #666;
  display: block;
}

.description {
  background-color: #fff;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 40rpx;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  width: 90%;
}

.test-options {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 30rpx;
  padding: 20rpx;
}

.test-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 40rpx;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.card-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
}

.card-subtitle {
  font-size: 28rpx;
  color: #666;
  margin-bottom: 20rpx;
}

.card-description {
  font-size: 26rpx;
  color: #888;
  margin-bottom: 30rpx;
}

.start-button {
  width: 80%;
  height: 80rpx;
  line-height: 80rpx;
  border-radius: 40rpx;
  background-color: #f0f0f0;
  color: #333;
  font-size: 28rpx;
  border: none;
}

.start-button.primary {
  background-color: #4CAF50;
  color: #fff;
}

.test-card:hover {
  transform: translateY(-2rpx);
  transition: all 0.3s ease;
}
</style>
