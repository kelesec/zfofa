<template>
  <div>
    <div>
      <a-table row-key="id"
               :columns="columns"
               :data="tableData"
               :stripe="true"
               :scrollbar="true"
               :scroll="{y: scrollH}"
               :size="'medium'"
               :pagination="false">
        <template #host="{ record }">
          <a-table-column>
            <a-link href="javascript:void(0);" @click="OpenUrlByDefaultBrowser(record.host)">
              <a-avatar :size="16" :style="{marginRight: '5px'}">
                <img v-if="record.icon" alt="avatar" :src="`data:image/png;base64, ${record.icon}`"/>
                <span v-else>ZFOFA</span>
              </a-avatar>
              <span>{{ record.host }}</span>
            </a-link>
          </a-table-column>
        </template>
        <template #ip="{ record }">
          <a-table-column>
            <a-link href="javascript:void(0);"
                    @click="OpenIpByIp138(record.ip)"
                    icon>{{ record.ip }}</a-link>
          </a-table-column>
        </template>
        <template #header="{ record }">
          <a-table-column>
            <span @click="ShowDrawer('Header', record.header)">{{ record.header }}</span>
          </a-table-column>
        </template>
        <template #cert="{ record }">
          <a-table-column>
            <span @click="ShowDrawer('Cert', record.certificate)">{{ record.certificate }}</span>
          </a-table-column>
        </template>
      </a-table>
    </div>
    <div class="lan_pagination">
      <a-pagination :size="'small'"
                    :simple="true"
                    :show-total="true"
                    :show-page-size="true"
                    :total="tableAssets.length"
                    @change="CurPageChange"
                    @pageSizeChange="PageSizeChange"
      >
      </a-pagination>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {reactive, ref} from 'vue';
import {TableColumnData, TableData} from "@arco-design/web-vue";
import {COLUMN_DATA, DataType, ShowDrawer} from "../utils/layout_content";
import {BrowserOpenURL, EventsOn} from "../../wailsjs/runtime";

/**
 * 表格中每一列的列类型
 */
const columns: TableColumnData[] = reactive(COLUMN_DATA as TableColumnData[])

/**
 * 填充表格数据
 */
let tableData = ref<TableData[]>([])    // 填入表格的数据
let tableAssets: TableData[] = reactive([])  // 总数据

/**
 * 构建 tableAssets 数据填充事件
 * @param assets DataType[] 数组
 * @constructor
 */
function PushAssets(assets: DataType[]) {
  tableAssets = assets.slice(0, assets.length)
  FlashTableData()
}
EventsOn("PushAssets", PushAssets)

/**
 * 分页数据绑定
 */
let curPage = ref<number>(1)
let perPageSize = ref<number>(10)

/**
 * 刷新表格数据
 * @constructor
 */
function FlashTableData() {
  let startIndex = (curPage.value - 1) * perPageSize.value
  let endIndex = startIndex + perPageSize.value
  tableData.value = tableAssets.slice(startIndex, endIndex)
}

FlashTableData()

/**
 * 页数该表示重新计算将要展示的数据
 * @param currentPage 当前页码
 * @constructor
 */
function CurPageChange(currentPage: number) {
  curPage.value = currentPage
  FlashTableData()
}

/**
 * 每页将要显示的数据量变化时，重新计算和刷新数据
 * @param pageSize 当前页将要显示多少数据
 * @constructor
 */
function PageSizeChange(pageSize: number) {
  perPageSize.value = pageSize
  FlashTableData()
}


/**
 * 表格滚动条大小计算
 */
let scrollH = ref<number>(280)

/**
 * 计算滚动条高度
 * @constructor
 */
function HandleWindowSize() {
  let appEle = document.getElementById('app') as HTMLElement
  scrollH.value = 280 + (appEle.offsetHeight - 585)
}

window.addEventListener('resize', HandleWindowSize)
HandleWindowSize()


/**
 * 使用默认浏览器打开URL
 * @param url
 * @constructor
 */
function OpenUrlByDefaultBrowser(url: string) {
  if (!url.startsWith("http")) {
    url = `http://${url}`
  }
  BrowserOpenURL(url)
}

/**
 * 点击IP地址链接时，跳转查询IP信息
 * @param ip
 * @constructor
 */
function OpenIpByIp138(ip: string) {
  BrowserOpenURL(`https://www.ip138.com/iplookup.php?ip=${ip}`)
}
</script>

<style scoped>
.lan_pagination {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>