import { CopyOutlined, DownloadOutlined, FileOutlined, FolderOutlined, ReloadOutlined, ShareAltOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Button, Form, Input, InputNumber, message, Modal, Select, Space, Switch, Table, Tag, Tooltip, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useEffect, useState } from 'react';
import { createShare, listFiles } from '@/services/ant-design-pro/api';

const { Paragraph } = Typography;

const FileList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<API.FileItem[]>([]);
  const [searchText, setSearchText] = useState('');
  const [shareModalVisible, setShareModalVisible] = useState(false);
  const [shareResultVisible, setShareResultVisible] = useState(false);
  const [shareResult, setShareResult] = useState<API.CreateShareResult | null>(null);
  const [currentFile, setCurrentFile] = useState<API.FileItem | null>(null);
  const [shareLoading, setShareLoading] = useState(false);
  const [useExtractCode, setUseExtractCode] = useState(true);
  const [form] = Form.useForm();

  const fetchData = async () => {
    setLoading(true);
    try {
      const res = await listFiles();
      setData(res.data || []);
    } catch (error) {
      message.error('获取文件列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleDownload = (fileName: string) => {
    window.open(`/api/v1/file/download?filename=${encodeURIComponent(fileName)}`);
  };

  const handleShare = (record: API.FileItem) => {
    setCurrentFile(record);
    setUseExtractCode(true);
    form.resetFields();
    form.setFieldsValue({ expireHours: 24, maxDownloads: 0 });
    setShareModalVisible(true);
  };

  const handleCreateShare = async () => {
    if (!currentFile) return;
    setShareLoading(true);
    try {
      const values = await form.validateFields();
      const params: API.CreateShareParams = {
        filePath: currentFile.FileName,
        fileName: currentFile.FileName,
        fileType: 'file',
        expireHours: values.expireHours || 0,
        maxDownloads: values.maxDownloads || 0,
        extractCode: useExtractCode ? 'auto' : '',
      };
      const res = await createShare(params);
      if (res.errorCode === 0 && res.data) {
        setShareResult(res.data);
        setShareModalVisible(false);
        setShareResultVisible(true);
        message.success('分享创建成功');
      } else {
        message.error(res.msg || '创建分享失败');
      }
    } catch {
      message.error('创建分享失败');
    } finally {
      setShareLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).then(() => {
      message.success('已复制到剪贴板');
    }).catch(() => {
      message.error('复制失败，请手动复制');
    });
  };

  const getFileIcon = (fileName: string) => {
    const ext = fileName.split('.').pop()?.toLowerCase();
    const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico'];
    if (imageExts.includes(ext || '')) return <FileOutlined style={{ color: '#52c41a' }} />;
    const archiveExts = ['zip', 'tar', 'gz', 'rar', '7z', 'bz2', 'xz'];
    if (archiveExts.includes(ext || '')) return <FolderOutlined style={{ color: '#faad14' }} />;
    return <FileOutlined style={{ color: '#1890ff' }} />;
  };

  const getFileTag = (fileName: string) => {
    const ext = fileName.split('.').pop()?.toLowerCase() || '';
    const colorMap: Record<string, string> = {
      pdf: 'red', doc: 'blue', docx: 'blue', xls: 'green', xlsx: 'green',
      zip: 'orange', tar: 'orange', gz: 'orange', rar: 'orange',
      jpg: 'cyan', jpeg: 'cyan', png: 'cyan', gif: 'cyan', webp: 'cyan', svg: 'cyan',
      mp4: 'purple', mp3: 'purple', avi: 'purple', mkv: 'purple',
      txt: 'default', log: 'default', json: 'geekblue', xml: 'geekblue',
      sh: 'volcano', py: 'volcano', go: 'volcano', js: 'volcano',
    };
    return ext ? <Tag color={colorMap[ext] || 'default'}>{ext.toUpperCase()}</Tag> : null;
  };

  const filteredData = searchText
    ? data.filter((item) => item.FileName.toLowerCase().includes(searchText.toLowerCase()))
    : data;

  const columns: ColumnsType<API.FileItem> = [
    {
      title: '文件名',
      dataIndex: 'FileName',
      key: 'FileName',
      ellipsis: true,
      sorter: (a, b) => a.FileName.localeCompare(b.FileName),
      render: (text: string) => (
        <Space>
          {getFileIcon(text)}
          <Tooltip title={text}>
            <span>{text}</span>
          </Tooltip>
          {getFileTag(text)}
        </Space>
      ),
    },
    {
      title: '文件大小',
      dataIndex: 'Size',
      key: 'Size',
      width: 120,
      sorter: (a, b) => {
        const parseSize = (s: string) => {
          const num = parseFloat(s);
          if (s.includes('GB')) return num * 1024 * 1024 * 1024;
          if (s.includes('MB')) return num * 1024 * 1024;
          if (s.includes('KB')) return num * 1024;
          return num;
        };
        return parseSize(a.Size) - parseSize(b.Size);
      },
    },
    {
      title: '最后修改时间',
      dataIndex: 'LastModified',
      key: 'LastModified',
      width: 180,
      defaultSortOrder: 'descend',
      sorter: (a, b) => new Date(a.LastModified).getTime() - new Date(b.LastModified).getTime(),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right' as const,
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<DownloadOutlined />}
            onClick={() => handleDownload(record.FileName)}
          >
            下载
          </Button>
          <Button
            type="link"
            icon={<ShareAltOutlined />}
            onClick={() => handleShare(record)}
          >
            分享
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Input.Search
          placeholder="搜索文件名"
          allowClear
          style={{ width: 300 }}
          onSearch={setSearchText}
          onChange={(e) => !e.target.value && setSearchText('')}
        />
        <Button icon={<ReloadOutlined />} onClick={fetchData}>
          刷新
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={filteredData}
        rowKey="FileName"
        loading={loading}
        scroll={{ x: 800 }}
        pagination={{
          defaultPageSize: 20,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          showTotal: (total) => `共 ${total} 个文件`,
        }}
        size="middle"
      />

      {/* 创建分享弹窗 */}
      <Modal
        title={`分享文件: ${currentFile?.FileName || ''}`}
        open={shareModalVisible}
        onOk={handleCreateShare}
        onCancel={() => setShareModalVisible(false)}
        confirmLoading={shareLoading}
        okText="创建分享"
        cancelText="取消"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item label="有效期" name="expireHours">
            <Select>
              <Select.Option value={1}>1 小时</Select.Option>
              <Select.Option value={6}>6 小时</Select.Option>
              <Select.Option value={24}>1 天</Select.Option>
              <Select.Option value={72}>3 天</Select.Option>
              <Select.Option value={168}>7 天</Select.Option>
              <Select.Option value={720}>30 天</Select.Option>
              <Select.Option value={0}>永不过期</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item label="提取码">
            <Switch
              checked={useExtractCode}
              onChange={setUseExtractCode}
              checkedChildren="启用"
              unCheckedChildren="关闭"
            />
            {useExtractCode && (
              <div style={{ marginTop: 8, color: 'rgba(0,0,0,0.45)', fontSize: 12 }}>
                创建后将自动生成 4 位提取码
              </div>
            )}
          </Form.Item>
          <Form.Item label="最大下载次数" name="maxDownloads" extra="设为 0 表示不限制">
            <InputNumber min={0} max={99999} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>

      {/* 分享结果弹窗 */}
      <Modal
        title="分享创建成功"
        open={shareResultVisible}
        onCancel={() => setShareResultVisible(false)}
        footer={[
          <Button key="close" onClick={() => setShareResultVisible(false)}>关闭</Button>,
          <Button
            key="copy"
            type="primary"
            icon={<CopyOutlined />}
            onClick={() => {
              const url = `${window.location.origin}${shareResult?.shareUrl || ''}`;
              const text = shareResult?.extractCode
                ? `链接: ${url}\n提取码: ${shareResult.extractCode}`
                : `链接: ${url}`;
              copyToClipboard(text);
            }}
          >
            复制分享信息
          </Button>,
        ]}
      >
        {shareResult && (
          <div style={{ lineHeight: 2.2 }}>
            <Paragraph>
              <strong>分享链接：</strong>
              <a href={shareResult.shareUrl} target="_blank" rel="noopener noreferrer">
                {window.location.origin}{shareResult.shareUrl}
              </a>
            </Paragraph>
            {shareResult.extractCode && (
              <Paragraph>
                <strong>提取码：</strong>
                <Tag color="blue" style={{ fontSize: 16, padding: '2px 12px' }}>{shareResult.extractCode}</Tag>
              </Paragraph>
            )}
            {shareResult.expireAt && (
              <Paragraph>
                <strong>过期时间：</strong>{new Date(shareResult.expireAt).toLocaleString()}
              </Paragraph>
            )}
          </div>
        )}
      </Modal>
    </PageContainer>
  );
};

export default FileList;
