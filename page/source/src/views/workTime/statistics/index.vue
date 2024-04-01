<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-space>
        <div class="flex flex-items-center">
          <router-link class="text-1.5em line-height-0.5em" :to="{ name: 'WorkTime' }">
            <n-icon>
              <KeyboardBackspaceOutlined />
            </n-icon>
          </router-link>
          返回
        </div>
        <n-text>工时统计表</n-text>
      </n-space>
    </n-gi>
    <n-gi>
      <n-space justify="space-between">
        <n-form inline label-placement="left" label-width="auto" size="small">
          <n-form-item label="筛选时间段：">
            <n-date-picker v-model:value="dateRange" type="daterange" />
          </n-form-item>
          <n-form-item label="选择部门：">
            <n-cascader v-model:value="deptId" :options="topDepartments" show-path remote :on-load="handleLoadChildren"
              value-field="id" label-field="name" children-field="children" />
          </n-form-item>
          <n-form-item>
            <n-button type="info" @click="search()">查询</n-button>
          </n-form-item>
        </n-form>
        <n-button type="primary" size="small" @click="handleOpenProjectConfigDrawer">项目配置</n-button>
      </n-space>
    </n-gi>
    <n-gi>
      <n-data-table :data="data" :single-line="false" :columns="columns" />
    </n-gi>
  </n-grid>

  <n-drawer title="项目配置" v-model:show="projectConfigDrawerActived" placement="right" width="30%">
    <n-drawer-content>
      <n-space vertical>
        <n-space v-for="project in projectListInstance" :key="project.id" justify="space-between">
          <n-text>{{ project.name }}</n-text>
          <n-select v-model:value="project.type" size="small" :options="projectConfigOptions" />
        </n-space>
      </n-space>
      <template #footer>
        <n-space>
          <n-button @click="projectConfigDrawerActived = false">取消</n-button>
          <n-button type="primary" @click="confirmProjectConfigUpdation()">确定</n-button>
        </n-space>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { KeyboardBackspaceOutlined } from "@vicons/material"
import dayjs from "dayjs"
import { getAllProjectListForWorkTimeStatistics } from "@/service"
import { Project } from "@/types/project"
import { Department } from "@/types/user"
import { WorkTime, WorkTimeStatisticsCondition } from "@/types/log"
import { getDepartmentList, getWorkTimeStatistics } from "@/service/api"

const dateRange = ref([dayjs().add(-1, "month").startOf("month").valueOf(), dayjs().add(-1, "month").endOf("month").valueOf()])
const deptId = ref("")
const topDepartments = ref<Department[]>([])

const handleLoadChildren = (node: Department) => {
  return new Promise<void>((resolve) => {
    getDepartmentList({ parentId: node.id, offset: -1, limit: -1 }).then((res) => {
      node.children = res.items.map((item) => {
        item.isLeaf = !item.childCount || item.childCount === 0
        return item
      })
      resolve()
    })
  })
}

type ProjectWithType = Project & {
  type: number
}
const projectList = ref<ProjectWithType[]>([])

// 表格数据处理
const columns = computed(() => {
  type Col = {
    title: string
    titleAlign?: string
    key?: string
    children?: Col[]
  }
  const cols: Col[] = [
    {
      title: "姓名",
      titleAlign: "center",
      key: "name"
    }
  ]

  const self: Col = {
    title: "自研项目",
    key: "self",
    titleAlign: "center",
    children: []
  }
  const share: Col = {
    title: "本月分摊",
    key: "share",
    titleAlign: "center",
    children: []
  }
  const cost: Col = {
    title: "待摊费用",
    key: "cost",
    titleAlign: "center",
    children: []
  }

  // 把项目分别放到不同的列中
  projectList.value.forEach((item) => {
    if (item.type === 1) {
      // 自研项目
      if (self.children) {
        self.children.push({
          title: item.name,
          titleAlign: "center",
          key: item.id
        })
      }
    } else if (item.type === 2) {
      // 本月分摊
      if (share.children) {
        share.children.push({
          title: item.name,
          titleAlign: "center",
          key: item.id
        })
      }
    } else if (item.type === 3) {
      // 待摊费用
      if (cost.children) {
        cost.children.push({
          title: item.name,
          titleAlign: "center",
          key: item.id
        })
      }
    }
  })
  cols.push(self, share, cost)
  return cols
})

