import { postChartListPageMy } from '@/services/MirrorBI/chart';
import { getUserGetLogin } from '@/services/MirrorBI/user';
import {
  AreaChartOutlined,
  BarChartOutlined,
  DeleteOutlined,
  EditOutlined,
  LineChartOutlined,
  PieChartOutlined,
  RadarChartOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import {
  Button,
  Card,
  Col,
  Empty,
  Input,
  message,
  Pagination,
  Popconfirm,
  Row,
  Spin,
  Tooltip,
  Typography,
} from 'antd';
import ReactECharts from 'echarts-for-react';
import { motion } from 'framer-motion';
import JSON5 from 'json5';
import React, { useEffect, useRef, useState } from 'react';
import './index.less';

const { Title, Paragraph } = Typography;
const { Meta } = Card;

/**
 * 获取图表图标
 * @param chartType 图表类型
 */
const getChartTypeIcon = (chartType?: string) => {
  if (!chartType) {
    return <BarChartOutlined />;
  }
  switch (chartType) {
    case '折线图':
      return <LineChartOutlined />;
    case '柱状图':
      return <BarChartOutlined />;
    case '饼图':
      return <PieChartOutlined />;
    case '雷达图':
      return <RadarChartOutlined />;
    case '堆叠图':
      return <AreaChartOutlined />;
    default:
      return <BarChartOutlined />;
  }
};

/**
 * 我的图表页面
 * @constructor
 */
const MyChartPage: React.FC = () => {
  // 初始化查询参数
  const initSearchParams = {
    pageSize: 6,
    current: 1,
  };

  const [searchParams, setSearchParams] = useState<API.ChartQueryRequest>({ ...initSearchParams });
  const [chartList, setChartList] = useState<API.Chart[]>();
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [searchName, setSearchName] = useState<string>('');

  // 用于跟踪图表是否已经完成初始化
  const [chartsInitialized, setChartsInitialized] = useState<boolean>(false);
  const chartsRef = useRef<Record<string, any>>({});

  // 添加用户信息状态
  const [currentUser, setCurrentUser] = useState<API.UserLoginVO>();

  // 处理图表实例事件
  const onChartReady = (instance: any, chartId: string) => {
    chartsRef.current[chartId] = instance;
    instance.resize();
  };

  // 在组件挂载和窗口大小变化时重新调整所有图表大小
  useEffect(() => {
    const handleResize = () => {
      if (chartsRef.current) {
        Object.values(chartsRef.current).forEach((chart: any) => {
          if (chart && chart.resize) {
            chart.resize();
          }
        });
      }
    };

    // 等待图表容器完全渲染后初始化图表
    if (chartList && chartList.length > 0 && !chartsInitialized) {
      setTimeout(() => {
        handleResize();
        setChartsInitialized(true);
      }, 300);
    }

    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [chartList, chartsInitialized]);

  // 获取当前登录用户信息
  const loadUserInfo = async () => {
    try {
      const res = await getUserGetLogin();
      if (res.data) {
        setCurrentUser(res.data);
      } else {
        message.error('获取用户信息失败');
      }
    } catch (e: any) {
      message.error('获取用户信息失败,' + e.message);
    }
  };

  // 首次加载时获取用户信息
  useEffect(() => {
    loadUserInfo();
  }, []);

  // 加载数据
  const loadData = async () => {
    setLoading(true);
    try {
      const res = await postChartListPageMy(searchParams);

      if (res.data) {
        setChartList(res.data.records ?? []);
        setTotal(res.data.total ?? 0);
      } else {
        message.error('获取我的图表失败');
      }
    } catch (e: any) {
      message.error('获取我的图表失败,' + e.message);
    } finally {
      setLoading(false);
    }
  };

  // 搜索方法
  const handleSearch = () => {
    // 设置搜索参数，重置页码到第一页
    setSearchParams({
      ...searchParams,
      name: searchName,
      current: 1,
    });
  };

  // 重置搜索
  const handleReset = () => {
    setSearchName('');
    setSearchParams(initSearchParams);
  };

  // 处理分页变化
  const handlePageChange = (page: number, pageSize: number) => {
    setSearchParams({
      ...searchParams,
      current: page,
      pageSize,
    });
  };

  // 数据变化时重新加载
  useEffect(() => {
    loadData();
  }, [searchParams]);

  // 页面动画变体
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.5,
      },
    },
  };

  // 处理图表配置，添加通用的样式优化
  const optimizeChartOption = (chartOption: any) => {
    if (!chartOption) return null;

    // 深拷贝以避免修改原对象
    const optimizedOption = JSON.parse(JSON.stringify(chartOption));

    // 调整标题配置
    if (optimizedOption.title) {
      optimizedOption.title = {
        ...optimizedOption.title,
        left: 'center',
        textStyle: {
          ...optimizedOption.title.textStyle,
          fontSize: 14,
        },
        padding: [10, 0],
      };
    }

    // 调整图例配置
    if (optimizedOption.legend) {
      optimizedOption.legend = {
        ...optimizedOption.legend,
        textStyle: {
          fontSize: 12,
        },
        bottom: 0,
        padding: [0, 0, 5, 0],
      };
    }

    // 确保图表有足够的边距
    optimizedOption.grid = {
      containLabel: true,
      left: '3%',
      right: '4%',
      bottom: '8%',
      top: optimizedOption.title ? '15%' : '8%',
      ...optimizedOption.grid,
    };

    // 确保文本不会被截断
    optimizedOption.tooltip = {
      confine: true,
      ...optimizedOption.tooltip,
    };

    return optimizedOption;
  };

  return (
    <div className="my-chart-container">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="page-header"
      >
        <div className="header-content">
          <div className="title-section">
            <Title level={2} className="page-title">
              我的图表
            </Title>
            <Paragraph className="page-description">展示您创建的所有数据可视化图表</Paragraph>
          </div>
          <div className="search-section">
            <Input
              placeholder="搜索图表名称"
              value={searchName}
              onChange={(e) => setSearchName(e.target.value)}
              onPressEnter={handleSearch}
              prefix={<SearchOutlined />}
              allowClear
              className="search-input"
            />
            <Button type="primary" onClick={handleSearch} className="search-button">
              搜索
            </Button>
            <Button onClick={handleReset} className="reset-button">
              重置
            </Button>
          </div>
        </div>
      </motion.div>

      <motion.div
        className="main-content"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        {loading ? (
          <div className="loading-container">
            <Spin size="large" />
            <p className="loading-text">加载中...</p>
          </div>
        ) : (
          <>
            {!chartList || chartList.length === 0 ? (
              <Empty description="暂无图表数据" className="empty-container" />
            ) : (
              <motion.div
                variants={containerVariants}
                initial="hidden"
                animate="visible"
                className="chart-grid"
              >
                <Row gutter={[40, 40]}>
                  {chartList.map((chart) => {
                    // 尝试解析图表配置
                    let chartOption;
                    try {
                      chartOption = chart.genChart ? JSON5.parse(chart.genChart) : null;
                      // 应用优化配置
                      chartOption = optimizeChartOption(chartOption);
                    } catch (e) {
                      chartOption = null;
                    }

                    return (
                      <Col xs={24} sm={24} md={24} lg={12} key={chart.id}>
                        <motion.div variants={itemVariants}>
                          <Card
                            className="chart-card"
                            hoverable
                            cover={
                              <div className="chart-preview">
                                {chartOption ? (
                                  <ReactECharts
                                    option={chartOption}
                                    style={{ height: '300px', width: '100%' }}
                                    opts={{
                                      renderer: 'canvas',
                                      devicePixelRatio: 2, // 提高清晰度
                                    }}
                                    notMerge={true}
                                    lazyUpdate={true}
                                    className="chart-instance"
                                    onChartReady={(instance) =>
                                      onChartReady(instance, chart.id || '')
                                    }
                                    onEvents={{
                                      finished: () => console.log('Chart render finished'),
                                    }}
                                  />
                                ) : (
                                  <div className="chart-fallback">
                                    {getChartTypeIcon(chart.chartType)}
                                    <p>图表预览不可用</p>
                                  </div>
                                )}
                              </div>
                            }
                            actions={[
                              <Tooltip title="编辑" key="edit">
                                <EditOutlined />
                              </Tooltip>,
                              <Tooltip title="删除" key="delete">
                                <Popconfirm
                                  title="确定要删除这个图表吗？"
                                  okText="确定"
                                  cancelText="取消"
                                >
                                  <DeleteOutlined />
                                </Popconfirm>
                              </Tooltip>,
                            ]}
                          >
                            <Meta
                              title={chart.name || '未命名图表'}
                              description={
                                <div className="chart-info">
                                  <div className="chart-type">
                                    {getChartTypeIcon(chart.chartType)}
                                    <span>{chart.chartType || '未知类型'}</span>
                                  </div>
                                  <div className="chart-goal">{chart.goal || '无分析目标'}</div>
                                  <div className="chart-author-time">
                                    <div className="user-info">
                                      {currentUser?.userAvatar ? (
                                        <img
                                          src={currentUser.userAvatar}
                                          alt="用户头像"
                                          className="user-avatar"
                                        />
                                      ) : (
                                        <div className="default-avatar">
                                          {currentUser?.userName?.substring(0, 1) || 'U'}
                                        </div>
                                      )}
                                      <span className="user-name">
                                        {currentUser?.userName ||
                                          currentUser?.userAccount ||
                                          '未知用户'}
                                      </span>
                                    </div>
                                    <span className="chart-time">
                                      创建于 {chart.createTime?.substring(0, 10) || '未知时间'}
                                    </span>
                                  </div>
                                </div>
                              }
                            />
                          </Card>
                        </motion.div>
                      </Col>
                    );
                  })}
                </Row>
              </motion.div>
            )}

            {total > 0 && (
              <div className="pagination-section">
                <Pagination
                  current={searchParams.current}
                  pageSize={searchParams.pageSize}
                  total={total}
                  showSizeChanger
                  showQuickJumper
                  showTotal={(t) => `共 ${t} 个图表`}
                  onChange={handlePageChange}
                  pageSizeOptions={['3', '6', '9', '12']}
                />
              </div>
            )}
          </>
        )}
      </motion.div>
    </div>
  );
};

export default MyChartPage;
