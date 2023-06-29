import { request } from "../request"
import { sm2 } from "sm-crypto"
import { LoginResponse } from "@/types/user"

/**
 * 登录
 * @param username - 用户名
 * @param password - 密码(需要加密)
 */
export function login(username: string, password: string) {
  return new Promise<LoginResponse>((resolve, reject) => {
    getPublicKey().then((res) => {
      const encryptedPassword = "04" + sm2.doEncrypt(password, res)
      resolve(
        request.post<LoginResponse>("/login", {
          username,
          password: encryptedPassword
        })
      )
    })
  })
}

/**
 * 钉钉中打开免登
 * @param code - 钉钉免登授权码
 * @param thirdSourceId - 三方源ID
 */
export function loginInDingTalk(thirdSourceId: string, code: string) {
  return request.post<LoginResponse>("/loginInDingtalk", {
    code,
    thirdSourceId
  })
}

/**
 * 获取公钥
 */
export function getPublicKey() {
  return request.get<string>("/pk")
}
