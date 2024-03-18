<template>
  <div class="lan_header">
    <div class="search">
      <span class="info_text">查询条件：</span>
      <span>
          <a-space direction="vertical" size="large">
            <a-input type="text"
                     :style="{width: searchBoxW}"
                     v-model="searchCondition"
                     placeholder='app="thinkphp"'
                     @pressEnter="Submit"
            />
          </a-space>
      </span>
      <span>
          <a-space>
            <a-button v-if="subBtn"
                      class="submit_btn"
                      status="success"
                      @click="Submit()"
            >开始查询</a-button>
            <a-button v-else
                      class="submit_btn"
                      status="danger"
                      @click="Submit()"
            >停止查询</a-button>
          </a-space>
      </span>
    </div>

    <div class="search">
      <span class="info_text">Cookie信息：</span>
      <span>
          <a-space direction="vertical" size="large">
            <a-input type="text"
                     :style="{width: searchBoxW}"
                     placeholder='不填则使用UnAuth模式'
                     v-model="cookie"
            />
          </a-space>
        </span>
      <span>
          <a-space>
            <a-button @click="Export"
                      class="submit_btn"
                      status="success"
            >
              导出数据
            </a-button>
          </a-space>
        </span>
    </div>

    <div class="search others">
      <div>
        <span class="info_number">线程数：</span>
        <span>
          <a-space direction="vertical" size="large">
            <a-input-number type="text"
                            :style="{width: '70px'}"
                            v-model="threads"
                            :size="'small'"
            />
          </a-space>
        </span>
      </div>

      <div>
        <span class="info_number">资产数：</span>
        <span>
          <a-space direction="vertical" size="large">
            <a-input-number type="text"
                            :style="{width: '90px'}"
                            v-model="assets"
                            :size="'small'"
            />
          </a-space>
        </span>
      </div>

      <div class="check_alive">
        <span>存活检测：</span>
        <span>
            <a-space size="small">
              <a-switch v-model="checkAlive"
                        :loading="!subBtn">
                <template #checked>
                  ON
                </template>
                <template #unchecked>
                  OFF
                </template>
              </a-switch>
          </a-space>
        </span>
      </div>

      <div class="outfile">
        <span>文件导出类型：</span>
        <span>
          <a-space size="large">
            <a-select :style="{width:'160px'}"
                      :size="'small'"
                      placeholder="文件导出类型"
                      multiple
                      :max-tag-count="1"
                      v-model="exportFileType"
            >
              <a-option value="txt">TXT</a-option>
              <a-option value="json">JSON</a-option>
              <a-option value="csv">CSV</a-option>
            </a-select>
          </a-space>
        </span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref} from "vue";
import {Message} from "@arco-design/web-vue";
import {Size, WindowGetSize} from "../../wailsjs/runtime";

/**
 * 计算搜索框宽度
 */
let searchBoxW = ref("400px")

/**
 * 计算搜索框大小
 * @constructor
 */
function HandleSearchBoxWidth() {
  WindowGetSize().then((size: Size) => {
    searchBoxW.value = `${size.w - 350}px`
  })
}

window.addEventListener("resize", HandleSearchBoxWidth)
HandleSearchBoxWidth()


/**
 * 查询事件
 */
let searchCondition = ref<string>("") // 查询条件
let cookie = ref<string>("")  // cookie信息
let subBtn = ref<boolean>(true) // 按钮状态，true显示开始查询，false显示停止查询
let threads = ref<number>(10) // 线程数
let assets = ref<number>(50)  // 资产数
let checkAlive = ref<boolean>(false)  // 存活检测

/**
 * 提交查询和停止查询
 * 通过 subBtn.value 值进行判断执行哪个事件，也即按钮状态
 * @constructor
 */
function Submit() {
  if (subBtn.value) {
    // 提交查询
    if (searchCondition.value == "") {
      Message.error("查询内容不能为空")
      return
    }

    Message.success("开始查询")
  } else {
    // 停止查询
    Message.error("停止查询")
  }

  subBtn.value = !subBtn.value
}


/**
 * 导出文件
 */
let exportFileType = ref<string[]>([])

/**
 * 导出资产到文件中
 * @constructor
 */
function Export() {
  if (exportFileType.value.length === 0) {
    Message.error("请选择文件导出类型")
    return
  }

  // 导出文件
  Message.success("导出成功，保存在/tmp/results")
}
</script>

<style scoped>
.lan_header {
  padding-top: 1px;
  margin-left: 10px;
}

.search {
  margin-bottom: 12px;
}

.submit_btn {
  margin-left: 30px;
}

.info_text {
  display: inline-block;
  width: 100px;
}

.info_number {
  display: inline-block;
  width: 60px;
}

.others {
  justify-content: center;
  display: flex;
}

.others > div {
  margin-right: 30px;
}

.check_alive {
  width: 135px;
  line-height: 35px;
}

.outfile {
  line-height: 35px;
}
</style>