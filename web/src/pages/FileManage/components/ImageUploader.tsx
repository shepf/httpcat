import { PlusOutlined, UploadOutlined } from '@ant-design/icons';
import { Modal, Upload, Button, message, Input } from 'antd';
import type { RcFile, UploadProps, UploadFile, ItemRender } from 'antd/lib/upload/interface';
import axios from 'axios';
import React, { useState } from 'react';
import { URL } from 'url';


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
  const [showUploadList, setShowUploadList] = useState(true);

  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState('');
  const [previewTitle, setPreviewTitle] = useState('');
  const [fileList, setFileList] = useState<CustomUploadFile[]>([
    // 添加更多的图片项...
  ]);

  const handleCancel = () => setPreviewOpen(false);

  const handlePreview = async (file: UploadFile<CustomUploadFile>) => {
    if (!file.url && !file.preview) {
      file.preview = await getBase64(file.originFileObj as RcFile);
    }

    setPreviewImage(file.url || (file.preview as string));
    setPreviewOpen(true);
    setPreviewTitle(file.name || file.url!.substring(file.url!.lastIndexOf('/') + 1));
  };

  const handleChange: UploadProps<CustomUploadFile>['onChange'] = ({ fileList: newFileList }) =>
    setFileList(newFileList as CustomUploadFile[]);

  const handleEditDescription = (fileUid: string, newDescription: string) => {
    const updatedFileList = fileList.map((file) => {
      if (file.uid === fileUid) {
        return {
          ...file,
          name: newDescription,
        };
      }
      return file;
    });
    setFileList(updatedFileList);
  };
  
  const useCustomItemRender = (originNode: React.ReactNode, file: CustomUploadFile) => {
    const [editing, setEditing] = useState(false);
    const [newName, setNewName] = useState(file.name);
  
    const handleEdit = () => {
      setEditing(true);
    };
  
    const handleSave = () => {
      setEditing(false);
      handleEditDescription(file.uid, newName);
    };
  
    const render = editing ? (
      <div>
        <Input
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          onPressEnter={handleSave}
          onBlur={handleSave}
          autoFocus
        />
      </div>
    ) : (
      <div onClick={handleEdit}>{file.name}</div>
    );
  
    return render;
  };
  
  const CustomItemRender: React.FC<{ originNode: React.ReactNode, file: CustomUploadFile }> = ({ originNode, file }) => {
    const render = useCustomItemRender(originNode, file);
    return render;
  };

  const customItemRender: ItemRender<CustomUploadFile> = (originNode, file) => {
    const render = <CustomItemRender originNode={originNode} file={file} />;
    return render;
  };

  const beforeUpload = (file: RcFile) => {
    const acceptedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml']; // 允许的图片类型
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
    console.log("image upload, file:", file)
    console.log("image upload, file.uid:", file.uid)
    console.log("image upload, file.url:", file.url)

    try {
      let token = localStorage.getItem('token');
      if (null === token) {
          token = '';
      }

      const response = await axios.post('/api/v1/imageManage/upload', formData, {
        headers: {
          Authorization: `Bearer ${token}`,  // 替换为您的 Bearer Token
        },
      });

      console.log("image upload, response:", response)
      console.log("image upload, response url:", response.data.data.url)
            



      const baseUrl = window.location.origin; // 获取当前网站的域名和端口
      const thumbUrl = `${baseUrl}${response.data.data.thumbUrl}`; 
      console.log("image upload, thumbUrl:", thumbUrl)
      const imageUrl = `${baseUrl}${response.data.data.url}`; 
      console.log("image upload, imageUrl:", imageUrl)



      const uploadedFile = {
        uid: file.uid,
        name: response.data.data.name,
        status: "done",
        //url: imageUrl, // 使用缓存的图片 URL 通过使用前端缓存的图片 URL，您可以确保在前端展示图片时，不需要再通过网络请求获取图片资源。
        // thumbUrl: thumbUrl, // 根据后端返回的数据字段进行调整
        description: "", // 添加 description 字段并设置默认值        
      };


      //  onSuccess(response.data); 是有效果的，不需要自己在handleUpload中调用setFileList 注释掉
      // setFileList((prevList) => [...prevList, uploadedFile]); // 添加上传成功的图片到 fileList

      //通常，onSuccess 回调函数由 Upload 组件内部调用，它会传递上传成功后的响应数据给回调函数，并自动更新 fileList 状态，以反映上传成功的文件。
      onSuccess(uploadedFile, file);


      // 上传成功后调用回调函数刷新列表
      onUploadSuccess();

      // 隐藏上传列表
      setShowUploadList(false);

      message.info("upload success")

    } catch (error) {
      console.log("上传图片报错：")
      console.error(error)
      onError(error);
      // 隐藏上传列表
      setShowUploadList(false);
      message.error("上传图失败")

    }
  };


  const uploadButton = (
    <div>
      <Button icon={<UploadOutlined />}>Upload</Button>
    </div>
  );

  return (
    <>
      <Upload<CustomUploadFile>
        //action="" // 可以注释掉 使用自定义的上传请求 customRequest
        //name="file"	  //发到后台的文件参数名
        listType="text"  // 先用text吧，现在picture-card发现没有缩略图   picture-card：以卡片形式展示上传文件，显示缩略图、文件名称和操作按钮。
        fileList={fileList}  
        onPreview={handlePreview}
        onChange={handleChange}
        beforeUpload={beforeUpload}
        itemRender={customItemRender}
        customRequest={handleUpload}  // 使用自定义的上传请求
        showUploadList={showUploadList} // 根据状态变量控制上传列表的显示与隐藏
      >
        {fileList.length >= 10 ? null : uploadButton}
      </Upload>
      <Modal open={previewOpen} title={previewTitle} footer={null} onCancel={handleCancel}>
        <img alt="example" style={{ width: '100%' }} src={previewImage} />
      </Modal>
    </>
  );
};

export default CustomImageUpload;
