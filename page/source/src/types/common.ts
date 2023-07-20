export type Result<T> = {
  code: number
  msg: string
  data: T
}

export type Paginate<T> = {
  total: number
  offset: number
  limit: number
  items: T[]
}

export type Condition = {
  offset?: number
  limit?: number
  orderBy?: string
}
