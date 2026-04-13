import {
  CopyOutlined,
  DeleteOutlined,
  DownloadOutlined,
  EditOutlined,
  FileOutlined,
  FolderAddOutlined,
  FolderOutlined,
  HomeOutlined,
  ReloadOutlined,
  ShareAltOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import {
  Breadcrumb,
  Button,
  Form,
  Input,
  InputNumber,
  message,
  Modal,
  Popconfirm,
  Select,
  Space,
  Switch,
  Table,
  Tag,
  Tooltip,
  Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useEffect, useState } from 'react';
import { createShare, deleteFiles, createFolder, renameFile, listFiles } from '@/services/ant-design-pro/api';

const { Paragraph } = Typography;

const FileList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<API.FileItem[]>([]);
  const [searchText, setSearchText] = useState('');
  const [currentDir, setCurrentDir] = useState('');
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // 分享弹窗
  const [shareModalVisible, setShareModalVisible] = useState(false);
  const [shareResultVisible, setShareResultVisible] = useState(false);
  const [shareResult, setShareResult] = useState<API.CreateShareResult | null>(null);
  const [currentFile, setCurrentFile] = useState<API.FileItem | null>(null);
  const [shareLoading, setShareLoading] = useState(false);
  const [useExtractCode, setUseExtractCode] = useState(true);
  const [form] = Form.useForm();

  // 新建文件夹弹窗
  const [mkdirModalVisible, setMkdirModalVisible] = useState(false);
  const [mkdirLoading, setMkdirLoading] = useState(false);
  const [mkdirForm] = Form.useForm();

  // 重命名弹窗
  const [renameModalVisible, setRenameModalVisible] = useState(false);
  const [renameLoading, setRenameLoading] = useState(false);
  const [renameForm] = Form.useForm();
  const [renameTarget, setRenameTarget] = useState<API.FileItem | null>(null);

  const fetchData = async (dir?: string) => {
    setLoading(true);
    try {
      const res = await listFiles({ dir: dir ?? (currentDir || undefined) });
      setData(res.data || []);
      setSelectedRowKeys([]);
    } catch {
      message.error('获取文件列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData(currentDir);
  }, [currentDir]);

  // ========== 目录导航 ==========
  const pathSegments = currentDir ? currentDir.split('/').filter(Boolean) : [];

  const navigateTo = (dir: string) => {
    setSearchText('');
    setCurrentDir(dir);
  };

  const handleFolderClick = (folderName: string) => {
    const newDir = currentDir ? `${currentDir}/${folderName}` : folderName;
    navigateTo(newDir);
  };

  // ========== 文件操作 ==========
  const handleDownload = (fileName: string) => {
    const filePath = currentDir ? `${currentDir}/${fileName}` : fileName;
    window.open(`/api/v1/file/download?filename=${encodeURIComponent(filePath)}`);
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
      const filePath = currentDir ? `${currentDir}/${currentFile.FileName}` : currentFile.FileName;
      const params: API.CreateShareParams = {
        filePath: filePath,
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

  const handleRename = (record: API.FileItem) => {
    setRenameTarget(record);
    renameForm.resetFields();
    renameForm.setFieldsValue({ newName: record.FileName });
    setRenameModalVisible(true);
  };

  const handleRenameSubmit = async () => {
    if (!renameTarget) return;
    setRenameLoading(true);
    try {
      const values = await renameForm.validateFields();
      const res = await renameFile({
        oldName: renameTarget.FileName,
        newName: values.newName,
        dir: currentDir || undefined,
      });
      if (res.errorCode === 0) {
        message.success('重命名成功');
        setRenameModalVisible(false);
        fetchData();
      } else {
        message.error(res.msg || '重命名失败');
      }
    } catch {
      message.error('重命名失败');
    } finally {
      setRenameLoading(false);
    }
  };

  const handleDeleteSingle = async (fileName: string) => {
    try {
      const res = await deleteFiles({ files: [fileName], dir: currentDir || undefined });
      if (res.errorCode === 0) {
        const result = res.data;
        if (result?.deleted && result.deleted.length > 0) {
          message.success('删除成功');
        } else if (result?.failed && result.failed.length > 0) {
          message.error(`删除失败: ${result.failed[0].error}`);
        }
        fetchData();
      } else {
        message.error(res.msg || '删除失败');
      }
    } catch {
      message.error('删除失败');
    }
  };

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) return;
    try {
      const res = await deleteFiles({
        files: selectedRowKeys as string[],
        dir: currentDir || undefined,
      });
      if (res.errorCode === 0) {
        const result = res.data;
        const deletedCount = result?.deleted?.length || 0;
        const failedCount = result?.failed?.length || 0;
        if (failedCount > 0) {
          message.warning(`成功删除 ${deletedCount} 项，失败 ${failedCount} 项`);
        } else {
          message.success(`成功删除 ${deletedCount} 项`);
        }
        fetchData();
      } else {
        message.error(res.msg || '批量删除失败');
      }
    } catch {
      message.error('批量删除失败');
    }
  };

  const handleCreateFolder = async () => {
    setMkdirLoading(true);
    try {
      const values = await mkdirForm.validateFields();
      const res = await createFolder({ name: values.folderName, dir: currentDir || undefined });
      if (res.errorCode === 0) {
        message.success('文件夹创建成功');
        setMkdirModalVisible(false);
        fetchData();
      } else {
        message.error(res.msg || '创建文件夹失败');
      }
    } catch {
      message.error('创建文件夹失败');
    } finally {
      setMkdirLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).then(
      () => message.success('已复制到剪贴板'),
      () => message.error('复制失败，请手动复制'),
    );
  };

  // ========== 图标和标签 ==========
  const getFileIcon = (record: API.FileItem) => {
    if (record.IsDir) return <FolderOutlined style={{ color: '#faad14', fontSize: 18 }} />;
    const ext = record.FileName.split('.').pop()?.toLowerCase();
    const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico'];
    if (imageExts.includes(ext || '')) return <FileOutlined style={{ color: '#52c41a' }} />;
    const archiveExts = ['zip', 'tar', 'gz', 'rar', '7z', 'bz2', 'xz'];
    if (archiveExts.includes(ext || '')) return <FileOutlined style={{ color: '#faad14' }} />;
    return <FileOutlined style={{ color: '#1890ff' }} />;
  };

  const getFileTag = (record: API.FileItem) => {
    if (record.IsDir) return <Tag color="gold">文件夹</Tag>;
    const ext = record.FileName.split('.').pop()?.toLowerCase() || '';
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

  // ========== 表格列 ==========
  const columns: ColumnsType<API.FileItem> = [
    {
      title: '文件名',
      dataIndex: 'FileName',
      key: 'FileName',
      ellipsis: true,
      sorter: (a, b) => {
        // 目录优先
        if (a.IsDir && !b.IsDir) return -1;
        if (!a.IsDir && b.IsDir) return 1;
        return a.FileName.localeCompare(b.FileName);
      },
      render: (_: string, record: API.FileItem) => (
        <Space>
          {getFileIcon(record)}
          {record.IsDir ? (
            <a onClick={() => handleFolderClick(record.FileName)} style={{ fontWeight: 500 }}>
              {record.FileName}
            </a>
          ) : (
            <Tooltip title={record.FileName}>
              <span>{record.FileName}</span>
            </Tooltip>
          )}
          {getFileTag(record)}
        </Space>
      ),
    },
    {
      title: '文件大小',
      dataIndex: 'Size',
      key: 'Size',
      width: 120,
      sorter: (a, b) => {
        if (a.IsDir && b.IsDir) return 0;
        if (a.IsDir) return -1;
        if (b.IsDir) return 1;
        const parseSize = (s: string) => {
          if (s === '-') return 0;
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
      width: 260,
      fixed: 'right' as const,
      render: (_, record) => (
        <Space size={0}>
          {record.IsDir ? (
            <Button type="link" size="small" onClick={() => handleFolderClick(record.FileName)}>
              打开
            </Button>
          ) : (
            <>
              <Button type="link" size="small" icon={<DownloadOutlined />} onClick={() => handleDownload(record.FileName)}>
                下载
              </Button>
              <Button type="link" size="small" icon={<ShareAltOutlined />} onClick={() => handleShare(record)}>
                分享
              </Button>
            </>
          )}
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleRename(record)}>
            重命名
          </Button>
          <Popconfirm
            title={`确定删除 "${record.FileName}"？`}
            description={record.IsDir ? '注意：只能删除空文件夹' : undefined}
            onConfirm={() => handleDeleteSingle(record.FileName)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const rowSelection = {
    selectedRowKeys,
    onChange: (newSelectedRowKeys: React.Key[]) => {
      setSelectedRowKeys(newSelectedRowKeys);
    },
  };

  return (
    <PageContainer>
      {/* 面包屑导航 */}
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item>
          <a onClick={() => navigateTo('')}>
            <HomeOutlined /> 根目录
          </a>
        </Breadcrumb.Item>
        {pathSegments.map((seg, idx) => {
          const path = pathSegments.slice(0, idx + 1).join('/');
          const isLast = idx === pathSegments.length - 1;
          return (
            <Breadcrumb.Item key={path}>
              {isLast ? (
                <span><FolderOutlined /> {seg}</span>
              ) : (
                <a onClick={() => navigateTo(path)}>
                  <FolderOutlined /> {seg}
                </a>
              )}
            </Breadcrumb.Item>
          );
        })}
      </Breadcrumb>

      {/* 工具栏 */}
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', flexWrap: 'wrap', gap: 8 }}>
        <Space wrap>
          <Input.Search
            placeholder="搜索文件名"
            allowClear
            style={{ width: 300 }}
            onSearch={setSearchText}
            onChange={(e) => !e.target.value && setSearchText('')}
          />
        </Space>
        <Space wrap>
          {selectedRowKeys.length > 0 && (
            <Popconfirm
              title={`确定删除选中的 ${selectedRowKeys.length} 项？`}
              onConfirm={handleBatchDelete}
              okText="确定"
              cancelText="取消"
            >
              <Button danger icon={<DeleteOutlined />}>
                批量删除 ({selectedRowKeys.length})
              </Button>
            </Popconfirm>
          )}
          <Button
            icon={<FolderAddOutlined />}
            onClick={() => {
              mkdirForm.resetFields();
              setMkdirModalVisible(true);
            }}
          >
            新建文件夹
          </Button>
          <Button icon={<ReloadOutlined />} onClick={() => fetchData()}>
            刷新
          </Button>
        </Space>
      </div>

      {/* 文件表格 */}
      <Table
        columns={columns}
        dataSource={filteredData}
        rowKey="FileName"
        loading={loading}
        scroll={{ x: 800 }}
        rowSelection={rowSelection}
        pagination={{
          defaultPageSize: 20,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          showTotal: (total) => `共 ${total} 项${selectedRowKeys.length > 0 ? `，已选 ${selectedRowKeys.length} 项` : ''}`,
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

      {/* 新建文件夹弹窗 */}
      <Modal
        title="新建文件夹"
        open={mkdirModalVisible}
        onOk={handleCreateFolder}
        onCancel={() => setMkdirModalVisible(false)}
        confirmLoading={mkdirLoading}
        okText="创建"
        cancelText="取消"
      >
        <Form form={mkdirForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            label="文件夹名称"
            name="folderName"
            rules={[
              { required: true, message: '请输入文件夹名称' },
              { pattern: /^[^/\\:*?"<>|]+$/, message: '名称不能包含特殊字符' },
            ]}
          >
            <Input placeholder="请输入文件夹名称" autoFocus />
          </Form.Item>
        </Form>
      </Modal>

      {/* 重命名弹窗 */}
      <Modal
        title={`重命名: ${renameTarget?.FileName || ''}`}
        open={renameModalVisible}
        onOk={handleRenameSubmit}
        onCancel={() => setRenameModalVisible(false)}
        confirmLoading={renameLoading}
        okText="确定"
        cancelText="取消"
      >
        <Form form={renameForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            label="新名称"
            name="newName"
            rules={[
              { required: true, message: '请输入新名称' },
              { pattern: /^[^/\\:*?"<>|]+$/, message: '名称不能包含特殊字符' },
            ]}
          >
            <Input placeholder="请输入新名称" autoFocus />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default FileList;
