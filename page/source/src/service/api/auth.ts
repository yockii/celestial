import {request} from "../request";
import {sm2} from "sm-crypto";
import type { User} from '../../store/user'

interface LoginResponse {
    token: string;
    user: User;
}
/**
 * 登录
 * @param username - 用户名
 * @param password - 密码(需要加密)
 */
export function login(username: string, password: string) {
    const encryptedPassword = "04" + sm2.doEncrypt(password, import.meta.env.VITE_SM2_PK)
    return request.post<LoginResponse>('/login', {
            username,
            password: encryptedPassword,
        }
    );
}