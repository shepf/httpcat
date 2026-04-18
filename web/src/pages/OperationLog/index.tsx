import {
  CalendarOutlined,
  ClockCircleOutlined,
  CloudDownloadOutlined,
  CloudUploadOutlined,
  DeleteOutlined,
  DesktopOutlined,
  EditOutlined,
  EyeOutlined,
  FileZipOutlined,
  FolderAddOutlined,
  KeyOutlined,
  LinkOutlined,
  LoginOutlined,
  PictureOutlined,
  ReloadOutlined,
  SearchOutlined,
  SettingOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import {
  Button,
  Card,
  Col,
  DatePicker,
  Input,
  message,
  Row,
  Select,
  Space,
  Statistic,
  Table,
  Tag,
  Tooltip,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import React, { useEffect, useState } from 'react';
import { getOperationLogs, getOperationStats } from '@/services/ant-design-pro/api';

const { RangePicker } = DatePicker;

// 操作类型映射
const actionMap: Record<string, { label: string; color: string; icon: React.ReactNode }> = {
  upload:          { label: '上传文件',   color: 'green',   icon: <CloudUploadOutlined /> },
  download:        { label: '下载文件',   color: 'blue',    icon: <CloudDownloadOutlined /> },
  delete:          { label: '删除文件',   color: 'red',     icon: <DeleteOutlined /> },
  rename:          { label: '重命名',     color: 'orange',  icon: <EditOutlined /> },
  mkdir:           { label: '新建文件夹', color: 'gold',    icon: <FolderAddOutlined /> },
  preview:         { label: '预览文件',   color: 'cyan',    icon: <EyeOutlined /> },
  download_zip:    { label: '打包下载',   color: 'purple',  icon: <FileZipOutlined /> },
  image_upload:    { label: '上传图片',   color: 'lime',    icon: <PictureOutlined /> },
  image_delete:    { label: '删除图片',   color: 'volcano', icon: <DeleteOutlined /> },
  image_rename:    { label: '重命名图片', color: 'orange',  icon: <EditOutlined /> },
  image_clear:     { label: '清空图片',   color: 'magenta', icon: <DeleteOutlined /> },
  share_create:    { label: '创建分享',   color: 'geekblue',icon: <LinkOutlined /> },
  share_delete:    { label: '删除分享',   color: 'red',     icon: <DeleteOutlined /> },
  login:           { label: '用户登录',   color: 'default',  icon: <LoginOutlined /> },
  change_password: { label: '修改密码',   color: 'warning', icon: <KeyOutlined /> },
  config_update:   { label: '更新配置',   color: 'processing', icon: <SettingOutlined /> },
  restart:         { label: '重启服务',   color: 'error',   icon: <ThunderboltOutlined /> },
};

const OperationLog: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<API.OperationLogItem[]>([]);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  // 筛选条件
  const [filterAction, setFilterAction] = useState<string | undefined>();
  const [filterIP, setFilterIP] = useState('');
  const [filterDetail, setFilterDetail] = useState('');
  const [filterDateRange, setFilterDateRange] = useState<[string, string] | null>(null);

  // 统计数据
  const [stats, setStats] = useState<API.OperationStats | null>(null);
  const [statsLoading, setStatsLoading] = useState(false);

  const fetchLogs = async (page?: number, size?: number) => {
    setLoading(true);
    try {
      const params: API.OperationLogParams = {
        current: page ?? current,
        pageSize: size ?? pageSize,
        action: filterAction,
        ip: filterIP || undefined,
        detail: filterDetail || undefined,
        dateFrom: filterDateRange?.[0],
        dateTo: filterDateRange?.[1],
      };
      const res = await getOperationLogs(params);
      if (res.errorCode === 0 && res.data) {
        setData(res.data.list || []);
        setTotal(res.data.total || 0);
        setCurrent(res.data.current || 1);
      } else {
        message.error(res.msg || '获取操作日志失败');
      }
    } catch {
      message.error('获取操作日志失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchStats = async () => {
    setStatsLoading(true);
    try {
      const res = await getOperationStats();
      if (res.errorCode === 0 && res.data) {
        setStats(res.data);
      }
    } catch {
      // ignore
    } finally {
      setStatsLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
    fetchStats();
  }, []);

  const handleSearch = () => {
    setCurrent(1);
    fetchLogs(1);
  };

  const handleReset = () => {
    setFilterAction(undefined);
    setFilterIP('');
    setFilterDetail('');
    setFilterDateRange(null);
    setCurrent(1);
    // 重置后重新查询
    setTimeout(() => fetchLogs(1), 0);
  };

  const getStatusTag = (status?: number) => {
    if (!status) return null;
    if (status >= 200 && status < 300) return <Tag color="green">{status}</Tag>;
    if (status >= 300 && status < 400) return <Tag color="blue">{status}</Tag>;
    if (status >= 400 && status < 500) return <Tag color="orange">{status}</Tag>;
    return <Tag color="red">{status}</Tag>;
  };

  const getMethodTag = (method?: string) => {
    const colors: Record<string, string> = {
      GET: 'green',
      POST: 'blue',
      PUT: 'orange',
      DELETE: 'red',
      PATCH: 'purple',
    };
    return <Tag color={colors[method || ''] || 'default'}>{method}</Tag>;
  };

  const columns: ColumnsType<API.OperationLogItem> = [
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 170,
      render: (text: string) => (
        <Tooltip title={text}>
          <Space size={4}>
            <ClockCircleOutlined style={{ color: '#999' }} />
            <span>{text}</span>
          </Space>
        </Tooltip>
      ),
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 130,
      render: (action: string) => {
        const info = actionMap[action];
        if (!info) return <Tag>{action}</Tag>;
        return (
          <Tag color={info.color} icon={info.icon}>
            {info.label}
          </Tag>
        );
      },
    },
    {
      title: '详情',
      dataIndex: 'detail',
      key: 'detail',
      ellipsis: true,
      render: (text: string) => (
        <Tooltip title={text}>
          <span>{text || '-'}</span>
        </Tooltip>
      ),
    },
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 100,
      render: (text: string) => text || '-',
    },
    {
      title: 'IP',
      dataIndex: 'ip',
      key: 'ip',
      width: 140,
      render: (text: string) => (
        <Tooltip title={text}>
          <Space size={4}>
            <DesktopOutlined style={{ color: '#999' }} />
            <span style={{ fontFamily: 'monospace', fontSize: 12 }}>{text}</span>
          </Space>
        </Tooltip>
      ),
    },
    {
      title: '方法',
      dataIndex: 'method',
      key: 'method',
      width: 80,
      render: (method: string) => getMethodTag(method),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 70,
      render: (status: number) => getStatusTag(status),
    },
    {
      title: '耗时',
      dataIndex: 'latency',
      key: 'latency',
      width: 80,
      render: (latency: number) => {
        if (!latency && latency !== 0) return '-';
        if (latency < 100) return <span style={{ color: '#52c41a' }}>{latency}ms</span>;
        if (latency < 1000) return <span style={{ color: '#faad14' }}>{latency}ms</span>;
        return <span style={{ color: '#ff4d4f' }}>{latency}ms</span>;
      },
    },
  ];

  // 统计卡片中操作类型前 5
  const topActions = stats?.actionCounts?.slice(0, 5) || [];

  return (
    <PageContainer>
      {/* 统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card size="small" loading={statsLoading}>
            <Statistic
              title="今日操作"
              value={stats?.todayCount ?? 0}
              prefix={<CalendarOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card size="small" loading={statsLoading}>
            <Statistic
              title="累计操作"
              value={stats?.totalCount ?? 0}
              prefix={<ThunderboltOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col span={12}>
          <Card size="small" title="热门操作 TOP5" loading={statsLoading}>
            <Space wrap>
              {topActions.map((item) => {
                const info = actionMap[item.action];
                return (
                  <Tag
                    key={item.action}
                    color={info?.color || 'default'}
                    icon={info?.icon}
                    style={{ cursor: 'pointer' }}
                    onClick={() => {
                      setFilterAction(item.action);
                      setCurrent(1);
                      fetchLogs(1);
                    }}
                  >
                    {info?.label || item.action}: {item.count}
                  </Tag>
                );
              })}
              {topActions.length === 0 && <span style={{ color: '#999' }}>暂无数据</span>}
            </Space>
          </Card>
        </Col>
      </Row>

      {/* 筛选栏 */}
      <Card size="small" style={{ marginBottom: 16 }}>
        <Space wrap>
          <Select
            placeholder="操作类型"
            allowClear
            style={{ width: 150 }}
            value={filterAction}
            onChange={setFilterAction}
          >
            {Object.entries(actionMap).map(([key, info]) => (
              <Select.Option key={key} value={key}>
                <Space size={4}>
                  {info.icon}
                  {info.label}
                </Space>
              </Select.Option>
            ))}
          </Select>
          <Input
            placeholder="IP 地址"
            allowClear
            style={{ width: 150 }}
            value={filterIP}
            onChange={(e) => setFilterIP(e.target.value)}
            onPressEnter={handleSearch}
          />
          <Input
            placeholder="操作详情"
            allowClear
            style={{ width: 200 }}
            value={filterDetail}
            onChange={(e) => setFilterDetail(e.target.value)}
            onPressEnter={handleSearch}
          />
          <RangePicker
            onChange={(_, dateStrings) => {
              if (dateStrings[0] && dateStrings[1]) {
                setFilterDateRange([dateStrings[0], dateStrings[1]]);
              } else {
                setFilterDateRange(null);
              }
            }}
          />
          <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
            搜索
          </Button>
          <Button icon={<ReloadOutlined />} onClick={handleReset}>
            重置
          </Button>
        </Space>
      </Card>

      {/* 日志表格 */}
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        scroll={{ x: 1000 }}
        pagination={{
          current,
          pageSize,
          total,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          showTotal: (t) => `共 ${t} 条操作记录`,
          onChange: (page, size) => {
            setCurrent(page);
            setPageSize(size);
            fetchLogs(page, size);
          },
        }}
        size="middle"
      />
    </PageContainer>
  );
};

export default OperationLog;
