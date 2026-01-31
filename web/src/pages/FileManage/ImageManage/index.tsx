import { DeleteOutlined, DownloadOutlined, EyeOutlined, FrownTwoTone, StopOutlined, WarningFilled } from '@ant-design/icons';
import { List, Card, Pagination, Button, Space, message, Modal } from 'antd';
import axios from 'axios';
import { useEffect, useState } from 'react';
import CustomImageUpload from '../components/ImageUploader';

const ImageList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const pageSize = 10; // 每页显示的图片数量
  const [totalRecords, setTotalRecords] = useState(0);


  const currentUrl = window.location.href;
  const currentIpAndPort = currentUrl.split('/')[2]; // 获取当前页面的 IP 和端口部分




  const fetchData = async () => {
    setLoading(true);

    try {
      let token = localStorage.getItem('token');
      if (null === token) {
          token = '';
      }

      const response = await axios.get('/api/v1/imageManage/listThumbImages', {
        params: {
          page: page,
          pageSize: pageSize,
        },
        headers: {
          Authorization: `Bearer ${token}`,  // 替换为您的 Bearer Token
        },
      });

      console.log("listThumbImages  response data: ", response.data)

      const responseData = response.data;
 

      // 获取图片数据
      const thumbnails = responseData.data;




		  //前端只通过 Base64 字符串无法确定图片的格式。Base64 只是一种表示图像数据的编码方式，并不能直接指示图像的格式。
		  //在前端展示 Base64 图片时，通常需要提供图像的 MIME 类型来指示图像的格式。MIME 类型是一种标识数据类型的字符串，例如 "image/jpeg" 表示 JPEG 图像，"image/png" 表示 PNG 图像。
		  // 将图像的 MIME 类型一并返回，这样前端就能够根据提供的 MIME 类型来正确解析和显示图像

      //根据文件名组装图像的 MIME 类型
      const updatedData = thumbnails.map((item: any) => {
        const fileExtension = item.FileName.split('.').pop().toLowerCase();
        let imageFormat = "image/jpeg"; // 默认格式为 JPEG
        if (fileExtension === "png") {
          imageFormat = "image/png";
        } else if (fileExtension === "gif") {
          imageFormat = "image/gif";
        }
        
        //拼接下载图片URL
        const downloadUrl = `http://${currentIpAndPort}/api/v1/imageManage/download?filename=${item.FileName}`;


        return {
          ...item,
          ImageFormat: imageFormat,
          ImageUrl: downloadUrl,
        };
      });


      setData(updatedData);
      setLoading(false);

    // 获取分页信息
      const pagination = responseData.pagination;
      console.log("获取分页信息 pagination: ",pagination)

      setPage(pagination.page);
      setTotalRecords(pagination.totalItems);

    } catch (error) {
      console.error('Error fetching thumbnails:', error);
      setLoading(false);
    }
  };


  useEffect(() => {
    fetchData();
  }, [page]);


  const handlePageChange = (pageNumber: number) => {
    setPage(pageNumber);
  };


  const handlePreview = (imageUrl: string) => {
    // 在这里处理预览图片的逻辑，例如弹出模态框显示大图
    console.log('Preview:', imageUrl);
  };

  const handleDownload = (imageUrl: string, FileName: string) => {
    let token = localStorage.getItem('token');
    if (null === token) {
        token = '';
    }

    
    axios({
      url: imageUrl,
      method: 'GET',
      responseType: 'blob', // 设置响应类型为 blob
      headers: {
        Authorization: `Bearer ${token}`,  // 替换为您的 Bearer Token
      },
    })
      .then((response) => {
        const url = URL.createObjectURL(new Blob([response.data]));
  
        const link = document.createElement('a');
        link.href = url;
        link.download = FileName; // 使用传入的文件名参数设置下载的文件名
        link.click();
  
        URL.revokeObjectURL(url);
      })
      .catch((error) => {
        console.error('Error downloading file:', error);
      });
  };

  const handleDelete = (FileName: string) => {
    // 获取认证 token
    let token = localStorage.getItem('token');
    if (null === token) {
      token = '';
    }
  
    // 发送请求给后端删除图片
    axios.delete(`/api/v1/imageManage/delete?filename=${FileName}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => {
        console.log('图片删除成功');
        // 在这里可以执行其他操作或刷新页面等
        message.success('图片删除成功');
        // 更新数据，将已删除的图片从列表中移除
        setData(prevData => prevData.filter(item => item.FileName !== FileName));
   
      })
      .catch((error) => {
        console.error('图片删除失败:', error);
        // 在这里可以给出错误提示或执行其他错误处理逻辑
        message.error('图片删除失败');
          // 更新数据，将已删除的图片从列表中移除
          setData(prevData => prevData.filter(item => item.FileName !== FileName));
      });
  };

  const handleClearAll = () => {
    // 在这里处理删除图片的逻辑，例如发送请求给后端删除图片
    console.log('handleClearAll:');
    Modal.confirm({
      title: '确认清空照片',
      content: '您确定要清空所有照片吗？',
      onOk() {
          // 在这里执行清空照片的操作
          console.log('清空照片');
          // 获取认证 token
          let token = localStorage.getItem('token');
          if (null === token) {
            token = '';
          }

          // 构建清空照片的请求 URL
          const clearUrl = '/api/v1/imageManage/clear';

          // 发送请求给后端清空照片
          axios.delete(clearUrl, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          })
          .then((response) => {
            console.log('照片清空成功');
            // 在这里可以执行其他操作或刷新页面等
            message.success('照片清空成功');
            // 清空数据，将列表中的所有图片移除
            setData([]);
          })
          .catch((error) => {
            console.error('照片清空失败:', error);
            // 在这里可以给出错误提示或执行其他错误处理逻辑
            message.error('照片清空失败');
          });


      },
    });


  };

  const handleRefresh = () => {
    fetchData(); // 执行刷新数据的逻辑
  };


  return (
    <div>

      <List
        header={
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '16px' }}>
            <Button type="primary" danger onClick={handleClearAll}>
              清空照片
            </Button>
            <div style={{ display: 'flex', gap: '8px' }}>
              <CustomImageUpload onUploadSuccess={handleRefresh} />
              <Button onClick={handleRefresh}>刷新</Button>
            </div>
          </div>
          
          
        }
        grid={{ gutter: 16, column: 5 }}
        dataSource={data}
        loading={loading}
        renderItem={(item: any) => (
          <List.Item>
            <Card 
              className="custom-card" 
              bodyStyle={{ display: 'none' }}  // 这个就是控制card内容区域样式的，需要内容展示的话，注意注释掉
              title={item.FileName}
              hoverable   //鼠标移过时可浮起
              cover={
                item.ThumbnailBase64 ? (
                  <div className="image-wrapper">
                      <img alt="Example Image"  className="custom-image" src={`data:${item.ImageFormat};base64,${item.ThumbnailBase64}`} />
                  </div>
                ) : (
                  <div className="no-thumbnail">
                    <div  className="custom-image">
                      <StopOutlined className="stop-icon" />
                    </div>
                  </div>
                )
              }
              actions={[
                <DownloadOutlined key="download" className="icon" onClick={() => handleDownload(item.ImageUrl, item.FileName)} />,
                <DeleteOutlined key="delete" className="icon" onClick={() => handleDelete(item.FileName)} />
              ]}
            >
              {/* {<Card.Meta title="" description={item.FileName} /> } */}
            </Card>
          </List.Item>
        )}
      />
      <div style={{ marginTop: '16px' }}> {/* 调整分页组件的上边距 */}
        <Pagination current={page} pageSize={pageSize} total={totalRecords} onChange={handlePageChange} />
      </div>
      <style>
        {`
          .custom-card {
            display: flex;
            flex-direction: column;
            height: 100%;
          }
          // .custom-card {
          //   height: 100%; /* 设置卡片的高度为图片的高度 */
          // }

          .image-wrapper {
            height: 100%; /* 设置图片容器的高度为卡片的高度 */
            display: flex;
            align-items: center;
            justify-content: center;
          }
          .custom-image {
            width: 100%;
            height: 100%;
            object-fit: cover; /* 图片按比例缩放，填充整个容器 */
          }

          .no-thumbnail {
            height: 100%;
            display: flex;
            justify-content: center;
            align-items: center;
            width: 100%; /* 或设置固定宽度，例如 width: 200px; */
          }

          .stop-icon {
            font-size: 8vw; /* 使用相对单位 vw 设置图标的大小 */
            transform: scale(1); /* 默认大小 */
            transition: transform 0.3s ease-in-out; /* 添加过渡效果 */
          }
          .no-thumbnail:hover .stop-icon {
            transform: scale(1.2); /* 鼠标悬停时放大图标 */
          }

        `}
      </style>
    </div>
  );
};

export default ImageList;
