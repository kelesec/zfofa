<template>
  <div class="lan_footer">
    <a-divider orientation="left" :margin="9" :size="1.5">
      日志信息
    </a-divider>

    <textarea class="lan_textarea"
              :value="logData"
              ref="scroll"
              readonly>
    </textarea>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from "vue";
import {EventsOn} from "../../wailsjs/runtime";

let logInfo = ref<string[]>(["zfofa by kele."])
let logData = ref(logInfo.value[0])
let scroll = ref<HTMLTextAreaElement>() // 滚动条

watch(logInfo.value, (newLog) => {
  if (newLog.length === 0) {
    return
  }

  if (newLog.length > 10) {
    logInfo.value.splice(1, 1)
  }

  logData.value = newLog.join('\n')

  // 确保滚动条始终在最下方
  if (scroll.value !== undefined) {
    scroll.value.scrollTop = scroll.value.scrollHeight + 10
  }
})

function PushLog(text: string) {
  logInfo.value.push(text)
}

function ClearLog() {
  logInfo.value.splice(1, logInfo.value.length)
}

// 注册监听事件
EventsOn('PushLog', PushLog)
EventsOn('ClearLog', ClearLog)
</script>

<style scoped>
.lan_footer {
  padding: 5px;
  background-color: #ffffff;
}

.lan_textarea {
  margin-top: 5px;
  background-color: #293134;
  color: #ffffff;
  padding: 0.5%;
  width: 99%;
  height: 40px;
  resize: none;
  font-size: 12px;
}
</style>