import { getFileInfo, removeUploadHistoryLog, uploadHistoryLogs } from '@/services/ant-design-pro/api';
import { DownloadOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, message, Tooltip } from 'antd';
import React, { useEffect, useRef, useState } from 'react';
import { FormattedMessage, useIntl } from 'umi';

const handleRemove = async (selectedRows: API.UploadHistoryLogItem[]) => {
  const hide = message.loading('正在删除');
  if (!selectedRows) return true;
  try {
    await removeUploadHistoryLog({
      id: selectedRows.map((row) => row.id),
    });
    hide();
    message.success('删除成功');
    return true;
  } catch (error) {
    hide();
    message.error('删除失败，请重试');
    return false;
  }
};

const TableList: React.FC = () => {
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.UploadHistoryLogItem>();
  const [selectedRowsState, setSelectedRows] = useState<API.UploadHistoryLogItem[]>([]);
  const [isDownloadButtonDisabled, setIsDownloadButtonDisabled] = useState(true);
  const intl = useIntl();

  useEffect(() => {
    if (!currentRow?.filename) return;
    const fetchData = async () => {
      const fileInfo = await getFileInfo({
        filename: currentRow?.filename,
        file_md5: currentRow?.file_md5,
      });
      setIsDownloadButtonDisabled(fileInfo?.errorCode === 22);
    };
    fetchData();
  }, [currentRow?.filename, currentRow?.file_md5]);

  const columns: ProColumns<API.UploadHistoryLogItem>[] = [
    {
      title: <FormattedMessage id="pages.searchTable.updateForm.fileName.nameLabel" />,
      dataIndex: 'filename',
      render: (dom, entity) => (
        <a
          onClick={() => {
            setCurrentRow(entity);
            setShowDetail(true);
          }}
        >
          {dom}
        </a>
      ),
    },
    {
      title: <FormattedMessage id="pages.searchTable.titleUpdatedAt" defaultMessage="上传时间" />,
      sorter: true,
      dataIndex: 'upload_time',
      valueType: 'dateTime',
      search: false,
    },
    {
      title: <FormattedMessage id="pages.searchTable.fileSize" defaultMessage="文件大小" />,
      dataIndex: 'file_size',
      sorter: true,
      hideInForm: true,
      search: false,
    },
    {
      title: <FormattedMessage id="pages.searchTable.ip" defaultMessage="IP" />,
      dataIndex: 'ip',
      valueType: 'textarea',
    },
    {
      title: <FormattedMessage id="pages.searchTable.appkey" defaultMessage="appkey" />,
      dataIndex: 'appkey',
      valueType: 'textarea',
    },
    {
      title: <FormattedMessage id="pages.searchTable.fileMD5" defaultMessage="文件MD5" />,
      dataIndex: 'file_md5',
      copyable: true,
      tip: '文件MD5是唯一标识',
      valueType: 'textarea',
    },
    {
      title: <FormattedMessage id="pages.searchTable.file_created_time" defaultMessage="文件创建时间" />,
      dataIndex: 'file_created_time',
      valueType: 'dateTime',
      hideInTable: true,
      render: (_text, record) => {
        const timestamp = Number(record.file_created_time);
        if (isNaN(timestamp)) return '';
        return new Date(timestamp * 1000).toLocaleString();
      },
    },
    {
      title: <FormattedMessage id="pages.searchTable.file_modified_time" defaultMessage="文件修改时间" />,
      dataIndex: 'file_modified_time',
      valueType: 'dateTime',
      hideInTable: true,
      render: (_text, record) => {
        const timestamp = Number(record.file_modified_time);
        if (isNaN(timestamp)) return '';
        return new Date(timestamp * 1000).toLocaleString();
      },
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.UploadHistoryLogItem, API.PageParams>
        headerTitle={intl.formatMessage({
          id: 'pages.searchTable.title',
          defaultMessage: '上传日志查询',
        })}
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 120 }}
        request={async (params) => {
          try {
            const result = await uploadHistoryLogs(params);
            return {
              data: result?.data?.list || [],
              success: result?.errorCode === 0,
              total: result?.data?.total || 0,
            };
          } catch (error) {
            return { data: [], success: false, total: 0 };
          }
        }}
        pagination={{
          pageSize: 15,
          showTotal: (total) => `共 ${total} 条`,
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
              <FormattedMessage id="pages.searchTable.chosen" defaultMessage="已选择" />{' '}
              <a style={{ fontWeight: 600 }}>{selectedRowsState.length}</a>{' '}
              <FormattedMessage id="pages.searchTable.item" defaultMessage="项" />
            </div>
          }
        >
          <Button
            danger
            onClick={async () => {
              await handleRemove(selectedRowsState);
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
          >
            <FormattedMessage id="pages.searchTable.batchDeletion" defaultMessage="批量删除" />
          </Button>
        </FooterToolbar>
      )}

      <Drawer
        width={600}
        open={showDetail}
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
            dataSource={currentRow}
            request={async () => ({
              data: { ...currentRow, isFileExist: true },
            })}
            columns={columns as ProDescriptionsItemProps<API.UploadHistoryLogItem>[]}
          >
            <ProDescriptions.Item key="download" label="下载">
              <Tooltip title={isDownloadButtonDisabled ? '文件已在服务端被删除' : ''}>
                <Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  disabled={isDownloadButtonDisabled}
                  onClick={() => {
                    const downloadUrl = `/api/v1/file/download?filename=${currentRow?.filename}`;
                    window.open(downloadUrl, '_blank');
                  }}
                >
                  下载
                </Button>
              </Tooltip>
            </ProDescriptions.Item>
          </ProDescriptions>
        )}
      </Drawer>
    </PageContainer>
  );
};

export default TableList;
