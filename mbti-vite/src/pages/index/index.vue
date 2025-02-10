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
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
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
  background-color: rgba(239, 239, 244, 0.9);
  padding: 30rpx;
  border-radius: 20rpx;
  margin-bottom: 60rpx;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
}

.description text {
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
}

.start-button {
  background-color: #007AFF;
  color: white;
  padding: 30rpx 80rpx;
  border-radius: 12rpx;
  font-size: 32rpx;
  border: none;
  margin-top: 40rpx;
  box-shadow: 0 2px 6px rgba(0, 122, 255, 0.2);
  transition: all 0.2s ease;
}

.start-button:active {
  background-color: #0051D5;
  transform: scale(0.98);
}
</style>
