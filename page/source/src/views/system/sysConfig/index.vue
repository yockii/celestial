<script setup lang="ts">
import { ref, h } from "vue"
import { NSwitch, NDataTable, NCard, useMessage } from "naive-ui"
import dayjs from "dayjs"
import { getSysConfigList, updateSysConfig, SysConfig } from "@/service/api/settings/sysConfig"
import { useUserStore } from "@/store/user"

const message = useMessage()
const userStore = useUserStore()

const loading = ref(false)
const list = ref<SysConfig[]>([])

// 前端展示用的配置项说明
const configLabels: Record<string, string> = {
  "login.usernamePasswordEnabled": "启用用户名密码登录/注册"
}

const refresh = () => {
  loading.value = true
  getSysConfigList()
    .then((res) => {
      list.value = res || []
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUpdate = (row: SysConfig, value: boolean) => {
  updateSysConfig(row.key, value ? "true" : "false").then((res) => {
    if (res) {
      row.value = value ? "true" : "false"
      message.success("修改成功，已实时生效")
    } else {
      refresh()
    }
  })
}

const columns = [
  {
    title: "配置项",
    key: "key",
    render: (row: SysConfig) => configLabels[row.key] || row.key
  },
  {
    title: "说明",
    key: "comment"
  },
  {
    title: "更新时间",
    key: "updateTime",
    render: (row: SysConfig) => (row.updateTime ? dayjs(row.updateTime).format("YYYY-MM-DD HH:mm:ss") : "-")
  },
  {
    title: "状态",
    key: "value",
    render: (row: SysConfig) => {
      if (row.key !== "login.usernamePasswordEnabled" || !userStore.hasResourceCode("system:sysConfig:update")) {
        return row.value === "true" ? "开启" : "关闭"
      }
      return h(NSwitch, {
        value: row.value === "true",
        onUpdateValue: (value: boolean) => handleUpdate(row, value)
      })
    }
  }
]

refresh()
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>安全设置</n-h3>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-card :bordered="false">
        <n-data-table size="small" :data="list" :loading="loading" :row-key="(row: SysConfig) => row.key" :columns="columns" />
        <n-alert type="info" style="margin-top: 8px"> 配置修改后实时生效，无需重启应用。关闭“用户名密码登录”后，登录页仅保留第三方账号登录方式。 </n-alert>
      </n-card>
    </n-gi>
  </n-grid>
</template>

<style scoped></style>
