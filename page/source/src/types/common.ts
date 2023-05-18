
export type Result<T> = {
    code: number;
    msg: string;
    data: T;
}

export type Paginate<T> = {
    total: number;
    items: T[];
}