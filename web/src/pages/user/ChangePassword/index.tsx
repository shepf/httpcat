import Footer from '@/components/Footer';
import { changePasswd } from '@/services/ant-design-pro/api';
import { LockOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { Alert, Card, message, Typography } from 'antd';
import React from 'react';
import { history, useModel } from 'umi';

const loginPath = '/user/login';

const ChangePasswordPage: React.FC = () => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const forceChangePassword = !!currentUser?.mustChangePassword;

  React.useEffect(() => {
    if (initialState && !currentUser) {
      history.replace(loginPath);
      return;
    }
    if (initialState && currentUser && !forceChangePassword) {
      history.replace('/');
    }
  }, [currentUser, forceChangePassword, initialState]);

  const handleSubmit = async (values: { oldPassword: string; newPassword: string; confirmPassword: string }) => {
    if (values.newPassword !== values.confirmPassword) {
      message.error('两次输入的密码不一致');
      return;
    }

    const response = await changePasswd({
      oldPassword: values.oldPassword,
      newPassword: values.newPassword,
    });

    if (response.errorCode !== 0) {
      message.error(response.msg || '密码修改失败');
      return;
    }

    message.success('密码修改成功，请重新登录');
    localStorage.removeItem('token');
    await setInitialState((s) => ({ ...s, currentUser: undefined }));
    history.replace(loginPath);
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
        background: '#f5f7fa',
      }}
    >
      <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: 24 }}>
        <Card style={{ width: '100%', maxWidth: 480, borderRadius: 12 }}>
          <Typography.Title level={3} style={{ marginBottom: 8 }}>
            首次登录请修改默认密码
          </Typography.Title>
          <Typography.Paragraph type="secondary" style={{ marginBottom: 24 }}>
            为了保证系统安全，当前管理员账号必须先完成密码修改后才能继续使用。
          </Typography.Paragraph>
          <Alert
            type="warning"
            showIcon
            style={{ marginBottom: 24 }}
            message={`当前账号：${currentUser?.name || 'admin'}，请将默认密码修改为新的安全密码。`}
          />
          <LoginForm
            submitter={{
              searchConfig: { submitText: '立即修改' },
              resetButtonProps: false,
            }}
            onFinish={async (values) => {
              await handleSubmit(values as { oldPassword: string; newPassword: string; confirmPassword: string });
            }}
          >
            <ProFormText.Password
              name="oldPassword"
              label="当前密码"
              fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
              rules={[{ required: true, message: '请输入当前密码' }]}
            />
            <ProFormText.Password
              name="newPassword"
              label="新密码"
              fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
              rules={[{ required: true, message: '请输入新密码' }]}
            />
            <ProFormText.Password
              name="confirmPassword"
              label="确认新密码"
              fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
              rules={[{ required: true, message: '请再次输入新密码' }]}
            />
          </LoginForm>
        </Card>
      </div>
      <Footer />
    </div>
  );
};

export default ChangePasswordPage;
