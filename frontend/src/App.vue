<template>
  <div class="layout">
    <a-layout :style="{height: `${windowScreenSize.H}px`}">
      <a-layout-header :style="{height: `${windowScreenSize.headerH}px`}">
        <layout_header/>
      </a-layout-header>

      <div id="hc_space"></div>
      <a-layout-content :style="{height: `${windowScreenSize.contentH}px`}">
        <layout_content/>
      </a-layout-content>

      <a-layout-footer :style="{height: `${windowScreenSize.footerH}`}">
        <layout_footer/>
      </a-layout-footer>
    </a-layout>
  </div>
</template>

<script lang="ts" setup>
import {ref} from "vue";
import {Size, WindowGetSize} from "../wailsjs/runtime";
import Layout_header from "./components/layout_header.vue";
import Layout_footer from "./components/layout_footer.vue";
import Layout_content from "./components/layout_content.vue";

// 设置窗口大小
let windowScreenSize = ref({
  H: 620,
  headerH: 150,
  contentH: 360,
  footerH: 70,
  cfH: 43
})

function HandleWindowSize() {
  WindowGetSize().then((size: Size) => {
    let bodyH = size.h - 40 // 去掉标题位置处的高度
    let headerH = 150
    let footerH = 70
    let contentH = bodyH - headerH - footerH

    windowScreenSize.value.H = bodyH
    windowScreenSize.value.headerH = headerH
    windowScreenSize.value.contentH = contentH
    windowScreenSize.value.footerH = footerH
    windowScreenSize.value.cfH = contentH + footerH
  })
}

window.addEventListener("resize", HandleWindowSize)
HandleWindowSize()
</script>

<style scoped>
.layout :deep(.arco-layout-header),
.layout :deep(.arco-layout-content) {
  display: flex;
  flex-direction: column;
  color: #17171a;
  font-size: 14px;
  font-stretch: condensed;
}

.layout :deep(.arco-layout-header) {
  padding: 10px;
  background-color: #e5e6eb;
}

.layout :deep(.arco-layout-content) {
  background-color: #f2f3f5;
}

#hc_space {
  border: 1px solid grey;
  border-bottom: 0;
}
</style>
