import {
  FileOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
} from '@ant-design/icons';
import { Modal, Spin, Typography, message } from 'antd';
import React, { useEffect, useState } from 'react';
import { getPreviewInfo } from '@/services/ant-design-pro/api';

const { Paragraph, Text } = Typography;

interface FilePreviewProps {
  visible: boolean;
  fileName: string; // 完整的文件路径（含子目录，如 "subdir/file.txt"）
  onClose: () => void;
}

const FilePreview: React.FC<FilePreviewProps> = ({ visible, fileName, onClose }) => {
  const [loading, setLoading] = useState(true);
  const [previewInfo, setPreviewInfo] = useState<API.PreviewInfo | null>(null);
  const [textContent, setTextContent] = useState<string>('');
  const [fullscreen, setFullscreen] = useState(false);

  const token = localStorage.getItem('token') || '';
  const previewUrl = `/api/v1/file/preview?filename=${encodeURIComponent(fileName)}&token=${encodeURIComponent(token)}`;

  useEffect(() => {
    if (visible && fileName) {
      setLoading(true);
      setTextContent('');
      setPreviewInfo(null);

      getPreviewInfo({ filename: fileName })
        .then((res) => {
          if (res.errorCode === 0 && res.data) {
            setPreviewInfo(res.data);
            // 如果是文本类型，获取内容进行高亮显示
            if (res.data.previewType === 'text' || res.data.previewType === 'markdown') {
              fetch(previewUrl, {
                headers: { Authorization: `Bearer ${token}` },
              })
                .then((r) => r.text())
                .then((text) => setTextContent(text))
                .catch(() => message.error('加载文件内容失败'));
            }
          } else {
            message.error(res.msg || '获取预览信息失败');
          }
        })
        .catch(() => message.error('获取预览信息失败'))
        .finally(() => setLoading(false));
    }
  }, [visible, fileName]);

  const renderPreview = () => {
    if (!previewInfo || !previewInfo.canPreview) {
      return (
        <div style={{ textAlign: 'center', padding: '60px 20px' }}>
          <FileOutlined style={{ fontSize: 64, color: '#999' }} />
          <Paragraph style={{ marginTop: 16, color: '#999' }}>
            该文件类型不支持在线预览
          </Paragraph>
        </div>
      );
    }

    switch (previewInfo.previewType) {
      case 'text':
      case 'markdown':
        return (
          <div
            style={{
              maxHeight: fullscreen ? 'calc(100vh - 120px)' : '65vh',
              overflow: 'auto',
              background: '#1e1e1e',
              borderRadius: 8,
              padding: 16,
            }}
          >
            <pre
              style={{
                margin: 0,
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-all',
                color: '#d4d4d4',
                fontSize: 13,
                lineHeight: 1.6,
                fontFamily: "'Fira Code', 'Cascadia Code', 'JetBrains Mono', Consolas, monospace",
              }}
            >
              {textContent || '加载中...'}
            </pre>
          </div>
        );

      case 'pdf':
        return (
          <iframe
            src={previewUrl}
            style={{
              width: '100%',
              height: fullscreen ? 'calc(100vh - 120px)' : '70vh',
              border: 'none',
              borderRadius: 8,
            }}
            title="PDF Preview"
          />
        );

      case 'image':
        return (
          <div style={{ textAlign: 'center', padding: 16 }}>
            <img
              src={previewUrl}
              alt={previewInfo.fileName}
              style={{
                maxWidth: '100%',
                maxHeight: fullscreen ? 'calc(100vh - 120px)' : '65vh',
                objectFit: 'contain',
                borderRadius: 8,
                boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
              }}
            />
          </div>
        );

      case 'video':
        return (
          <div style={{ textAlign: 'center', padding: 16 }}>
            <video
              controls
              autoPlay={false}
              style={{
                maxWidth: '100%',
                maxHeight: fullscreen ? 'calc(100vh - 120px)' : '65vh',
                borderRadius: 8,
                boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
              }}
            >
              <source src={previewUrl} />
              您的浏览器不支持视频播放
            </video>
          </div>
        );

      case 'audio':
        return (
          <div style={{ textAlign: 'center', padding: '40px 20px' }}>
            <FileOutlined style={{ fontSize: 64, color: '#1890ff', marginBottom: 24 }} />
            <Paragraph>
              <Text strong>{previewInfo.fileName}</Text>
            </Paragraph>
            <Paragraph type="secondary">
              {previewInfo.sizeFormatted} · {previewInfo.extension?.toUpperCase()}
            </Paragraph>
            <audio
              controls
              autoPlay={false}
              style={{ width: '100%', maxWidth: 500, marginTop: 16 }}
            >
              <source src={previewUrl} />
              您的浏览器不支持音频播放
            </audio>
          </div>
        );

      default:
        return (
          <div style={{ textAlign: 'center', padding: '60px 20px' }}>
            <FileOutlined style={{ fontSize: 64, color: '#999' }} />
            <Paragraph style={{ marginTop: 16, color: '#999' }}>
              该文件类型不支持在线预览
            </Paragraph>
          </div>
        );
    }
  };

  const modalWidth = fullscreen ? '100vw' : 900;
  const modalStyle = fullscreen
    ? { top: 0, maxWidth: '100vw', paddingBottom: 0 }
    : {};

  return (
    <Modal
      title={
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', paddingRight: 32 }}>
          <span>
            预览: {previewInfo?.fileName || fileName.split('/').pop()}
            {previewInfo && (
              <Text type="secondary" style={{ fontSize: 12, marginLeft: 8 }}>
                {previewInfo.sizeFormatted} · {previewInfo.extension?.toUpperCase()}
              </Text>
            )}
          </span>
          <span
            style={{ cursor: 'pointer', fontSize: 16 }}
            onClick={() => setFullscreen(!fullscreen)}
            title={fullscreen ? '退出全屏' : '全屏'}
          >
            {fullscreen ? <FullscreenExitOutlined /> : <FullscreenOutlined />}
          </span>
        </div>
      }
      open={visible}
      onCancel={onClose}
      footer={null}
      width={modalWidth}
      style={modalStyle}
      bodyStyle={{ padding: fullscreen ? '8px' : '16px' }}
      destroyOnClose
      centered={!fullscreen}
    >
      <Spin spinning={loading}>{renderPreview()}</Spin>
    </Modal>
  );
};

export default FilePreview;
