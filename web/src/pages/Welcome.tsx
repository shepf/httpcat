import { CloudUploadOutlined, DownloadOutlined, InboxOutlined, RocketOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Alert, Card, message, Space, Spin, Tag, Typography, UploadProps } from 'antd';
import Dragger from 'antd/lib/upload/Dragger';
import React, { useEffect, useState } from 'react';
import { getFirstUploadToken, getVersion } from '@/services/ant-design-pro/api';
import request from 'umi-request';
import styles from './Welcome.less';

const { Text } = Typography;

const CodePreview: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <pre className={styles.pre}>
    <code>
      <Typography.Text copyable>{children}</Typography.Text>
    </code>
  </pre>
);

const Welcome: React.FC = () => {
  const [uploadToken, setUploadToken] = useState<string>('');
  const [tokenLoading, setTokenLoading] = useState<boolean>(true);
  const [currentVersion, setCurrentVersion] = useState<string>('');
  const [latestVersion, setLatestVersion] = useState<string>('');
  const [hasNewVersion, setHasNewVersion] = useState<boolean>(false);

  const host = window.location.host;
  const apiUrl = `http://${host}/api/v1/file/`;

  useEffect(() => {
    const fetchToken = async () => {
      try {
        const token = await getFirstUploadToken();
        setUploadToken(token);
      } catch (error) {
        console.error('获取上传Token失败:', error);
      } finally {
        setTokenLoading(false);
      }
    };
    fetchToken();

    // 获取当前版本
    const fetchVersion = async () => {
      try {
        const res = await getVersion();
        const ver = res.data?.version || '';
        setCurrentVersion(ver);
        // 检查 GitHub 最新版本
        try {
          const ghRes = await fetch('https://api.github.com/repos/shepf/httpcat/releases/latest');
          if (ghRes.ok) {
            const ghData = await ghRes.json();
            const latest = ghData.tag_name || '';
            setLatestVersion(latest);
            if (latest && ver && latest.replace(/^v/, '') !== ver.replace(/^v/, '')) {
              setHasNewVersion(true);
            }
          }
        } catch (_) { /* 网络不通忽略 */ }
      } catch (_) { /* 忽略 */ }
    };
    fetchVersion();
  }, []);

  const props: UploadProps = {
    name: 'file',
    multiple: true,
    showUploadList: true,
    customRequest: async ({ file, onSuccess, onError, onProgress }) => {
      if (!uploadToken) {
        message.error('暂无可用的上传Token，请先在管理页面创建');
        onError?.(new Error('No upload token'));
        return;
      }

      const formData = new FormData();
      formData.append('f1', file);

      try {
        const response = await request('/api/v1/file/upload', {
          method: 'POST',
          data: formData,
          headers: { UploadToken: uploadToken },
          requestType: 'form',
        });
        onSuccess?.(response);
      } catch (error) {
        onError?.(error as Error);
        message.error('上传失败，请检查网络连接');
      }
    },
    onChange(info) {
      const { status } = info.file;
      if (status === 'done') {
        message.success(`${info.file.name} 上传成功`);
      } else if (status === 'error') {
        message.error(`${info.file.name} 上传失败`);
      }
    },
  };

  return (
    <PageContainer>
      {hasNewVersion ? (
        <Alert
          type="info"
          banner
          showIcon
          icon={<RocketOutlined />}
          message={
            <span>
              当前版本 <Tag color="default">{currentVersion}</Tag>
              ，发现新版本 <Tag color="blue">{latestVersion}</Tag>
              <a
                href="https://github.com/shepf/httpcat/releases"
                target="_blank"
                rel="noopener noreferrer"
                style={{ marginLeft: 8 }}
              >
                前往更新
              </a>
            </span>
          }
          style={{ marginBottom: 16, borderRadius: 6 }}
        />
      ) : currentVersion ? (
        <Alert
          type="success"
          banner
          showIcon
          icon={<RocketOutlined />}
          message={
            <span>
              当前版本 <Tag color="green">{currentVersion}</Tag> 已是最新
            </span>
          }
          style={{ marginBottom: 16, borderRadius: 6 }}
        />
      ) : null}

      <Card title="快捷上传" className={styles.uploadCard} style={{ marginTop: 16 }}>
        {tokenLoading ? (
          <div style={{ textAlign: 'center', padding: '40px 0' }}>
            <Spin tip="加载中..." />
          </div>
        ) : !uploadToken ? (
          <Alert
            message="暂无可用的上传Token"
            description="请先前往「管理 > 上传Token管理」页面创建并生成 Token"
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
          />
        ) : null}
        <Dragger {...props} disabled={!uploadToken}>
          <p className="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
          <p className="ant-upload-hint">支持单个或批量文件上传</p>
        </Dragger>
      </Card>

      <Card title="命令行使用示例" className={styles.cmdCard} style={{ marginTop: 16 }}>
        <Space direction="vertical" style={{ width: '100%' }}>
          <div>
            <Text strong><CloudUploadOutlined /> 上传文件：</Text>
            <CodePreview>
              {`curl -F "f1=@/path/to/file" -H "UploadToken: ${uploadToken || '<your_token>'}" ${apiUrl}upload`}
            </CodePreview>
          </div>
          <div>
            <Text strong><DownloadOutlined /> 下载文件：</Text>
            <CodePreview>
              {`wget ${apiUrl}download?filename=example.txt`}
            </CodePreview>
          </div>
        </Space>
      </Card>
    </PageContainer>
  );
};

export default Welcome;
