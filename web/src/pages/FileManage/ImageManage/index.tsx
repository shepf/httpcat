import {
  ClearOutlined,
  CopyOutlined,
  DeleteOutlined,
  DownloadOutlined,
  EyeOutlined,
  LinkOutlined,
  ReloadOutlined,
  ShareAltOutlined,
  StopOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import { List, Card, Pagination, Button, Space, message, Modal, Select, Tooltip, Input, Form, InputNumber, Switch, Tag, Typography } from 'antd';
import { useEffect, useState, useRef, useCallback } from 'react';
import { listThumbImages, downloadImage, deleteImage, clearImages, createShare } from '@/services/ant-design-pro/api';
import CustomImageUpload from '../components/ImageUploader';
import styles from './index.less';

const AUTO_REFRESH_OPTIONS = [
  { label: '10秒', value: 10 },
  { label: '30秒', value: 30 },
  { label: '60秒', value: 60 },
];

const { Paragraph } = Typography;

const ImageList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const pageSize = 10;
  const [totalRecords, setTotalRecords] = useState(0);
  const [previewVisible, setPreviewVisible] = useState(false);
  const [previewImage, setPreviewImage] = useState('');
  const [searchText, setSearchText] = useState('');

  const [autoRefresh, setAutoRefresh] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState(30);
  const [countdown, setCountdown] = useState(0);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const countdownRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // 分享相关状态
  const [shareModalVisible, setShareModalVisible] = useState(false);
  const [shareResultVisible, setShareResultVisible] = useState(false);
  const [shareResult, setShareResult] = useState<API.CreateShareResult | null>(null);
  const [currentImage, setCurrentImage] = useState<string>('');
  const [shareLoading, setShareLoading] = useState(false);
  const [useExtractCode, setUseExtractCode] = useState(true);
  const [shareForm] = Form.useForm();

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const response = await listThumbImages({ page, pageSize, search: searchText || undefined });
      const thumbnails = response.data || [];

      const updatedData = thumbnails.map((item: API.ImageItem) => {
        const ext = item.FileName.split('.').pop()?.toLowerCase();
        const formatMap: Record<string, string> = {
          png: 'image/png',
          gif: 'image/gif',
          webp: 'image/webp',
          svg: 'image/svg+xml',
        };
        return {
          ...item,
          ImageFormat: formatMap[ext || ''] || 'image/jpeg',
        };
      });

      setData(updatedData);
      const pagination = response.pagination;
      if (pagination) {
        setPage(pagination.page);
        setTotalRecords(pagination.totalItems);
      }
    } catch (error) {
      message.error('获取图片列表失败');
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, searchText]);

  useEffect(() => {
    fetchData();
  }, [page, searchText]);

  // 自动刷新逻辑
  useEffect(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
    if (countdownRef.current) {
      clearInterval(countdownRef.current);
      countdownRef.current = null;
    }

    if (autoRefresh) {
      setCountdown(refreshInterval);

      countdownRef.current = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) return refreshInterval;
          return prev - 1;
        });
      }, 1000);

      timerRef.current = setInterval(() => {
        fetchData();
        setCountdown(refreshInterval);
      }, refreshInterval * 1000);
    } else {
      setCountdown(0);
    }

    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
      if (countdownRef.current) clearInterval(countdownRef.current);
    };
  }, [autoRefresh, refreshInterval, fetchData]);

  const handlePreview = (base64: string, format: string) => {
    setPreviewImage(`data:${format};base64,${base64}`);
    setPreviewVisible(true);
  };

  const handleDownload = async (fileName: string) => {
    try {
      const blob = await downloadImage(fileName);
      const url = URL.createObjectURL(new Blob([blob]));
      const link = document.createElement('a');
      link.href = url;
      link.download = fileName;
      link.click();
      URL.revokeObjectURL(url);
    } catch (error) {
      message.error('下载失败');
    }
  };

  const handleDelete = async (fileName: string) => {
    try {
      await deleteImage(fileName);
      message.success('图片删除成功');
      setData((prev) => prev.filter((item) => item.FileName !== fileName));
    } catch (error) {
      message.error('图片删除失败');
    }
  };

  const handleCopyImage = async (fileName: string) => {
    try {
      const blob = await downloadImage(fileName);
      const ext = fileName.split('.').pop()?.toLowerCase();
      // Clipboard API 只支持 image/png，其他格式需要转换
      let pngBlob: Blob;
      if (ext === 'png') {
        pngBlob = new Blob([blob], { type: 'image/png' });
      } else {
        // 通过 canvas 转换为 PNG
        const img = new window.Image();
        const url = URL.createObjectURL(new Blob([blob]));
        pngBlob = await new Promise<Blob>((resolve, reject) => {
          img.onload = () => {
            const canvas = document.createElement('canvas');
            canvas.width = img.naturalWidth;
            canvas.height = img.naturalHeight;
            const ctx = canvas.getContext('2d');
            ctx?.drawImage(img, 0, 0);
            canvas.toBlob((b) => {
              URL.revokeObjectURL(url);
              b ? resolve(b) : reject(new Error('转换失败'));
            }, 'image/png');
          };
          img.onerror = () => {
            URL.revokeObjectURL(url);
            reject(new Error('图片加载失败'));
          };
          img.src = url;
        });
      }
      await navigator.clipboard.write([
        new ClipboardItem({ 'image/png': pngBlob }),
      ]);
      message.success('图片已复制到剪贴板');
    } catch (error) {
      message.error('复制失败，请检查浏览器权限');
    }
  };

  const handleCopyLink = (fileName: string) => {
    const link = `${window.location.origin}/api/v1/imageManage/download?filename=${encodeURIComponent(fileName)}`;
    navigator.clipboard.writeText(link).then(
      () => message.success('图片链接已复制'),
      () => message.error('复制链接失败'),
    );
  };

  const handleShare = (fileName: string) => {
    setCurrentImage(fileName);
    setUseExtractCode(true);
    shareForm.resetFields();
    shareForm.setFieldsValue({ expireHours: 24, maxDownloads: 0 });
    setShareModalVisible(true);
  };

  const handleCreateShare = async () => {
    if (!currentImage) return;
    setShareLoading(true);
    try {
      const values = await shareForm.validateFields();
      const params: API.CreateShareParams = {
        filePath: currentImage,
        fileName: currentImage,
        fileType: 'image',
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

  const copyShareToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).then(
      () => message.success('已复制到剪贴板'),
      () => message.error('复制失败，请手动复制'),
    );
  };

  const handleSearch = (value: string) => {
    setSearchText(value);
    setPage(1); // 搜索时重置到第1页
  };

  const handleClearAll = () => {
    Modal.confirm({
      title: '确认清空照片',
      content: '您确定要清空所有照片吗？此操作不可撤销。',
      okText: '确认清空',
      okType: 'danger',
      cancelText: '取消',
      async onOk() {
        try {
          await clearImages();
          message.success('照片清空成功');
          setData([]);
          setTotalRecords(0);
        } catch (error) {
          message.error('照片清空失败');
        }
      },
    });
  };

  return (
    <div className={styles.container}>
      <List
        header={
          <div className={styles.header}>
            <Space>
              <Button type="primary" danger icon={<ClearOutlined />} onClick={handleClearAll}>
                清空照片
              </Button>
              <Input.Search
                placeholder="搜索图片名"
                allowClear
                style={{ width: 250 }}
                onSearch={handleSearch}
                onChange={(e) => !e.target.value && handleSearch('')}
              />
            </Space>
            <Space>
              <CustomImageUpload onUploadSuccess={fetchData} />
              <div className={styles.autoRefreshGroup}>
                <Tooltip title={autoRefresh ? `${countdown}秒后自动刷新` : '开启自动刷新'}>
                  <Button
                    type={autoRefresh ? 'primary' : 'default'}
                    icon={<SyncOutlined spin={autoRefresh} />}
                    onClick={() => setAutoRefresh(!autoRefresh)}
                  >
                    {autoRefresh ? `自动刷新 ${countdown}s` : '自动刷新'}
                  </Button>
                </Tooltip>
                <Select
                  value={refreshInterval}
                  onChange={(val) => setRefreshInterval(val)}
                  options={AUTO_REFRESH_OPTIONS}
                  size="middle"
                  className={styles.intervalSelect}
                />
              </div>
              <Button icon={<ReloadOutlined />} onClick={fetchData}>
                刷新
              </Button>
            </Space>
          </div>
        }
        grid={{ gutter: 16, column: 5 }}
        dataSource={data}
        loading={loading}
        locale={{ emptyText: '暂无图片' }}
        renderItem={(item: any) => (
          <List.Item>
            <Card
              className={styles.imageCard}
              bodyStyle={{ display: 'none' }}
              title={<span className={styles.cardTitle}>{item.FileName}</span>}
              hoverable
              cover={
                item.ThumbnailBase64 ? (
                  <div className={styles.imageWrapper}>
                    <img
                      alt={item.FileName}
                      className={styles.image}
                      src={`data:${item.ImageFormat};base64,${item.ThumbnailBase64}`}
                    />
                  </div>
                ) : (
                  <div className={styles.noThumbnail}>
                    <StopOutlined className={styles.stopIcon} />
                  </div>
                )
              }
              actions={[
                item.ThumbnailBase64 ? (
                  <EyeOutlined
                    key="preview"
                    onClick={() => handlePreview(item.ThumbnailBase64, item.ImageFormat)}
                  />
                ) : null,
                <CopyOutlined key="copy" onClick={() => handleCopyImage(item.FileName)} />,
                <LinkOutlined key="link" onClick={() => handleCopyLink(item.FileName)} />,
                <ShareAltOutlined key="share" onClick={() => handleShare(item.FileName)} />,
                <DownloadOutlined key="download" onClick={() => handleDownload(item.FileName)} />,
                <DeleteOutlined key="delete" onClick={() => handleDelete(item.FileName)} />,
              ].filter(Boolean)}
            />
          </List.Item>
        )}
      />

      {totalRecords > 0 && (
        <div className={styles.pagination}>
          <Pagination
            current={page}
            pageSize={pageSize}
            total={totalRecords}
            onChange={setPage}
            showTotal={(total) => `共 ${total} 张图片`}
          />
        </div>
      )}

      <Modal
        open={previewVisible}
        footer={null}
        onCancel={() => setPreviewVisible(false)}
        width={800}
        centered
      >
        <img alt="预览" style={{ width: '100%' }} src={previewImage} />
      </Modal>

      {/* 创建分享弹窗 */}
      <Modal
        title={`分享图片: ${currentImage}`}
        open={shareModalVisible}
        onOk={handleCreateShare}
        onCancel={() => setShareModalVisible(false)}
        confirmLoading={shareLoading}
        okText="创建分享"
        cancelText="取消"
      >
        <Form form={shareForm} layout="vertical" style={{ marginTop: 16 }}>
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
              copyShareToClipboard(text);
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
    </div>
  );
};

export default ImageList;
