
export type User = {
    id: string;
    username: string;
    realName: string;
    password: string;
    email: string;
    mobile: string;
    status: number;
    createTime: number;
}

export type UserCondition = {
    username: string;
    realName: string;
    email: string;
    mobile: string;
    status: number;
    offset: number;
    limit: number;
    orderBy?: string;
}

export type LoginResponse = {
    token: string;
    user: User;
}

export type Role = {
    id: string;
    name: string;
    desc: string;
    type: number;
    style?: string;
    dataPermission: number;
    status: number;
    createTime?: number;
}

export type RoleCondition = {
    name: string;
    type: number;
    dataPermission: number;
    status: number;
    offset: number;
    limit: number;
    orderBy?: string;
}