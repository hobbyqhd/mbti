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
      }
    }
  },
  computed: {
    formattedReport() {
      return this.result.report
        .replace(/### ([^\n]+)/g, '<h3>$1</h3>')
        .replace(/#### ([^\n]+)/g, '<h4>$1</h4>')
        .replace(/- \*\*([^*]+)\*\*/g, '<li><strong>$1</strong>')
        .replace(/- ([^\n]+)/g, '<li>$1</li>')
        .replace(/\n\n/g, '</p><p>')
        .replace(/\n/g, '<br>')
        .replace(/<li>/g, '<p><li>')
        .replace(/<\/li>/g, '</li></p>')
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
  font-size: 20px;
  color: #333;
  margin-bottom: 10px;
  display: block;
}

.mbti-type {
  font-size: 64rpx;
  font-weight: bold;
  color: #007AFF;
  display: block;
  letter-spacing: 4rpx;
  text-shadow: 2rpx 2rpx 4rpx rgba(0, 0, 0, 0.1);
}

.report-section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20rpx;
  padding: 40rpx;
  margin: 40rpx 0;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  line-height: 1.8;
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
}

.report-section >>> h3 {
  font-size: 36rpx;
  font-weight: 600;
  color: #2c3e50;
  margin: 40rpx 0 20rpx;
  padding-bottom: 16rpx;
  border-bottom: 2rpx solid rgba(0, 0, 0, 0.1);
}

.report-section >>> h4 {
  font-size: 32rpx;
  font-weight: 600;
  color: #34495e;
  margin: 30rpx 0 16rpx;
}

.report-section >>> p {
  margin: 16rpx 0;
  color: #333;
  font-size: 28rpx;
  line-height: 1.8;
}

.report-section >>> li {
  margin: 12rpx 0;
  list-style-type: disc;
  margin-left: 40rpx;
  color: #333;
  font-size: 28rpx;
  line-height: 1.8;
}

.report-section >>> strong {
  color: #2c3e50;
  font-weight: 600;
}

.share-button,
.restart-button {
  width: 100%;
  margin-bottom: 20rpx;
  padding: 24rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  border: none;
  transition: all 0.2s ease;
  box-shadow: 0 2rpx 6rpx rgba(0, 122, 255, 0.2);
}

.share-button {
  background-color: #007AFF;
  color: white;
}

.restart-button {
  background-color: #34C759;
  color: white;
}

.share-button:active,
.restart-button:active,
.export-button:active {
  transform: scale(0.98);
  opacity: 0.9;
}

.export-button {
  width: 100%;
  margin-bottom: 20rpx;
  padding: 24rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  border: none;
  transition: all 0.2s ease;
  box-shadow: 0 2rpx 6rpx rgba(0, 122, 255, 0.2);
  background-color: #9254de;
  color: white;
}
</style>