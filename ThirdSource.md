# 三方用户接入说明

## 三方用户接入配置

matchConfigurations 配置，用来匹配统一用户中心的用户，json中如果要使用该字段做匹配则填写，否则留空或不填写。所有填写的字段将全部进行匹配，若未匹配到用户，将自动新增用户

内容如下：
```json
{
    "username": "用户名字段",
    "realName": "真实姓名的字段",
    "mobile": "手机号字段",
    "email": "邮箱字段",
    "match": ["realName", "mobile", "email"] // 用于匹配的字段
}
```
### 钉钉
configuration 配置, 用来接入钉钉服务，内容如下：
```json
{
    "appKey": "钉钉应用的appKey",
    "appSecret": "钉钉应用的appSecret",
    "agentId": "钉钉应用的agentId",
    "callbackToken": "钉钉应用的callbackToken",
    "callbackAesKey": "钉钉应用的callbackAesKey"
}
```
根据文档钉钉返回的用户json应该为如下字符串
```json
{
    "errcode":"0",
    "result":{
        "extension":"{\"爱好\":\"旅游\",\"年龄\":\"24\"}",
        "unionid":"z21HjQliSzpw0YWCNxmxxxxx",
        "boss":"true",
        "role_list":{
            "group_name":"职务",
            "name":"总监",
            "id":"100"
        },
        "exclusive_account":false,
        "manager_userid":"manager240",
        "admin":"true",
        "remark":"备注备注",
        "title":"技术总监",
        "hired_date":"1597573616828",
        "userid":"zhangsan",
        "work_place":"未来park",
        "dept_order_list":{
            "dept_id":"2",
            "order":"1"
        },
        "real_authed":"true",
        "dept_id_list":"[2,3,4]",
        "job_number":"4",
        "email":"test@xxx.com",
        "leader_in_dept":{
            "leader":"true",
            "dept_id":"2"
        },
        "mobile":"18513027676",
        "active":"true",
        "telephone":"010-86123456-2345",
        "avatar":"xxx",
        "hide_mobile":"false",
        "senior":"true",
        "name":"张三",
        "union_emp_ext":{
            "union_emp_map_list":{
                "userid":"5000",
                "corp_id":"dingxxx"
            },
            "userid":"500",
            "corp_id":"dingxxx"
        },
        "state_code":"86"
    },
    "errmsg":"ok"
}
```