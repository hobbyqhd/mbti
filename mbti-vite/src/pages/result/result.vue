<template>
  <view class="result-container">
    <view v-if="loading" class="loading">
      <text>加载中...</text>
    </view>

    <view v-else class="result-content">
      <view class="result-header">
        <text class="result-title">您的MBTI类型是</text>
        <text class="mbti-type">{{ result.type }}</text>
      </view>

      <view class="dimensions-section">
        <view class="dimension-item" v-for="(dimension, index) in result.dimensions" :key="index">
          <view class="dimension-labels">
            <text class="dimension-label">{{dimension.left}}</text>
            <text class="dimension-label">{{dimension.right}}</text>
          </view>
          <view class="dimension-bar">
            <view class="dimension-progress left" :style="{ width: dimension.leftValue + '%' }"></view>
            <view class="dimension-progress right" :style="{ width: dimension.rightValue + '%' }"></view>
          </view>
          <view class="dimension-scores">
            <text class="score-text">{{dimension.leftValue}}%</text>
            <text class="score-text">{{dimension.rightValue}}%</text>
          </view>
        </view>
      </view>

      <view class="report-section">
        <rich-text :nodes="formattedReport"></rich-text>
      </view>

      <view class="button-group">
        <button class="share-button" @click="shareResult">分享结果</button>
        <button class="export-button" @click="exportPDF">导出PDF</button>
        <button class="restart-button" @click="restartTest">重新测试</button>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      loading: true,
      result: {
        type: '',
        dimensions: [
          { left: 'E (外向)', right: 'I (内向)', leftValue: 0, rightValue: 0 },
          { left: 'S (感觉)', right: 'N (直觉)', leftValue: 0, rightValue: 0 },
          { left: 'T (思维)', right: 'F (情感)', leftValue: 0, rightValue: 0 },
          { left: 'J (判断)', right: 'P (知觉)', leftValue: 0, rightValue: 0 }
        ],
        report: ''
      }
    }
  },
  computed: {
    formattedReport() {
      return this.result.report;
    },
  },
  onLoad(options) {
    if (options.id) {
      this.fetchResult(options.id)
    }
  },
  onMounted() {
    const loadScript = (src) => {
      return new Promise((resolve, reject) => {
        const script = document.createElement('script')
        script.src = src
        script.onload = resolve
        script.onerror = reject
        document.head.appendChild(script)
      })
    }

    Promise.all([
      loadScript('https://cdnjs.cloudflare.com/ajax/libs/html2canvas/1.4.1/html2canvas.min.js'),
      loadScript('https://cdnjs.cloudflare.com/ajax/libs/jspdf/2.5.1/jspdf.umd.min.js')
    ]).catch(error => {
      console.error('加载CDN资源失败:', error)
    })
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
      uni.showToast({
        title: '分享功能开发中',
        icon: 'none'
      })
    },
    restartTest() {
      uni.redirectTo({
        url: '/pages/test/test'
      })
    },
    async exportPDF() {
      try {
        uni.showLoading({
          title: '正在生成PDF...'
        });

        const element = document.querySelector('.result-content');
        const canvas = await window.html2canvas(element, {
          scale: 2,
          useCORS: true,
          logging: false
        });

        const { jsPDF } = window.jspdf;
        const pdf = new jsPDF({
          orientation: 'portrait',
          unit: 'mm',
          format: 'a4'
        });

        const imgWidth = 210;
        const pageHeight = 297;
        const imgHeight = canvas.height * imgWidth / canvas.width;

        const imgData = canvas.toDataURL('image/jpeg', 1.0);
        pdf.addImage(imgData, 'JPEG', 0, 0, imgWidth, imgHeight);

        const fileName = `MBTI测试结果_${this.result.type}_${new Date().toLocaleDateString()}.pdf`;

        pdf.save(fileName);

        uni.hideLoading();
        uni.showToast({
          title: 'PDF导出成功',
          icon: 'success'
        });
      } catch (error) {
        console.error('PDF导出失败:', error);
        uni.hideLoading();
        uni.showToast({
          title: 'PDF导出失败',
          icon: 'none'
        });
      }
    }
  }
}
</script>

<style>
.result-container {
  padding: 30rpx;
  min-height: 100vh;
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.result-content {
  max-width: 800rpx;
  margin: 0 auto;
  padding: 40rpx;
}

.result-header {
  text-align: center;
  margin-bottom: 50rpx;
  animation: fadeIn 0.8s ease-out;
}

.result-title {
  font-size: 32rpx;
  color: #666;
  margin-bottom: 20rpx;
  display: block;
}

.mbti-type {
  font-size: 72rpx;
  font-weight: bold;
  color: #007AFF;
  display: block;
  letter-spacing: 4rpx;
  text-shadow: 2rpx 2rpx 4rpx rgba(0, 0, 0, 0.1);
}

.dimensions-section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20rpx;
  padding: 40rpx;
  margin: 40rpx 0;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  animation: slideUp 0.8s ease-out;
}

.dimension-item {
  margin-bottom: 40rpx;
}

.dimension-labels {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15rpx;
}

.dimension-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.dimension-bar {
  height: 24rpx;
  background: #f0f0f0;
  border-radius: 12rpx;
  overflow: hidden;
  display: flex;
  box-shadow: inset 0 2rpx 4rpx rgba(0, 0, 0, 0.1);
}

.dimension-progress {
  height: 100%;
  transition: width 0.6s ease;
}

.dimension-progress.left {
  background: linear-gradient(to right, #4CAF50, #81C784);
}

.dimension-progress.right {
  background: linear-gradient(to right, #2196F3, #64B5F6);
}

.dimension-scores {
  display: flex;
  justify-content: space-between;
  margin-top: 10rpx;
}

.score-text {
  font-size: 24rpx;
  color: #666;
}

.report-section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20rpx;
  padding: 40rpx;
  margin: 40rpx 0;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  animation: slideUp 1s ease-out;
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  margin-top: 40rpx;
  animation: slideUp 1.2s ease-out;
}

.share-button,
.restart-button,
.export-button {
  width: 100%;
  padding: 25rpx;
  border-radius: 12rpx;
  font-size: 32rpx;
  font-weight: 500;
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
  color: white;
}

.share-button {
  background: linear-gradient(135deg, #007AFF, #0056b3);
}

.restart-button {
  background: linear-gradient(135deg, #34C759, #28a745);
}

.export-button {
  background: linear-gradient(135deg, #9254de, #722ed1);
}

.share-button:active,
.restart-button:active,
.export-button:active {
  transform: scale(0.98);
  opacity: 0.9;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>