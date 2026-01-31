import { InboxOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { Alert, Card, message, Typography, UploadProps } from 'antd';
import Dragger from 'antd/lib/upload/Dragger';
import axios, { AxiosProgressEvent, AxiosResponse } from 'axios';
import React from 'react';
import { FormattedMessage, request, useIntl } from 'umi';
import styles from './Welcome.less';

const CodePreview: React.FC = ({ children }) => (
  <pre className={styles.pre}>
    <code>
      <Typography.Text copyable>{children}</Typography.Text>
    </code>
  </pre>
);

const Welcome: React.FC = () => {
  const intl = useIntl();


  const host = window.location.host;
  const apiUrl = `http://${host}/api/v1/file/`;

  type UploadRequestError = Error;
  const props: UploadProps = {
    name: 'file',
    multiple: true,
     showUploadList: true, // 显示上传列表
    //覆盖默认的上传行为，自定义上传实现
    customRequest: async ({ file, onSuccess, onError, onProgress }) => {
      const formData = new FormData();
      formData.append('f1', file);
  
      const headers = {
        UploadToken: 'httpcat:dZE8NVvimYNbV-YpJ9EFMKg3YaM=:eyJkZWFkbGluZSI6MH0=',
      };
  
      try {
        const response: AxiosResponse = await axios.post(`${apiUrl}upload`, formData, {
          headers: headers,
          onUploadProgress: (progressEvent: AxiosProgressEvent) => {
            const percent = Math.round((progressEvent.loaded * 100) / (progressEvent.total ?? 1));
            onProgress?.({ percent });
          },
        });
  
        const data = response.data;
        onSuccess?.(data);
      } catch (error) {
        onError?.(error as UploadRequestError);
        console.log("customRequest error:", error);
        message.error("网络连接失败，请检查网络状态。");
        return;
      }
    },
    // 在antd的Upload组件中，action属性用于指定上传文件的URL，但无法直接设置请求头。如果您需要设置请求头，可以使用beforeUpload属性来自定义上传行为
    // action: `${apiUrl}upload`,
    onChange(info) { //上传中、完成、失败都会调用这个函数
      const { status } = info.file;
      if (status !== 'uploading') {
        console.log(info.file, info.fileList);
      }
      if (status === 'done') {
        message.success(`${info.file.name} file uploaded successfully.`);
      } else if (status === 'error') {
        message.error(`${info.file.name} file upload failed.`);
      }
    },
    onDrop(e) {
      console.log('Dropped files', e.dataTransfer.files);
    },


  };

  return (
    <PageContainer>
      <Card>
        <Alert
          message={intl.formatMessage({
            id: 'pages.welcome.alertMessage',
            defaultMessage: 'Faster and stronger heavy-duty components have been released.',
          })}
          type="success"
          showIcon
          banner
          style={{
            margin: -12,
            marginBottom: 24,
          }}
        />
        <Typography.Text strong>
          <a
            // href="https://procomponents.ant.design/components/table"
            href="/welcome"
            rel="noopener noreferrer"
            target="__blank"
          >
            <FormattedMessage id="pages.welcome.link" defaultMessage="Welcome" />
          </a>
        </Typography.Text>
        <p>上传命令demo：</p>
        <CodePreview>
        {`curl -v -F "f1=@/root/test.md" -H "UploadToken: httpcat:dZE8NVvimYNbV-YpJ9EFMKg3YaM=:eyJkZWFkbGluZSI6MH0=" ${apiUrl}upload`}
        </CodePreview>
        <p>下载命令demo：</p>
        <CodePreview>
        {`wget ${apiUrl}download?filename=test.md`}
        </CodePreview>

      </Card>

      <Dragger {...props}>
        <p className="ant-upload-drag-icon">
          <InboxOutlined />
        </p>
        <p className="ant-upload-text">Click or drag file to this area to upload</p>
        <p className="ant-upload-hint">
          Support for a single or bulk upload. Strictly prohibited from uploading company data or other
          banned files.
        </p>
      </Dragger>
    </PageContainer>
  );
};

export default Welcome;
