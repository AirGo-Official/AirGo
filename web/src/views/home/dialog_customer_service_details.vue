<template>
  <div>
    <el-dialog v-model="state.isShowDialog" width="90%" destroy-on-close>
        <div class="home-container layout-pd">
          <div class="home-card-item" style="height: 200px">
            <div style="height: 100%" ref="homeLineRef"></div>
          </div>
        </div>
    </el-dialog>
  </div>

</template>

<script setup lang="ts" name="home">
import { reactive, onMounted, ref, watch, nextTick, onActivated, markRaw } from 'vue';
import * as echarts from 'echarts';
import { storeToRefs } from 'pinia';
import { useThemeConfig } from '/@/stores/themeConfig';
import { useTagsViewRoutes } from '/@/stores/tagsViewRoutes';
import { useTrafficStore } from "/@/stores/user_logic/trafficStore";
import { useI18n } from "vue-i18n";

// 定义变量内容
const homeLineRef = ref();
const storesTagsViewRoutes = useTagsViewRoutes();
const storesThemeConfig = useThemeConfig();
const { themeConfig } = storeToRefs(storesThemeConfig);
const { isTagsViewCurrenFull } = storeToRefs(storesTagsViewRoutes);
const trafficStore = useTrafficStore()
const trafficStoreData = storeToRefs(trafficStore)
const {t} = useI18n()
const state = reactive({
  isShowDialog:false,
  global: {
    homeChartOne: null,
    homeChartTwo: null,
    homeCharThree: null,
    dispose: [null, '', undefined],
  } as any,
  myCharts: [] as EmptyArrayType,
  charts: {
    theme: '',
    bgColor: '',
    color: '#303133',
  },
  customerServiceID:0,
});
const openDialog = (customerServiceID:number) => {
  console.log("customerServiceID:",customerServiceID)
  state.customerServiceID = customerServiceID
  state.isShowDialog = true;
  trafficStore.getSubTrafficList({id:state.customerServiceID } as UserTrafficLog).then(()=>{
    initLineChart(trafficStoreData.trafficLineChart.value.xAxis,trafficStoreData.trafficLineChart.value.u,trafficStoreData.trafficLineChart.value.d);
  })
};
const closeDialog = () => {
  state.isShowDialog = false;
};

// 暴露变量
defineExpose({
  openDialog,
});

// 折线图
const initLineChart = (xAxis:string[],u:number[],d:number[]) => {
  if (!state.global.dispose.some((b: any) => b === state.global.homeChartOne)) state.global.homeChartOne.dispose();
  state.global.homeChartOne = markRaw(echarts.init(homeLineRef.value, state.charts.theme));
  const option = {
    backgroundColor: state.charts.bgColor,
    title: {
      text: t('message.home.traffic_log'),
      x: 'left',
      textStyle: { fontSize: '15', color: state.charts.color },
    },
    grid: { top: 70, right: 20, bottom: 30, left: 30 },
    tooltip: { trigger: 'axis' },
    legend: { data: [t('message.home.upstream_traffic'), t('message.home.downstream_traffic')], right: 0 },
    xAxis: {
      axisLabel: { //设置x轴的字
        show: true,
        interval: 0,//使x轴横坐标全部显示
        rotate: 20,
        textStyle: {//x轴字体样式
          fontSize: 8,
        },
      },
      data: xAxis,
    },
    // dataZoom:[
    //   {
    //     show: true,
    //     start: 0,
    //     end: 50
    //   }
    // ],
    yAxis: [
      {
        type: 'value',
        name: 'GB',
        splitLine: { show: true, lineStyle: { type: 'dashed', color: '#f5f5f5' } },
      },
    ],
    series: [
      {
        name: t('message.home.upstream_traffic'),
        type: 'line',
        symbolSize: 6,
        symbol: 'circle',
        smooth: true,
        data: u,
        lineStyle: { color: '#fe9a8b' },
        itemStyle: { color: '#fe9a8b', borderColor: '#fe9a8b' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#fe9a8bb3' },
            { offset: 1, color: '#fe9a8b03' },
          ]),
        },
      },
      {
        name: t('message.home.downstream_traffic'),
        type: 'line',
        symbolSize: 6,
        symbol: 'circle',
        smooth: true,
        data: d,
        lineStyle: { color: '#9E87FF' },
        itemStyle: { color: '#9E87FF', borderColor: '#9E87FF' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#9E87FFb3' },
            { offset: 1, color: '#9E87FF03' },
          ]),
        },
        emphasis: {
          itemStyle: {
            color: {
              type: 'radial',
              x: 0.5,
              y: 0.5,
              r: 0.5,
              colorStops: [
                { offset: 0, color: '#9E87FF' },
                { offset: 0.4, color: '#9E87FF' },
                { offset: 0.5, color: '#fff' },
                { offset: 0.7, color: '#fff' },
                { offset: 0.8, color: '#fff' },
                { offset: 1, color: '#fff' },
              ],
            },
            borderColor: '#9E87FF',
            borderWidth: 2,
          },
        },
      },
    ],
  };
  state.global.homeChartOne.setOption(option);
  state.myCharts.push(state.global.homeChartOne);
};

