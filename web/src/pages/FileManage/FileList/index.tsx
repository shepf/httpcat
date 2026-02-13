import { DownloadOutlined, FileOutlined, FolderOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Button, Input, Space, Table, Tag, Tooltip, message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useEffect, useState } from 'react';
import { listFiles } from '@/services/ant-design-pro/api';

const FileList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<API.FileItem[]>([]);
  const [searchText, setSearchText] = useState('');

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
      width: 100,
      render: (_, record) => (
        <Button
          type="link"
          icon={<DownloadOutlined />}
          onClick={() => handleDownload(record.FileName)}
        >
          下载
        </Button>
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
        pagination={{
          defaultPageSize: 20,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          showTotal: (total) => `共 ${total} 个文件`,
        }}
        size="middle"
      />
    </PageContainer>
  );
};

export default FileList;
