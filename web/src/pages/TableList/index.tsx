import { getFileInfo, removeUploadHistoryLog, uploadHistoryLogs} from '@/services/ant-design-pro/api';
import { DownloadOutlined, PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, Input, message, Tooltip } from 'antd';
import React, { useEffect, useRef, useState } from 'react';
import { FormattedMessage, useIntl } from 'umi';
import UpdateForm from './components/UpdateForm';



/**
 *  Delete Upload History Log
 * @zh-CN 删除上传文件历史日志
 *
 * @param selectedRows
 */
const handleRemove = async (selectedRows: API.UploadHistoryLogItem[]) => {
  const hide = message.loading('正在删除');
  if (!selectedRows) return true;
  try {
    await removeUploadHistoryLog({
      id: selectedRows.map((row) => row.id),
    });
    // await removeUploadHistoryLog({ id: ['123']});

    hide();
    message.success('Deleted successfully and will refresh soon');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again');
    return false;
  }
};

const TableList: React.FC = () => {
  /**
   * @en-US Pop-up window of new window
   * @zh-CN 新建窗口的弹窗
   *  */
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  /**
   * @en-US The pop-up window of the distribution update window
   * @zh-CN 分布更新窗口的弹窗
   * */
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);

  const [showDetail, setShowDetail] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.UploadHistoryLogItem>();
  const [selectedRowsState, setSelectedRows] = useState<API.UploadHistoryLogItem[]>([]);

  // 下载button是否显示
  const [isDownloadButtonDisabled, setIsDownloadButtonDisabled] = useState(true);
  useEffect(() => {
    const fetchData = async () => {
      const fileInfo = await getFileInfo({ filename: currentRow?.filename, file_md5: currentRow?.file_md5 });
      setIsDownloadButtonDisabled(fileInfo?.errorCode === 22);
    };
  
    fetchData();
  }, [currentRow?.filename, currentRow?.file_md5]);
  

  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */
  const intl = useIntl();

  const columns: ProColumns<API.UploadHistoryLogItem>[] = [
    {
      title: (
        <FormattedMessage
          id="pages.searchTable.updateForm.fileName.nameLabel"
        />
      ),
      dataIndex: 'filename',
      render: (dom, entity) => {
        return (
          <a
            onClick={() => {
              setCurrentRow(entity);
              setShowDetail(true);
            }}
          >
            {dom}
          </a>
        );
      },
    },
    {
      title: (
        <FormattedMessage
          id="pages.searchTable.titleUpdatedAt"
          defaultMessage="file upload time"
        />
      ),
      sorter: true,
      dataIndex: 'upload_time',
      valueType: 'dateTime',
      search: false, //不让显示到查询表单中
      renderFormItem: (item, { defaultRender, ...rest }, form) => {
        const status = form.getFieldValue('status');
        if (`${status}` === '0') {
          return false;
        }
        if (`${status}` === '3') {
          return (
            <Input
              {...rest}
              placeholder={intl.formatMessage({
                id: 'pages.searchTable.exception',
                defaultMessage: 'Please enter the reason for the exception!',
              })}
            />
          );
        }
        return defaultRender(item);
      },
    },
    {
      title: (
        <FormattedMessage
          id="pages.searchTable.fileSize"
          defaultMessage="Number of service calls"
        />
      ),
      dataIndex: 'file_size',
      sorter: true,
      hideInForm: true,
      search: false, //不让显示到查询表单中
      // renderText: (val: string) =>
      //   `${val}${intl.formatMessage({
      //     id: 'pages.searchTable.tenThousand',
      //     defaultMessage: ' 万 ',
      //   })}`,
    },    
    {
      title: <FormattedMessage id="pages.searchTable.ip" defaultMessage="ip" />,
      dataIndex: 'ip',
      valueType: 'textarea',
    },
    {
      title: <FormattedMessage id="pages.searchTable.appkey" defaultMessage="appkey" />,
      dataIndex: 'appkey',
      valueType: 'textarea',
    },
    {
      title: <FormattedMessage id="pages.searchTable.fileMD5" defaultMessage="file md5" />,
      dataIndex: 'file_md5',
      copyable: true,
      tip: 'The file md5 is the unique key',

      valueType: 'textarea',
    },
    {
      title: <FormattedMessage id="pages.searchTable.file_created_time" defaultMessage="ip" />,
      dataIndex: 'file_created_time',
      valueType: 'dateTime',
      hideInTable: true, // 初始隐藏列
      render: (text, record) => {
        const timestamp = Number(record.file_created_time); // 将 timestamp 转换为数字类型
        if (isNaN(timestamp) || timestamp === undefined) {
          return '';
        }
        const dateTime = new Date(timestamp * 1000);
        const formattedDateTime = dateTime.toLocaleString(); // 格式化日期时间
        return formattedDateTime;
      },
    },
    {
      title: <FormattedMessage id="pages.searchTable.file_modified_time" defaultMessage="ip" />,
      dataIndex: 'file_modified_time',
      valueType: 'dateTime',
      hideInTable: true, // 初始隐藏列
      render: (text, record) => {
        const timestamp = Number(record.file_modified_time); // 将 timestamp 转换为数字类型
        if (isNaN(timestamp) || timestamp === undefined) {
          return '';
        }
        const dateTime = new Date(timestamp * 1000);
        const formattedDateTime = dateTime.toLocaleString(); // 格式化日期时间
        return formattedDateTime;
      },
    },


  ];



  const currentUrl = window.location.href;
  const currentIpAndPort = currentUrl.split('/')[2]; // 获取当前页面的 IP 和端口部分


  return (
    <PageContainer>
      <ProTable<API.UploadHistoryLogItem, API.PageParams>
        headerTitle={intl.formatMessage({
          id: 'pages.searchTable.title',
          defaultMessage: 'File upload log',
        })}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
          
        }}
        toolBarRender={() => [
          // <Button
          //   type="primary"
          //   key="primary"
          //   onClick={() => {
          //     handleModalVisible(true);
          //   }}
          // >
          //   <PlusOutlined /> <FormattedMessage id="pages.searchTable.new" defaultMessage="New" />
          // </Button>,
        ]}
        // ProTable组件的request属性是一个非常重要的API，它接收一个对象作为参数。

        // 该对象中必须包含以下属性：
        // * data：用于存储数据列表的数组。
        // * success：表示数据请求是否成功的布尔值。
        // * total（可选）：用于手动分页的总记录数。
        request={async (params) => {
            try {
            const result = await uploadHistoryLogs(params);
            console.log("uploadHistoryLogs,response: ", result);
            console.log("xxx: ", result?.data);
              return {
                data: result?.data?.list || [],
                success: result?.errorCode === 0, // 根据实际情况判断请求是否成功
                total: result?.data?.total || 0, // 如果需要手动分页，可以提供总记录数
              };
            } catch (error) {
            console.log("Error occurred: ", error);
              return {
                data: [],
                success: false,
                total: 0,
              };
            }
        }}
        // 自定义分页器
        pagination={{
          pageSize: 15, /*每页条数：如果不配置的话，默认是每页20条显示*/
          showTotal: (total) => `共${total}条`, /*总条数的展示*/
        }}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              <FormattedMessage id="pages.searchTable.chosen" defaultMessage="Chosen" />{' '}
              <a style={{ fontWeight: 600 }}>{selectedRowsState.length}</a>{' '}
              <FormattedMessage id="pages.searchTable.item" defaultMessage="项" />
              &nbsp;&nbsp;
              {/* <span>
                <FormattedMessage
                  id="pages.searchTable.totalServiceCalls"
                  defaultMessage="Total number of file"
                />{' '}
                {selectedRowsState.reduce((pre, item) => pre + item.callNo!, 0)}{' '}
                <FormattedMessage id="pages.searchTable.tenThousand" defaultMessage="万" />
              </span> */}
            </div>
          }
        >
          <Button
            onClick={async () => {
              await handleRemove(selectedRowsState);
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
          >
            <FormattedMessage
              id="pages.searchTable.batchDeletion"
              defaultMessage="Batch deletion"
            />
          </Button>

        </FooterToolbar>
      )}
      <UpdateForm
        onSubmit={async (value) => {
          // const success = await handleUpdate(value);
          // if (success) {
          //   handleUpdateModalVisible(false);
          //   setCurrentRow(undefined);
          //   if (actionRef.current) {
          //     actionRef.current.reload();
          //   }
          // }
        }}
        onCancel={() => {
          handleUpdateModalVisible(false);
          if (!showDetail) {
            setCurrentRow(undefined);
          }
        }}
        updateModalVisible={updateModalVisible}
        values={currentRow || {}}
      />

      <Drawer
        width={600}
        visible={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {currentRow?.id && (


          <ProDescriptions<API.UploadHistoryLogItem>
            column={2}
            title={currentRow?.filename}
            dataSource={currentRow} // 设置数据源为 currentRow
            request={async () => {
              // 发送请求判断文件是否存在
              const response = await getFileInfo({ filename: currentRow?.filename, file_md5: currentRow?.file_md5 });
              console.log("getFileInfo response: ", response)
              console.log('currentRow:', currentRow);
              if (response?.errorCode === 22) {
                console.log('currentRow2:', currentRow);
                // 文件不存在，返回错误信息
                return {
                  data: {...currentRow, isFileExist: true} ,
                };
              }

              console.log('currentRow3:', currentRow);

              return {
                data: {...currentRow, isFileExist: true} ,
              };
            }}
            // 右上角显示
            // extra={<Button type="link">修改</Button>}

            columns={columns as ProDescriptionsItemProps<API.UploadHistoryLogItem>[]}
          >
            {/* <ProDescriptions.Item
              dataIndex="percent"
              label="百分比"
              valueType="percent"
            >
              100
            </ProDescriptions.Item>
            <div>多余的dom</div> */}
            <ProDescriptions.Item key="download" label="下载">
              {/* Tooltip文字提示组件是鼠标移动到某个组件悬浮显示字符串用！ */}
              <Tooltip title={isDownloadButtonDisabled ? '文件已经在服务端被删除☘️' : ''}>
                <div>
                  <Button
                    type="primary"
                    icon={<DownloadOutlined />}
                    disabled={isDownloadButtonDisabled} // 根据 isButtonDisabled 决定按钮是否置灰
                    onClick={() => {
                      // 处理下载逻辑
                      const downloadUrl = `http://${currentIpAndPort}/api/v1/file/download?filename=${currentRow?.filename}`;
                      window.open(downloadUrl, '_blank');
                    }}
                  >
                    下载
                  </Button>
                </div>
              </Tooltip>
            </ProDescriptions.Item>
          </ProDescriptions>

          
        )}
      </Drawer>
    </PageContainer>
  );
};

export default TableList;
