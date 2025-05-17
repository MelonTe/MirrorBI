import { getChartData, getChartGet } from '@/services/MirrorBI/chart';
import { history, useParams } from '@umijs/max';
import { Button, Card, Descriptions, Empty, message, Spin, Table, Tag, Typography } from 'antd';
import ReactECharts from 'echarts-for-react';
import JSON5 from 'json5';
import React, { useEffect, useState } from 'react';
import './index.less';

const { Title, Paragraph } = Typography;

const statusColorMap: Record<string, string> = {
  succeed: 'success',
  wait: 'processing',
  running: 'processing',
  failed: 'error',
};

const statusTextMap: Record<string, string> = {
  succeed: '执行成功',
  wait: '等待执行',
  running: '正在执行',
  failed: '执行失败',
};

const ChartDetail: React.FC = () => {
  const { id } = useParams();
  const [loading, setLoading] = useState(true);
  const [chart, setChart] = useState<API.Chart | null>(null);
  const [error, setError] = useState('');
  const [showRaw, setShowRaw] = useState(false);
  const [rawLoading, setRawLoading] = useState(false);
  const [rawData, setRawData] = useState<any[]>([]);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    getChartGet({ id })
      .then((res) => {
        if (res.code === 0 && res.data) {
          setChart(res.data);
        } else {
          setError(res.message || '未找到图表');
        }
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [id]);

  const handleShowRaw = async () => {
    if (!id) return;
    setRawLoading(true);
    try {
      const res = await getChartData({ id });
      if (res.code === 0 && Array.isArray(res.data)) {
        setRawData(res.data);
        setShowRaw(true);
      } else {
        message.error(res.message || '获取原始数据失败');
      }
    } catch (e: any) {
      message.error(e.message || '获取原始数据失败');
    }
    setRawLoading(false);
  };

  let chartOption: any = null;
  try {
    chartOption = chart?.genChart ? JSON5.parse(chart.genChart) : null;
  } catch {
    chartOption = null;
  }

  return (
    <div className="chart-detail-container">
      <Card className="detail-card" bordered={false}>
        <Button className="back-btn" onClick={() => history.back()}>
          &lt; 返回
        </Button>
        {loading ? (
          <div className="loading-area">
            <Spin size="large" />
          </div>
        ) : error ? (
          <Empty description={error} />
        ) : chart ? (
          <>
            <Title level={2} className="chart-title">
              {chart.name || '未命名图表'}
            </Title>
            <Descriptions column={2} className="chart-descs">
              <Descriptions.Item label="分析目标">{chart.goal || '-'}</Descriptions.Item>
              <Descriptions.Item label="图表类型">{chart.chartType || '-'}</Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {chart.createTime?.substring(0, 10) || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="执行状态">
                <Tag color={statusColorMap[chart.status || 'wait']}>
                  {statusTextMap[chart.status || 'wait']}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="执行信息" span={2}>
                {chart.execMessage || '-'}
              </Descriptions.Item>
            </Descriptions>
            <div className="chart-section">
              {chartOption ? (
                <ReactECharts option={chartOption} style={{ height: 360, width: '100%' }} />
              ) : (
                <Empty description="图表预览不可用" />
              )}
            </div>
            <div className="ai-result-section">
              <Title level={4}>AI 分析结论</Title>
              <Paragraph className="ai-result-text">{chart.genResult || '暂无分析结论'}</Paragraph>
            </div>
            <div className="raw-data-section">
              <Button
                type="primary"
                onClick={handleShowRaw}
                loading={rawLoading}
                disabled={showRaw}
              >
                {showRaw ? '已展示原始数据' : '点击展示原始数据'}
              </Button>
              {showRaw && (
                <div className="raw-table-area">
                  <Table
                    dataSource={rawData}
                    columns={
                      rawData[0]
                        ? Object.keys(rawData[0]).map((key) => ({ title: key, dataIndex: key }))
                        : []
                    }
                    rowKey={(_, idx) => String(idx)}
                    pagination={{ pageSize: 10 }}
                    bordered
                    size="middle"
                  />
                </div>
              )}
            </div>
          </>
        ) : null}
      </Card>
    </div>
  );
};

export default ChartDetail;
