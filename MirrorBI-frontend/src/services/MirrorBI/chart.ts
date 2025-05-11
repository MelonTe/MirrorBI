// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 添加一个图表 POST /api/chart/add */
export async function postChartAdd(body: API.ChartAddRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: string }>('/api/chart/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 添加一个图表 POST /api/chart/delete */
export async function postChartOpenApiDelete(
  body: API.DeleteRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: boolean }>('/api/chart/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 编辑图表 POST /api/chart/edit */
export async function postChartEdit(body: API.ChartEditRequest, options?: { [key: string]: any }) {
  return request<API.Response & { data?: boolean }>('/api/chart/edit', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取一个图表 GET /api/chart/get */
export async function getChartGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getChartGetParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.Chart }>('/api/chart/get', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 根据页数查询图表列表 POST /api/chart/list/page */
export async function postChartListPage(
  body: API.ChartQueryRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.ListChartResponse }>('/api/chart/list/page', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 根据页数查询图表列表 POST /api/chart/list/page/my */
export async function postChartListPageMy(
  body: API.ChartQueryRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.ListChartResponse }>('/api/chart/list/page/my', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
