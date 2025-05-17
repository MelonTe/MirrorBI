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

/** 根据chart表的ID，获取上传的原始EXCEL的JSON格式数据 GET /api/chart/data */
export async function getChartData(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getChartDataParams,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: boolean }>('/api/chart/data', {
    method: 'GET',
    params: {
      ...params,
    },
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

/** 上传excel文件和目标信息，使用AI生成信息。 POST /api/chart/gen/ai */
export async function postChartGenAi(
  body: {
    /** 图表名称 */
    name: string;
    /** 分析目标 */
    goal: string;
    /** 图表类型 */
    chartType: string;
  },
  file?: File,
  options?: { [key: string]: any },
) {
  const formData = new FormData();

  if (file) {
    formData.append('file', file);
  }

  Object.keys(body).forEach((ele) => {
    const item = (body as any)[ele];

    if (item !== undefined && item !== null) {
      if (typeof item === 'object' && !(item instanceof File)) {
        if (item instanceof Array) {
          item.forEach((f) => formData.append(ele, f || ''));
        } else {
          formData.append(ele, JSON.stringify(item));
        }
      } else {
        formData.append(ele, item);
      }
    }
  });

  return request<API.Response & { data?: API.ChartGenByAiResponse }>('/api/chart/gen/ai', {
    method: 'POST',
    data: formData,
    requestType: 'form',
    ...(options || {}),
  });
}

/** 上传excel文件和目标信息，异步执行AI生成信息。 POST /api/chart/gen/ai/async */
export async function postChartGenAiAsync(
  body: {
    /** 图表名称 */
    name: string;
    /** 分析目标 */
    goal: string;
    /** 图表类型 */
    chartType: string;
  },
  file?: File,
  options?: { [key: string]: any },
) {
  const formData = new FormData();

  if (file) {
    formData.append('file', file);
  }

  Object.keys(body).forEach((ele) => {
    const item = (body as any)[ele];

    if (item !== undefined && item !== null) {
      if (typeof item === 'object' && !(item instanceof File)) {
        if (item instanceof Array) {
          item.forEach((f) => formData.append(ele, f || ''));
        } else {
          formData.append(ele, JSON.stringify(item));
        }
      } else {
        formData.append(ele, item);
      }
    }
  });

  return request<API.Response & { data?: string }>('/api/chart/gen/ai/async', {
    method: 'POST',
    data: formData,
    requestType: 'form',
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

/** 根据页数查询图表列表，是未成功分析的列表 POST /api/chart/list/page/my/no */
export async function postChartListPageMyNo(
  body: API.ChartQueryRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: API.ListChartResponse }>('/api/chart/list/page/my/no', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