// 柱状图

// 批量设置 echarts resize
const initEchartsResizeFun = () => {
  nextTick(() => {
    for (let i = 0; i < state.myCharts.length; i++) {
      setTimeout(() => {
        state.myCharts[i].resize();
      }, i * 1000);
    }
  });
};
// 批量设置 echarts resize
const initEchartsResize = () => {
  window.addEventListener('resize', initEchartsResizeFun);
};
// 页面加载时
onMounted(() => {
  initEchartsResize();
});
// 由于页面缓存原因，keep-alive
onActivated(() => {
  initEchartsResizeFun();
});
// 监听 pinia 中的 tagsview 开启全屏变化，重新 resize 图表，防止不出现/大小不变等
watch(
  () => isTagsViewCurrenFull.value,
  () => {
    initEchartsResizeFun();
  }
);
</script>

<style scoped lang="scss">

$homeNavLengh: 8;
.home-container {
  overflow: hidden;
  .home-card-one,
  .home-card-two,
  .home-card-three {
    .home-card-item {
      width: 100%;
      height: 130px;
      border-radius: 4px;
      transition: all ease 0.3s;
      padding: 20px;
      overflow: hidden;
      background: var(--el-color-white);
      color: var(--el-text-color-primary);
      border: 1px solid var(--next-border-color-light);
      &:hover {
        box-shadow: 0 2px 12px var(--next-color-dark-hover);
        transition: all ease 0.3s;
      }
      &-icon {
        width: 70px;
        height: 70px;
        border-radius: 100%;
        flex-shrink: 1;
        i {
          color: var(--el-text-color-placeholder);
        }
      }
      &-title {
        font-size: 15px;
        font-weight: bold;
        height: 30px;
      }
    }
  }
  .home-card-one {
    @for $i from 0 through 3 {
      .home-one-animation#{$i} {
        opacity: 0;
        animation-name: error-num;
        animation-duration: 0.5s;
        animation-fill-mode: forwards;
        animation-delay: calc($i/4) + s;
      }
    }
  }
  .home-card-two,
  .home-card-three {
    .home-card-item {
      height: 400px;
      width: 100%;
      overflow: hidden;
      .home-monitor {
        height: 100%;
        .flex-warp-item {
          width: 25%;
          height: 111px;
          display: flex;
          .flex-warp-item-box {
            margin: auto;
            text-align: center;
            color: var(--el-text-color-primary);
            display: flex;
            border-radius: 5px;
            background: var(--next-bg-color);
            cursor: pointer;
            transition: all 0.3s ease;
            &:hover {
              background: var(--el-color-primary-light-9);
              transition: all 0.3s ease;
            }
          }
          @for $i from 0 through $homeNavLengh {
            .home-animation#{$i} {
              opacity: 0;
              animation-name: error-num;
              animation-duration: 0.5s;
              animation-fill-mode: forwards;
              animation-delay: calc($i/10) + s;
            }
          }
        }
      }
    }
  }
}
</style>
