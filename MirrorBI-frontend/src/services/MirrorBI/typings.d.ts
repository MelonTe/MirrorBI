declare namespace API {
  type Chart = {
    chartData?: string;
    chartType?: string;
    createTime?: string;
    genChart?: string;
    genResult?: string;
    goal?: string;
    id?: string;
    updateTime?: string;
    userId?: string;
  };

  type ChartAddRequest = {
    /** 图表数据 */
    chartData?: string;
    /** 图表类型 */
    chartType?: string;
    /** 目标 */
    goal?: string;
  };

  type ChartEditRequest = {
    /** 图表数据 */
    chartData?: string;
    /** 图表类型 */
    chartType?: string;
    /** 目标 */
    goal?: string;
    /** 图表ID */
    id?: string;
  };

  type ChartQueryRequest = {
    /** 图表数据 */
    chartData?: string;
    /** 图表类型 */
    chartType?: string;
    /** 当前页数 */
    current?: number;
    /** 目标 */
    goal?: string;
    /** 页面大小 */
    pageSize?: number;
    /** 排序字段 */
    sortField?: string;
    /** 排序顺序（默认升序） */
    sortOrder?: string;
    /** 用户Id */
    userId?: string;
  };

  type DeleteRequest = {
    id: string;
  };

  type getChartGetParams = {
    /** 图表的ID */
    id: string;
  };

  type getUserGetParams = {
    /** 用户的ID */
    id: string;
  };

  type getUserGetVoParams = {
    /** 用户的ID */
    id: string;
  };

  type ListChartResponse = {
    /** 当前页数 */
    current?: number;
    /** 总页数 */
    pages?: number;
    /** 图表列表 */
    records?: Chart[];
    /** 页面大小 */
    size?: number;
    /** 总记录数 */
    total?: number;
  };

  type ListUserVOResponse = {
    /** 当前页数 */
    current?: number;
    /** 总页数 */
    pages?: number;
    records?: UserVO[];
    /** 页面大小 */
    size?: number;
    /** 总记录数 */
    total?: number;
  };

  type Response = {
    code?: number;
    data?: Record<string, any>;
    message?: string;
  };

  type User = {
    createTime?: string;
    id?: string;
    updateTime?: string;
    userAccount?: string;
    userAvatar?: string;
    userName?: string;
    userPassword?: string;
    userProfile?: string;
    userRole?: string;
  };

  type UserAddRequest = {
    /** 用户账号 */
    userAccount: string;
    /** 用户头像 */
    userAvatar?: string;
    /** 用户昵称 */
    userName?: string;
    /** 用户简介 */
    userProfile?: string;
    /** 用户权限 */
    userRole?: string;
  };

  type UserEditRequest = {
    /** 用户ID */
    id?: string;
    /** 用户昵称 */
    userName?: string;
    /** 用户简介 */
    userProfile?: string;
  };

  type UserLoginRequest = {
    userAccount: string;
    userPassword: string;
  };

  type UserLoginVO = {
    createTime?: string;
    id?: string;
    updateTime?: string;
    userAccount?: string;
    userAvatar?: string;
    userName?: string;
    userProfile?: string;
    userRole?: string;
  };

  type UserQueryRequest = {
    /** 当前页数 */
    current?: number;
    /** 用户ID */
    id?: string;
    /** 页面大小 */
    pageSize?: number;
    /** 排序字段 */
    sortField?: string;
    /** 排序顺序（默认升序） */
    sortOrder?: string;
    /** 用户账号 */
    userAccount?: string;
    /** 用户昵称 */
    userName?: string;
    /** 用户简介 */
    userProfile?: string;
    /** 用户权限 */
    userRole?: string;
  };

  type UserRegsiterRequest = {
    checkPassword: string;
    userAccount: string;
    userPassword: string;
  };

  type UserUpdateRequest = {
    /** 用户ID */
    id?: string;
    /** 用户头像 */
    userAvatar?: string;
    /** 用户昵称 */
    userName?: string;
    /** 用户简介 */
    userProfile?: string;
    /** 用户权限 */
    userRole?: string;
  };

  type UserVO = {
    createTime?: string;
    id?: string;
    userAccount?: string;
    userAvatar?: string;
    userName?: string;
    userProfile?: string;
    userRole?: string;
  };
}
