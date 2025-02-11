<template>
  <view class="result-container">
    <view v-if="loading" class="loading">
      <text>加载中...</text>
    </view>

    <view v-else class="result-content">
      <view class="result-header">
        <text class="result-title">您的MBTI类型是 {{ result.mbtiType }}</text>
      </view>

      <view class="dimensions-section">
        <view class="dimension-item" v-for="(dimension, index) in dimensions" :key="index">
          <view class="dimension-labels">
            <text class="dimension-label">{{dimension.left}}</text>
            <text class="dimension-label">{{dimension.right}}</text>
          </view>
          <view class="dimension-bar">
            <view class="dimension-progress left" :style="{ width: dimension.leftScore + '%' }"></view>
            <view class="dimension-progress right" :style="{ width: dimension.rightScore + '%' }"></view>
          </view>
          <view class="dimension-scores">
            <text class="score-text">{{dimension.leftScore}}%</text>
            <text class="score-text">{{dimension.rightScore}}%</text>
          </view>
        </view>
      </view>

      <view class="report-section">
        <rich-text :nodes="formattedReport"></rich-text>
      </view>

      <button class="share-button" @click="shareResult">分享结果</button>
      <button class="export-button" @click="exportPDF">导出PDF</button>
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
      },
      dimensions: [
        { left: 'E (外向)', right: 'I (内向)', leftScore: 60, rightScore: 40 },
        { left: 'S (感觉)', right: 'N (直觉)', leftScore: 45, rightScore: 55 },
        { left: 'T (思维)', right: 'F (情感)', leftScore: 70, rightScore: 30 },
        { left: 'J (判断)', right: 'P (知觉)', leftScore: 35, rightScore: 65 }
      ]
    }
  },
  computed: {
    formattedReport() {
      return this.result.report
        .replace(/^### ([^\n]+)/gm, '<h3>$1</h3>')
        .replace(/^#### ([^\n]+)/gm, '<h4>$1</h4>')
        .replace(/^\s*[-•]\s*\*\*([^*]+)\*\*/gm, '<li><strong>$1</strong></li>')
        .replace(/^\s*[-•]\s*([^\n]+)/gm, '<li>$1</li>')
        .replace(/(<li>[\s\S]*?<\/li>\s*)+/g, '<ul>$&</ul>')
        .replace(/\n\n(?!<)/g, '</p><p>')
        .replace(/\n(?!<)/g, '<br>')
        .split(/\n(?=<h[34]>)/g)
        .map(section => section.trim())
        .join('\n')
    },
  },
  onLoad(options) {
    if (options.id) {
      this.fetchResult(options.id)
    }
  },
  onMounted() {
    // 动态加载CDN资源
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
    },
    async exportPDF() {
      try {
        uni.showLoading({
          title: '正在生成PDF...'
        });

        // 获取要导出的内容区域
        const element = document.querySelector('.result-content');
        const canvas = await window.html2canvas(element, {
          scale: 2,
          useCORS: true,
          logging: false
        });

        // 创建PDF文档
        const { jsPDF } = window.jspdf;
        const pdf = new jsPDF({
          orientation: 'portrait',
          unit: 'mm',
          format: 'a4'
        });

        // 获取画布尺寸和PDF页面尺寸
        const imgWidth = 210; // A4纸的宽度（单位：mm）
        const pageHeight = 297; // A4纸的高度（单位：mm）
        const imgHeight = canvas.height * imgWidth / canvas.width;

        // 将画布内容添加到PDF
        const imgData = canvas.toDataURL('image/jpeg', 1.0);
        pdf.addImage(imgData, 'JPEG', 0, 0, imgWidth, imgHeight);

        // 生成PDF文件名
        const fileName = `MBTI测试结果_${this.result.mbtiType}_${new Date().toLocaleDateString()}.pdf`;

        // 保存PDF
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
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
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
  font-size: 16px;
  color: #333;
  margin-bottom: 10px;
  display: block;
}

.mbti-type {
  font-size: 48rpx;
  font-weight: bold;
  color: #007AFF;
  display: block;
  letter-spacing: 4rpx;
  text-shadow: 2rpx 2rpx 4rpx rgba(0, 0, 0, 0.1);
}

.report-section {
  padding: 20rpx;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
  margin: 20rpx 0;
}

.report-section :deep(h3) {
  font-size: 36rpx;
  color: #333;
  margin: 30rpx 0 20rpx;
  font-weight: 600;
}

.report-section :deep(h4) {
  font-size: 32rpx;
  color: #444;
  margin: 24rpx 0 16rpx;
  font-weight: 500;
}

.report-section :deep(p) {
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
  margin: 16rpx 0;
}

.report-section :deep(ul) {
  padding-left: 30rpx;
  margin: 16rpx 0;
}

.report-section :deep(li) {
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
  margin: 12rpx 0;
  list-style-type: disc;
}

.report-section :deep(strong) {
  color: #333;
  font-weight: 600;
}

.dimensions-section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20rpx;
  padding: 40rpx;
  margin: 40rpx 0;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
}

.dimension-item {
  margin-bottom: 30rpx;
}

.dimension-labels {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10rpx;
}

.dimension-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.dimension-bar {
  height: 20rpx;
  background: #f0f0f0;
  border-radius: 10rpx;
  overflow: hidden;
  display: flex;
}

.dimension-progress {
  height: 100%;
  transition: width 0.3s ease;
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
  margin-top: 8rpx;
}

.score-text {
  font-size: 24rpx;
  color: #666;
}
.share-button,
.restart-button,
.export-button {
  width: 100%;
  margin-bottom: 20rpx;
  padding: 16rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  border: none;
  transition: all 0.2s ease;
  box-shadow: 0 2rpx 6rpx rgba(0, 122, 255, 0.2);
  color: white;
}

.share-button {
  background-color: #007AFF;
}

.restart-button {
  background-color: #34C759;
}

.export-button {
  background-color: #9254de;
}

.share-button:active,
.restart-button:active,
.export-button:active {
  transform: scale(0.98);
  opacity: 0.9;
}
</style>