const condition = ref<WorkTimeStatisticsCondition>({
  departmentId: "",
  dateCondition: {
    start: 0,
    end: 0
  }
})
const workTimeList = ref<WorkTime[]>([])
const getWorkTimeDataList = () => {
  getWorkTimeStatistics(condition.value).then((res) => {
    const list = res || []
    // 把list中的projectId做一下处理，如果有topParentId的，则替换为topParentId
    list.forEach((item) => {
      if (item.topParentId) item.projectId = item.topParentId
    })
    workTimeList.value = list
  })
}
type WorkTimeData = {
  name?: string
  [key: string]: number | string | undefined
}
const data = computed(() => {
  const data: WorkTimeData[] = []
  workTimeList.value.forEach((item) => {
    // 先看data中是否存在该用户
    let user = data.find((d) => d.name === item.name)
    if (!user) {
      user = {
        name: item.name || ""
      }
      data.push(user)
    }
    if (user) {
      if (!user[item.projectId]) {
        user[item.projectId] = 0
      }

      let wt = item.workTime
      if (item.startDate < dateRange.value[0] || item.endDate > dateRange.value[1]) {
        // 如果开始或结束时间落在查询时间段之外，那么就要按照天数计算比例
        // 计算出开始日期和结束日期在查询时间段内的天数
        const start = dayjs(item.startDate).isBefore(dayjs(dateRange.value[0])) ? dayjs(dateRange.value[0]) : dayjs(item.startDate)
        const end = dayjs(item.endDate).isAfter(dayjs(dateRange.value[1])) ? dayjs(dateRange.value[1]) : dayjs(item.endDate)
        const days = end.diff(start, "day") + 1

        // item的总天数
        const totalDays = dayjs(item.endDate).diff(dayjs(item.startDate), "day") + 1
        wt = (item.workTime * days) / totalDays
      }

      user[item.projectId] = ((wt + (user[item.projectId] as number) * 3600) / 3600).toFixed(2)
    }
  })
  return data
})
const search = () => {
  condition.value.departmentId = deptId.value
  condition.value.dateCondition = {
    start: dateRange.value[0],
    end: dateRange.value[1] + 86400000 - 1
  }
  getWorkTimeDataList()
}

// 项目配置抽屉
const projectListInstance = ref<ProjectWithType[]>([])
const projectConfigDrawerActived = ref(false)
const handleOpenProjectConfigDrawer = () => {
  projectListInstance.value = JSON.parse(JSON.stringify(projectList.value))
  projectConfigDrawerActived.value = true
}
const projectConfigOptions = [
  {
    label: "不参与统计",
    value: -1
  },
  {
    label: "自研项目",
    value: 1
  },
  {
    label: "本月分摊",
    value: 2
  },
  {
    label: "待摊费用",
    value: 3
  }
]
const confirmProjectConfigUpdation = () => {
  projectList.value = projectListInstance.value
  localStorage.setItem("projectList", JSON.stringify(projectList.value))
  projectConfigDrawerActived.value = false
}

// 加载页面
onMounted(() => {
  // 获取部门树
  getDepartmentList({ onlyParent: true, offset: -1, limit: -1 }).then((res) => {
    topDepartments.value = res.items.map((item) => {
      item.isLeaf = !item.childCount || item.childCount === 0
      return item
    })
    deptId.value = topDepartments.value[0]?.id
  })

  // 从localStorage中获取项目配置
  const projectConfigListFromLocalStorage = localStorage.getItem("projectList")
  if (projectConfigListFromLocalStorage) {
    projectList.value = JSON.parse(projectConfigListFromLocalStorage)
  }

  // 获取项目列表
  getAllProjectListForWorkTimeStatistics().then((res) => {
    if (res) {
      // 如果projectList中有res中不存在的项目，则删除
      for (let i = projectList.value.length - 1; i >= 0; i--) {
        const indexInRes = res.findIndex((p) => p.id === projectList.value[i].id)
        if (indexInRes === -1) {
          projectList.value.splice(i, 1)
        }
      }
      // res中如果存在projectList中不存在的项目，则添加到projectList中，type=-1，表示不参与统计；
      const projectListFromRes = res.map((item) => {
        return {
          ...item,
          type: -1
        }
      })
      projectList.value.forEach((item) => {
        const index = projectListFromRes.findIndex((p) => p.id === item.id)
        if (index > -1) {
          projectListFromRes.splice(index, 1)
        }
      })
      projectList.value = projectList.value.concat(projectListFromRes)
    }
  })
})
</script>
