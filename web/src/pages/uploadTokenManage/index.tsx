import { createUploadToken, generateAppSecret, removeUploadToken, saveUploadToken, uploadTokenLists } from '@/services/ant-design-pro/api';
import type {
  EditableFormInstance,
  ProColumns,
  ProFormInstance,
} from '@ant-design/pro-components';
import {
  EditableProTable,
  ProCard,
  ProForm,
  ProFormDependency,
  ProFormField,
  ProFormSegmented,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { Button, Input, Modal, message } from 'antd';
import React, { useRef, useState } from 'react';
import { FormattedMessage } from 'umi';



const defaultData: API.UploadTokenItem[] = [
  // {
  //   id: 624748504,
  //   appkey: 'asdfafa',
  //   appsecret: 'aadfafdafaa',
  //   desc: '',
  //   state: 'open',
  //   created_at: 1590486176000,
  // },
  // {
  //   id: 624748505,
  //   appkey: 'agxvdfaaafs',
  //   appsecret: 'adsadsfadfafadfadf',
  //   desc: '',
  //   state: 'open',
  //   created_at: 1590486176000,
  // },
];

const handleSave =async ( key: RecordKey,
  record: API.UploadTokenItem,
  originRow: API.UploadTokenItem, ) => {

    try {
      // 在这里编写保存数据的逻辑，可以发送请求将数据保存到后台
      console.log('要保存的行数据：', record);
      console.log('原始行数据：', originRow);
      console.log('要更新的行的 key：', key);
  
      // 这里发送请求将数据保存到后台
      await saveUploadToken(record);
  
      // 这里返回保存成功后的提示信息，可以根据需要进行修改
      message.success('Save successfully');

      return '保存成功';
    } catch (error) {
      // 这里处理保存失败的情况
      console.error('保存失败：', error);
      message.error('Save failed, please try again');
      throw new Error('保存失败');
    }
}


let i = 0;

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
  const [isModalOpen, setIsModalOpen] = useState(false);
  const showModal = () => {
    setIsModalOpen(true);
  };

  const handleOk = () => {
    setIsModalOpen(false);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };


  const [editableKeys, setEditableRowKeys] = useState<React.Key[]>(() => []);
  const [position, setPosition] = useState<'top' | 'bottom' | 'hidden'>(
    'bottom',
  );
  const [controlled, setControlled] = useState<boolean>(false);
  const formRef = useRef<ProFormInstance<any>>();
  const editorFormRef = useRef<EditableFormInstance<API.UploadTokenItem>>();
  const columns: ProColumns<API.UploadTokenItem>[] = [
    {
      title: (
        <FormattedMessage
          id="pages.uploadTokenManage.title.appKey"
        />
      ),
      copyable: true,
      dataIndex: 'appkey',
      formItemProps: () => {
        return {
          rules: [
            { required: true, message: '此项为必填项' },
            { pattern: /^[a-zA-Z0-9_]+$/, message: '请输入数字、字母和下划线' }
          ],
        };
      },
      width: '15%',
    },
    {
      title: (
        <FormattedMessage
          id="pages.uploadTokenManage.title.appSecret"
        />
      ),
      copyable: true,
      dataIndex: 'appsecret',
      width: '23%',
      formItemProps: {
        rules: [
          {
            required: true,
            whitespace: true,
            message: '此项是必填项',
          },
          { pattern: /^[a-zA-Z0-9_]+$/, message: '请输入数字、字母和下划线' },
          {
            max: 32,
            whitespace: true,
            message: '最长为 32 位',
          },
          {
            min: 6,
            whitespace: true,
            message: '最小为 6 位',
          },
        ],
      },
    },       
    {
      title: (
        <FormattedMessage
          id="pages.uploadTokenManage.title.state"
        />
      ),
      key: 'state',
      width: '10%',
      dataIndex: 'state',
      valueType: 'select',
      valueEnum: {
        // all: { text: '全部', status: 'Default' },
        open: {
          text: '启用',
          status: 'Success',
        },
        closed: {
          text: '禁用',
          status: 'Error',
        },
      },
    },
    {
      title: (
        <FormattedMessage
          id="pages.uploadTokenManage.title.desc"
        />
      ),
      width: '30%',
      dataIndex: 'desc',
    },
    {
      title: (
        <FormattedMessage
          id="pages.uploadTokenManage.title.uploadToken"
        />
      ),
      copyable: true,
      editable: false,
      dataIndex: 'uploadToken',
      width: '23%',
      tip: '使用appkey、appsecret生成上传文件token，生成后请及时复制！',
      formItemProps: {
        rules: [

          { pattern: /^[a-zA-Z0-9_:]+$/, message: '请输入数字、字母和下划线' },

        ],
      },
    },  
    // {
    //   title: (
    //     <FormattedMessage
    //       id="pages.uploadTokenManage.title.createTime"
    //       defaultMessage="create time"
    //     />
    //   ),
    //   width: '15%',
    //   dataIndex: 'created_at',
    //   valueType: 'date',
    // },  
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      render: (text, record, _, action) => {
        // 判断是否为系统内置记录
        const isSysBuilt = record.is_sys_built;
        
        
       return [
        isSysBuilt ? null : (
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

        isSysBuilt ? null : (
          <a
            key="delete"
            onClick={async () => {
              const tableDataSource = formRef.current?.getFieldValue('table') as API.UploadTokenItem[];
              formRef.current?.setFieldsValue({
                table: tableDataSource.filter((item) => item.appkey !== record.appkey),
              });
  
              console.log("tableDataSourc: ", tableDataSource);
              console.log("delete record.appkey: ", record.appkey);
              console.log("delete record: ", record);
  
              await removeUploadToken(record);
            }}
          >
            删除
          </a>
        ),

        <a
        key="getUploadToken"
        onClick={async () => {
          const hide = message.loading('正在生成上传token');

          try {

            console.log("getUploadToken record.appkey: ", record.appkey)
            console.log("getUploadToken record: ", record)

            const result = await createUploadToken(record)
            console.log("result: ",result)

            if (result) {
              console.log("result.data: ",result.data)


              editorFormRef.current?.setRowData?.(record.appkey!, {
                uploadToken: result.data,
              });

            }
            

            hide();

            message.success('successfully');
            return true;
          } catch (error) {
            hide();
            message.error('生成上传token, please try again');
            return false;
          }


        }}
      >
        获取上传token
        
      </a>,

      ]
    }
    },
  ];

  return (
    <ProForm<{
      table: API.UploadTokenItem[];
    }>
      submitter={false} // 不展示提交按钮和重置按钮
      formRef={formRef}
      // initialValues={{
      //   table: defaultData,  //来自requst，这里暂时不需要了
      // }}
      validateTrigger="onBlur"
    >
      <EditableProTable<API.UploadTokenItem>
        rowKey="appkey"
        scroll={{
          x: 960,
        }}
        editableFormRef={editorFormRef}
        headerTitle="上传token管理" // 表格标题
        maxLength={8} // 表格项最大多少个
        name="table"
        controlled={controlled}
        recordCreatorProps={
          position !== 'hidden'
            ? {
                position: position as 'top',
                record: () => ({ 
                  appkey: "admin_" + (Math.random() * 1000000).toFixed(0),
                  state: "open",
                  appsecret: createAppSecret(),
                }),
              }
            : false
        }
        // 上边操作也暂时不需要了
        // toolBarRender={() => [
        //   <ProFormSwitch
        //     key="render"
        //     fieldProps={{
        //       style: {
        //         marginBlockEnd: 0,
        //       },
        //       checked: controlled,
        //       onChange: (value) => {
        //         setControlled(value);
        //       },
        //     }}
        //     checkedChildren="数据更新通知 Form"
        //     unCheckedChildren="保存后通知 Form"
        //     noStyle
        //   />,
        //   <ProFormSegmented
        //     key="render"
        //     fieldProps={{
        //       style: {
        //         marginBlockEnd: 0,
        //       },
        //       value: position,
        //       onChange: (value) => {
        //         setPosition(value as 'top');
        //       },
        //     }}
        //     noStyle
        //     request={async () => [
        //       {
        //         label: '添加到顶部',
        //         value: 'top',
        //       },
        //       {
        //         label: '添加到底部',
        //         value: 'bottom',
        //       },
        //       {
        //         label: '隐藏',
        //         value: 'hidden',
        //       },
        //     ]}
        //   />,
        //   <Button
        //     key="rows"
        //     onClick={() => {
        //       const rows = editorFormRef.current?.getRowsData?.();
        //       console.log(rows);
        //     }}
        //   >
        //     获取 table 的数据
        //   </Button>,
        // ]}
        request={uploadTokenLists}
        columns={columns}
        editable={{
          type: 'single',
          editableKeys,
          onSave: handleSave,
          onChange: setEditableRowKeys,
          actionRender: (row, config, defaultDom) => {
            return [
              defaultDom.save,
              defaultDom.cancel,
              <a
                key="id"
                onClick={() => {
                  console.log(config.index);
                  i++;
                  editorFormRef.current?.setRowData?.(config.index!, {
                    appsecret: createAppSecret(),
                  });
                }}
              >
                更新appSecret
              </a>,
            ];
          },
        }}
      />
      
      {/* <ProForm.Item>
        <ProCard title="表格数据" headerBordered collapsible defaultCollapsed>
          <ProFormDependency name={['table']}>
            {({ table }) => {
              return (
                <ProFormField
                  ignoreFormItem
                  fieldProps={{
                    style: {
                      width: '100%',
                    },
                  }}
                  mode="read"
                  valueType="jsonCode"
                  text={JSON.stringify(table)}
                />
              );
            }}
          </ProFormDependency>
        </ProCard>
      </ProForm.Item> */}
    </ProForm>
  );
};