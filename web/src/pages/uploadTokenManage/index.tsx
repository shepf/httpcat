import { createUploadToken, removeUploadToken, saveUploadToken, uploadTokenLists } from '@/services/ant-design-pro/api';
import type { EditableFormInstance, ProColumns, ProFormInstance } from '@ant-design/pro-components';
import { EditableProTable, ProForm } from '@ant-design/pro-components';
import { message } from 'antd';
import React, { useRef, useState } from 'react';
import { FormattedMessage } from 'umi';

type RecordKey = React.Key;

const handleSave = async (
  key: RecordKey,
  record: API.UploadTokenItem,
  originRow: API.UploadTokenItem,
) => {
  try {
    await saveUploadToken(record);
    message.success('保存成功');
    return '保存成功';
  } catch (error) {
    message.error('保存失败，请重试');
    throw new Error('保存失败');
  }
};

function createAppSecret() {
  const characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  let randomString = '';
  for (let j = 0; j < 32; j++) {
    const randomIndex = Math.floor(Math.random() * characters.length);
    randomString += characters[randomIndex];
  }
  return randomString;
}

export default () => {
  const [editableKeys, setEditableRowKeys] = useState<React.Key[]>(() => []);
  const formRef = useRef<ProFormInstance<any>>();
  const editorFormRef = useRef<EditableFormInstance<API.UploadTokenItem>>();

  const columns: ProColumns<API.UploadTokenItem>[] = [
    {
      title: <FormattedMessage id="pages.uploadTokenManage.title.appKey" />,
      copyable: true,
      dataIndex: 'appkey',
      formItemProps: () => ({
        rules: [
          { required: true, message: '此项为必填项' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: '请输入数字、字母和下划线' },
        ],
      }),
      width: '15%',
    },
    {
      title: <FormattedMessage id="pages.uploadTokenManage.title.appSecret" />,
      copyable: true,
      dataIndex: 'appsecret',
      width: '23%',
      formItemProps: {
        rules: [
          { required: true, whitespace: true, message: '此项是必填项' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: '请输入数字、字母和下划线' },
          { max: 32, whitespace: true, message: '最长为 32 位' },
          { min: 6, whitespace: true, message: '最小为 6 位' },
        ],
      },
    },
    {
      title: <FormattedMessage id="pages.uploadTokenManage.title.state" />,
      key: 'state',
      width: '10%',
      dataIndex: 'state',
      valueType: 'select',
      valueEnum: {
        open: { text: '启用', status: 'Success' },
        closed: { text: '禁用', status: 'Error' },
      },
    },
    {
      title: <FormattedMessage id="pages.uploadTokenManage.title.desc" />,
      width: '30%',
      dataIndex: 'desc',
    },
    {
      title: <FormattedMessage id="pages.uploadTokenManage.title.uploadToken" />,
      copyable: true,
      editable: false,
      dataIndex: 'uploadToken',
      width: '23%',
      tip: '使用 appkey、appsecret 生成上传Token，生成后请及时复制',
    },
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      render: (_text, record, _, action) => {
        const isSysBuilt = record.is_sys_built;

        return [
          !isSysBuilt && (
            <a
              key="editable"
              onClick={() => {
                if (record.appkey !== undefined) {
                  action?.startEditable?.(record.appkey);
                }
              }}
            >
              编辑
            </a>
          ),
          !isSysBuilt && (
            <a
              key="delete"
              onClick={async () => {
                const tableDataSource = formRef.current?.getFieldValue('table') as API.UploadTokenItem[];
                formRef.current?.setFieldsValue({
                  table: tableDataSource.filter((item) => item.appkey !== record.appkey),
                });
                await removeUploadToken(record);
              }}
            >
              删除
            </a>
          ),
          <a
            key="getUploadToken"
            onClick={async () => {
              const hide = message.loading('正在生成上传Token');
              try {
                const result = await createUploadToken(record);
                if (result?.data) {
                  editorFormRef.current?.setRowData?.(record.appkey!, {
                    uploadToken: result.data,
                  });
                }
                hide();
                message.success('Token 生成成功');
              } catch (error) {
                hide();
                message.error('Token 生成失败，请重试');
              }
            }}
          >
            获取上传Token
          </a>,
        ].filter(Boolean);
      },
    },
  ];

  return (
    <ProForm<{ table: API.UploadTokenItem[] }>
      submitter={false}
      formRef={formRef}
      validateTrigger="onBlur"
    >
      <EditableProTable<API.UploadTokenItem>
        rowKey="appkey"
        scroll={{ x: 960 }}
        editableFormRef={editorFormRef}
        headerTitle="上传Token管理"
        maxLength={8}
        name="table"
        recordCreatorProps={{
          position: 'bottom',
          record: () => ({
            appkey: 'admin_' + (Math.random() * 1000000).toFixed(0),
            state: 'open',
            appsecret: createAppSecret(),
          }),
        }}
        request={uploadTokenLists}
        columns={columns}
        editable={{
          type: 'single',
          editableKeys,
          onSave: handleSave,
          onChange: setEditableRowKeys,
          actionRender: (row, config, defaultDom) => [
            defaultDom.save,
            defaultDom.cancel,
            <a
              key="updateSecret"
              onClick={() => {
                editorFormRef.current?.setRowData?.(config.index!, {
                  appsecret: createAppSecret(),
                });
              }}
            >
              更新appSecret
            </a>,
          ],
        }}
      />
    </ProForm>
  );
};
