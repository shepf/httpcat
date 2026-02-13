import { UploadOutlined } from '@ant-design/icons';
import { Modal, Upload, Button, message, Input } from 'antd';
import type { RcFile, UploadProps, UploadFile } from 'antd/lib/upload/interface';
import React, { useState } from 'react';
import { uploadImage } from '@/services/ant-design-pro/api';

interface CustomUploadFile extends UploadFile {
  description?: string;
}

const getBase64 = (file: RcFile): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = (error) => reject(error);
  });

interface CustomUploadProps {
  onUploadSuccess: () => void;
}

const CustomImageUpload: React.FC<CustomUploadProps> = ({ onUploadSuccess }) => {
  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState('');
  const [previewTitle, setPreviewTitle] = useState('');
  const [fileList, setFileList] = useState<CustomUploadFile[]>([]);

  const handlePreview = async (file: UploadFile) => {
    if (!file.url && !file.preview) {
      file.preview = await getBase64(file.originFileObj as RcFile);
    }
    setPreviewImage(file.url || (file.preview as string));
    setPreviewOpen(true);
    setPreviewTitle(file.name || file.url!.substring(file.url!.lastIndexOf('/') + 1));
  };

  const handleChange: UploadProps<CustomUploadFile>['onChange'] = ({ fileList: newFileList }) =>
    setFileList(newFileList as CustomUploadFile[]);

  const beforeUpload = (file: RcFile) => {
    const acceptedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml'];
    const isAcceptedType = acceptedTypes.includes(file.type);
    const isLt100M = file.size / 1024 / 1024 < 100;

    if (!isAcceptedType) {
      message.error('只支持上传 JPG/PNG/GIF/WebP/SVG 格式的图片!');
    }
    if (!isLt100M) {
      message.error('图片大小不能超过 100MB!');
    }

    return isAcceptedType && isLt100M;
  };

  const handleUpload = async (options: any) => {
    const { file, onSuccess, onError } = options;

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await uploadImage(formData);

      const uploadedFile = {
        uid: file.uid,
        name: response.data.name,
        status: 'done' as const,
      };

      onSuccess(uploadedFile, file);
      onUploadSuccess();
      setFileList([]);
      message.success('上传成功');
    } catch (error) {
      onError(error);
      setFileList([]);
      message.error('上传失败');
    }
  };

  return (
    <>
      <Upload<CustomUploadFile>
        listType="text"
        fileList={fileList}
        onPreview={handlePreview}
        onChange={handleChange}
        beforeUpload={beforeUpload}
        customRequest={handleUpload}
        showUploadList={fileList.length > 0}
      >
        {fileList.length >= 10 ? null : (
          <Button icon={<UploadOutlined />}>上传图片</Button>
        )}
      </Upload>
      <Modal open={previewOpen} title={previewTitle} footer={null} onCancel={() => setPreviewOpen(false)}>
        <img alt="preview" style={{ width: '100%' }} src={previewImage} />
      </Modal>
    </>
  );
};

export default CustomImageUpload;
