import { postChartGenAiAsync } from '@/services/MirrorBI/chart';
import {
  AreaChartOutlined,
  BarChartOutlined,
  LineChartOutlined,
  PieChartOutlined,
  RadarChartOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import { useNavigate } from '@umijs/max';
import {
  Button,
  Card,
  Col,
  Divider,
  Form,
  Input,
  message,
  Row,
  Select,
  Space,
  Spin,
  Typography,
  Upload,
} from 'antd';
import TextArea from 'antd/es/input/TextArea';
import { motion } from 'framer-motion';
import React, { useState } from 'react';
import './index.less';

const { Title, Paragraph } = Typography;

const ChartTypeOption = ({ icon, label }: { icon: React.ReactNode; label: string }) => (
  <Space>
    {icon}
    <span>{label}</span>
  </Space>
);

const AddChartAsync: React.FC = () => {
  const [option, setOption] = useState<any>();
  const [chart, setChart] = useState<any>();
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [analyzing, setAnalyzing] = useState<boolean>(false);
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const onFinish = async (values: any) => {
    if (submitting) return;
    setSubmitting(true);
    setAnalyzing(true);
    const params = { ...values, file: undefined };
    try {
      const res = await postChartGenAiAsync(params, values.file?.file, {});
      if (res.code === 0) {
        message.success('提交成功，AI分析任务已开始！');
        navigate('/my_chart');
      } else {
        message.error('提交失败：' + res.message);
      }
    } catch (e: any) {
      message.error('提交失败：' + e.message);
    }
    setSubmitting(false);
    setAnalyzing(false);
  };

  const resetForm = () => {
    form.resetFields();
    setOption(undefined);
    setChart(undefined);
  };

  return (
    <div className="add-chart-container">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Card className="form-card" bordered={false}>
          <Title level={2} className="page-title">
            <BarChartOutlined /> 智能数据分析（异步）
          </Title>
          <Paragraph className="page-description">
            上传您的Excel文件，描述您的分析需求，AI将异步为您生成数据可视化图表。提交后可在"我的图表"中查看进度。
          </Paragraph>
          <Divider />

          <Form
            form={form}
            name="addChartAsync"
            onFinish={onFinish}
            initialValues={{}}
            layout="vertical"
            className="analysis-form"
          >
            <Row gutter={24}>
              <Col xs={24} lg={12}>
                <Form.Item
                  name="goal"
                  label="分析目标"
                  rules={[{ required: true, message: '请输入分析目标!' }]}
                >
                  <TextArea
                    placeholder="请输入您的分析需求，比如：分析网站用户的增长情况"
                    autoSize={{ minRows: 3, maxRows: 6 }}
                    className="custom-textarea"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} lg={12}>
                <Form.Item name="name" label="图表名称">
                  <Input placeholder="请输入图表名称" className="custom-input" />
                </Form.Item>

                <Form.Item name="chartType" label="图表类型">
                  <Select
                    placeholder="请选择图表类型"
                    className="chart-type-select"
                    options={[
                      {
                        value: '折线图',
                        label: <ChartTypeOption icon={<LineChartOutlined />} label="折线图" />,
                      },
                      {
                        value: '柱状图',
                        label: <ChartTypeOption icon={<BarChartOutlined />} label="柱状图" />,
                      },
                      {
                        value: '堆叠图',
                        label: <ChartTypeOption icon={<AreaChartOutlined />} label="堆叠图" />,
                      },
                      {
                        value: '饼图',
                        label: <ChartTypeOption icon={<PieChartOutlined />} label="饼图" />,
                      },
                      {
                        value: '雷达图',
                        label: <ChartTypeOption icon={<RadarChartOutlined />} label="雷达图" />,
                      },
                    ]}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item name="file" label="原始数据" className="upload-container">
              <Upload.Dragger
                name="file"
                accept=".xlsx,.xls,.csv"
                beforeUpload={() => false}
                maxCount={1}
              >
                <p className="ant-upload-drag-icon">
                  <UploadOutlined />
                </p>
                <p className="ant-upload-text">点击或拖拽Excel文件到此区域上传</p>
                <p className="ant-upload-hint">支持 .xlsx, .xls 格式文件</p>
              </Upload.Dragger>
            </Form.Item>

            <Form.Item className="form-actions">
              <Space>
                <Button
                  type="primary"
                  htmlType="submit"
                  loading={submitting}
                  disabled={submitting}
                  className="analysis-btn"
                  size="large"
                >
                  开始智能分析
                </Button>
                <Button htmlType="button" onClick={resetForm} size="large">
                  重置
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </Card>
      </motion.div>

      {analyzing && !chart && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="loading-container">
          <Spin size="large" />
          <p className="loading-text">正在分析数据，请稍候...</p>
        </motion.div>
      )}
    </div>
  );
};

export default AddChartAsync;
