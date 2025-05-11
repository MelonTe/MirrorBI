// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 创建一个用户「管理员」 默认密码为12345678 POST /api/user/add */
export async function postUserAdd(body: API.UserAddRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: string }>('/api/user/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 根据ID软删除用户「管理员」 POST /api/user/delete */
export async function postUserOpenApiDelete(
  body: API.DeleteRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: boolean }>('/api/user/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新用户个人资料 若用户不存在，则返回失败 POST /api/user/edit */
export async function postUserEdit(body: API.UserEditRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: boolean }>('/api/user/edit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 根据ID获取用户「管理员」 GET /api/user/get */
export async function getUserGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getUserGetParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.User }>('/api/user/get', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取登录的用户信息 GET /api/user/get/login */
export async function getUserGetLogin(options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.UserLoginVO }>('/api/user/get/login', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 根据ID获取简略信息用户 GET /api/user/get/vo */
export async function getUserGetVo(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getUserGetVoParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.UserVO }>('/api/user/get/vo', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 分页获取一系列用户信息「管理员」 根据用户关键信息进行模糊查询 POST /api/user/list/page/vo */
export async function postUserListPageVo(
  body: API.UserQueryRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.ListUserVOResponse }>('/api/user/list/page/vo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户登录 根据账号密码进行登录 POST /api/user/login */
export async function postUserLogin(body: API.UserLoginRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: API.UserLoginVO }>('/api/user/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 执行用户注销（退出） POST /api/user/logout */
export async function postUserLogout(options?: { [key: string]: any }) {
  return request<API.Response & { data?: boolean }>('/api/user/logout', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 注册用户 根据账号密码进行注册 POST /api/user/register */
export async function postUserRegister(
  body: API.UserRegsiterRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: string }>('/api/user/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新用户信息「管理员」 若用户不存在，则返回失败 POST /api/user/update */
export async function postUserUpdate(
  body: API.UserUpdateRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: boolean }>('/api/user/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
