import { CopyOutlined, DeleteOutlined, FileOutlined, LinkOutlined, ReloadOutlined, ShareAltOutlined, CloudDownloadOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Button, Card, Col, message, Popconfirm, Row, Space, Statistic, Table, Tag, Tooltip, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useEffect, useState } from 'react';
import { deleteShare, getShareStats, getShareConfig, listShares } from '@/services/ant-design-pro/api';

const ShareManage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<API.ShareItem[]>([]);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 });
  const [stats, setStats] = useState<API.ShareStats>({});
  const [config, setConfig] = useState<API.ShareConfig>({});
  const [statsLoading, setStatsLoading] = useState(false);

  const fetchData = async (page?: number, pageSize?: number) => {
    setLoading(true);
    try {
      const res = await listShares({
        current: page || pagination.current,
        pageSize: pageSize || pagination.pageSize,
      });
      if (res.errorCode === 0 && res.data) {
        setData(res.data.list || []);
        setPagination({
          current: res.data.current || 1,
          pageSize: res.data.pageSize || 20,
          total: res.data.total || 0,
        });
      }
    } catch {
      message.error('获取分享列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchStats = async () => {
    setStatsLoading(true);
    try {
      const [statsRes, configRes] = await Promise.all([getShareStats(), getShareConfig()]);
      if (statsRes.errorCode === 0 && statsRes.data) {
        setStats(statsRes.data);
      }
      if (configRes.errorCode === 0 && configRes.data) {
        setConfig(configRes.data);
      }
    } catch {
      // ignore
    } finally {
      setStatsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    fetchStats();
  }, []);

  const handleDelete = async (code: string) => {
    try {
      const res = await deleteShare(code);
      if (res.errorCode === 0) {
        message.success('已取消分享');
        fetchData();
        fetchStats();
      } else {
        message.error(res.msg || '操作失败');
      }
    } catch {
      message.error('操作失败');
    }
  };

  const copyShareInfo = (record: API.ShareItem) => {
    const url = `${window.location.origin}/s/${record.shareCode}`;
    const text = record.extractCode
      ? `链接: ${url}\n提取码: ${record.extractCode}`
      : `链接: ${url}`;
    navigator.clipboard.writeText(text).then(() => {
      message.success('已复制分享信息');
    }).catch(() => {
      message.error('复制失败');
    });
  };

  const getStatusTag = (record: API.ShareItem) => {
    if (!record.isActive) {
      return <Tag color="default">已取消</Tag>;
    }
    if (record.expireAt && new Date(record.expireAt) < new Date()) {
      return <Tag color="red">已过期</Tag>;
    }
    if (record.maxDownloads && record.maxDownloads > 0 && (record.curDownloads || 0) >= record.maxDownloads) {
      return <Tag color="orange">已达上限</Tag>;
    }
    return <Tag color="green">有效</Tag>;
  };

  const columns: ColumnsType<API.ShareItem> = [
    {
      title: '文件名',
      dataIndex: 'fileName',
      key: 'fileName',
      ellipsis: true,
      render: (text: string, record) => (
        <Tooltip title={text}>
          <Space>
            <FileOutlined style={{ color: record.fileType === 'image' ? '#52c41a' : '#1890ff' }} />
            <span>{text}</span>
            {record.fileType === 'image' && <Tag color="green" style={{ fontSize: 11 }}>图片</Tag>}
          </Space>
        </Tooltip>
      ),
    },
    {
      title: '分享码',
      dataIndex: 'shareCode',
      key: 'shareCode',
      width: 120,
      render: (text: string) => (
        <Tag color="blue" style={{ fontFamily: 'monospace' }}>{text}</Tag>
      ),
    },
    {
      title: '提取码',
      dataIndex: 'extractCode',
      key: 'extractCode',
      width: 100,
      render: (text: string) => text ? <Tag color="gold">{text}</Tag> : <span style={{ color: '#999' }}>无</span>,
    },
    {
      title: '下载次数',
      key: 'downloads',
      width: 120,
      render: (_, record) => {
        const cur = record.curDownloads || 0;
        const max = record.maxDownloads || 0;
        return max > 0 ? `${cur} / ${max}` : `${cur} / 不限`;
      },
    },
    {
      title: '过期时间',
      dataIndex: 'expireAt',
      key: 'expireAt',
      width: 170,
      render: (text: string) => text ? new Date(text).toLocaleString() : '永不过期',
    },
    {
      title: '状态',
      key: 'status',
      width: 90,
      render: (_, record) => getStatusTag(record),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 170,
      render: (text: string) => text ? new Date(text).toLocaleString() : '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 140,
      render: (_, record) => (
        <Space>
          <Tooltip title="复制分享信息">
            <Button
              type="link"
              size="small"
              icon={<CopyOutlined />}
              onClick={() => copyShareInfo(record)}
            />
          </Tooltip>
          {record.isActive && (
            <Popconfirm
              title="确定要取消此分享吗？"
              onConfirm={() => handleDelete(record.shareCode!)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                取消
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      {/* 统计概览卡片 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card loading={statsLoading} size="small" hoverable>
            <Statistic
              title="总分享数"
              value={stats.totalShares || 0}
              prefix={<ShareAltOutlined style={{ color: '#1890ff' }} />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={statsLoading} size="small" hoverable>
            <Statistic
              title="有效分享"
              value={stats.activeShares || 0}
              prefix={<LinkOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={statsLoading} size="small" hoverable>
            <Statistic
              title="已失效"
              value={stats.expiredShares || 0}
              prefix={<DeleteOutlined style={{ color: '#ff4d4f' }} />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={statsLoading} size="small" hoverable>
            <Statistic
              title="总下载量"
              value={stats.totalDownloads || 0}
              prefix={<CloudDownloadOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Space>
          {config.anonymousAccess !== undefined && (
            <Tag color={config.anonymousAccess ? 'green' : 'orange'}>
              {config.anonymousAccess ? '✓ 允许匿名访问' : '✗ 需登录访问'}
            </Tag>
          )}
        </Space>
        <Button icon={<ReloadOutlined />} onClick={() => { fetchData(); fetchStats(); }}>
          刷新
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50'],
          showTotal: (total) => `共 ${total} 条分享`,
          onChange: (page, pageSize) => fetchData(page, pageSize),
        }}
        size="middle"
      />
    </PageContainer>
  );
};

export default ShareManage;
