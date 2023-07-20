export type OfficeConfig = {
  document: DocumentConfig
  documentType?: string
  token?: string
  type?: string
  editorConfig?: EditorConfig
}

export type DocumentConfig = {
  fileType: string
  key: string
  title: string
  url: string
  permissions?: OfficePermissions
}

export type EditorConfig = {
  mode: string
  lang: string
  callbackUrl: string
  user: OfficeUser
  customization?: unknown
}

export type OfficePermissions = {
  comment?: boolean
  commentGroups?: CommentGroups
  copy?: boolean
  deleteCommentAuthorOnly?: boolean
  download?: boolean
  edit?: boolean
  editCommentAuthorOnly?: boolean
  fillForms?: boolean
  modifyContentControl?: boolean
  modifyFilter?: boolean
  print?: boolean
  review?: boolean
  reviewGroups?: string[]
}

export type OfficeUser = {
  id: string
  name: string
  group?: string
}

export type CommentGroups = {
  edit?: string[]
  remove?: string[]
  view?: string
}
