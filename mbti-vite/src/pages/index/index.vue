<template>
  <view class="content">
    <view class="header">
      <text class="title">MBTI性格测试</text>
      <text class="subtitle">了解你的性格类型</text>
    </view>
    <view class="description">
      <text>MBTI（迈尔斯-布里格斯类型指标）是一个帮助你了解自己性格特征的测试工具。通过回答一系列问题，你将发现自己的性格类型。</text>
    </view>
    <button class="start-button" @click="startTest">开始测试</button>
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
    startTest() {
      uni.navigateTo({
        url: '/pages/test/test'
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
  justify-content: center;
  padding: 40rpx;
  min-height: 100vh;
}

.header {
  text-align: center;
  margin-bottom: 60rpx;
}

.title {
  font-size: 48rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.subtitle {
  font-size: 32rpx;
  color: #666;
  display: block;
}

.description {
  background-color: #f5f5f5;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 60rpx;
  text-align: center;
}

.description text {
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
}

.start-button {
  background-color: #4CAF50;
  color: white;
  padding: 30rpx 80rpx;
  border-radius: 50rpx;
  font-size: 32rpx;
  border: none;
  margin-top: 40rpx;
}

.start-button:active {
  background-color: #45a049;
}
</style>
