import { defineStore } from "pinia"
import { Stage } from "@/types/project"

export const useProjectStore = defineStore("project", {
    state: () :{
        stageList: Stage[]
    } => ({
        stageList: []
    }),
    getters: {
        stageListWithNone: (state) => {
            return [{id: "", name: "无", orderNum: 0, status: 1}, ...state.stageList]
        }
    }
})